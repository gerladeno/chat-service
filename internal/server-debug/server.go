package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/gerladeno/chat-service/internal/buildinfo"
	"github.com/gerladeno/chat-service/internal/logger"
	"github.com/gerladeno/chat-service/internal/validator"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr      string      `option:"mandatory" validate:"required,hostname_port"`
	v1Swagger *openapi3.T `option:"mandatory" validate:"required"`
}

type Server struct {
	lg      *zap.Logger
	srv     *http.Server
	swagger *openapi3.T
}

func New(opts Options) (*Server, error) {
	if err := validator.Validator.Struct(opts); err != nil {
		return nil, fmt.Errorf("validating server-debug options: %w", err)
	}
	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.Use(middleware.Recover())

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		swagger: opts.v1Swagger,
	}
	index := newIndexPage()
	e.GET("/version", s.Version)
	index.addPage("/version", "Get build information")
	e.PUT("/log/level", s.LogLevel)
	wrap(e)
	index.addPage("/debug/pprof/", "Go std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")
	e.GET("/debug/error", s.SendErrorEvent)
	index.addPage("/debug/error", "Debug Sentry error event")
	e.GET("/schema/client", s.SchemaClient)
	index.addPage("/schema/client", "Get client OpenAPI specification")

	e.GET("/", index.handler)
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}

func (s *Server) Version(eCtx echo.Context) error {
	if err := eCtx.JSON(http.StatusOK, buildinfo.BuildInfo); err != nil {
		return fmt.Errorf("sending version: %w", err)
	}
	return nil
}

func (s *Server) LogLevel(eCtx echo.Context) error {
	level := eCtx.FormValue("level")
	if err := logger.Atom.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("parsing level %s: %w", level, err)
	}
	s.lg.Info(fmt.Sprintf("Log level changed to %s", level))
	return nil
}

func (s *Server) SendErrorEvent(eCtx echo.Context) error {
	s.lg.Error("look for me in the sentry")
	if err := eCtx.String(http.StatusOK, "event sent"); err != nil {
		return fmt.Errorf("sending error event text: %w", err)
	}
	return nil
}

func (s *Server) SchemaClient(eCtx echo.Context) error {
	data, err := s.swagger.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshalling swager json: %w", err)
	}
	if err = eCtx.Blob(http.StatusOK, "application/json", data); err != nil {
		return fmt.Errorf("sending data: %w", err)
	}
	return nil
}

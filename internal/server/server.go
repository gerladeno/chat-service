package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/gerladeno/chat-service/internal/middlewares"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
	bodyLimit         = "12KB"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger               *zap.Logger              `option:"mandatory" validate:"required"`
	addr                 string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins         []string                 `option:"mandatory" validate:"min=1"`
	registerHandlersFunc func(e *echo.Echo)       `option:"mandatory" validate:"required"`
	client               middlewares.Introspector `option:"mandatory" validate:"required"`
	resource             string                   `option:"mandatory" validate:"required"`
	role                 string                   `option:"mandatory" validate:"required"`
	wsSecProtocol        string                   `option:"mandatory" validate:"required"`
	errHandler           echo.HTTPErrorHandler    `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating server options: %w", err)
	}

	e := echo.New()
	e.HTTPErrorHandler = opts.errHandler
	e.Use(
		middleware.Recover(),
		middleware.BodyLimit(bodyLimit),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.allowOrigins,
			AllowMethods: []string{http.MethodPost},
		}),
		middlewares.NewKeycloakTokenAuth(opts.client, opts.resource, opts.role, opts.wsSecProtocol),
		middlewares.RequestLogger(opts.logger),
	)
	opts.registerHandlersFunc(e)

	s := Server{
		lg: opts.logger,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
	return &s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer s.lg.Info("stopped")
		<-ctx.Done()
		gfCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := s.srv.Shutdown(gfCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %w", err)
		}
		return nil
	})

	return eg.Wait()
}

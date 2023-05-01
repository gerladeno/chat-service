package serverclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/gerladeno/chat-service/internal/middlewares"
	clientv1 "github.com/gerladeno/chat-service/internal/server-client/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
	bodyLimit         = "12KB"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger       *zap.Logger              `option:"mandatory" validate:"required"`
	addr         string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins []string                 `option:"mandatory" validate:"min=1"`
	v1Swagger    *openapi3.T              `option:"mandatory" validate:"required"`
	v1Handlers   clientv1.ServerInterface `option:"mandatory" validate:"required"`
	client       middlewares.Introspector `option:"mandatory" validate:"required"`
	resource     string                   `option:"mandatory" validate:"required"`
	role         string                   `option:"mandatory" validate:"required"`
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
	e.Use(
		middleware.Recover(),
		middleware.BodyLimit(bodyLimit),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if c.Request().Method == http.MethodOptions {
					return nil
				}
				opts.logger.Info("request",
					zap.Duration("latency", v.Latency),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("host", v.Host),
					zap.String("method", v.Method),
					zap.String("path", v.RoutePath),
					zap.String("request_id", v.RequestID),
					zap.String("user_agent", v.UserAgent),
					zap.Int("status", v.Status),
					zap.String("userId", middlewares.MustUserID(c).String()),
				)
				return nil
			},
			LogLatency:   true,
			LogRemoteIP:  true,
			LogHost:      true,
			LogMethod:    true,
			LogRoutePath: true,
			LogRequestID: true,
			LogUserAgent: true,
			LogStatus:    true,
			LogError:     true,
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.allowOrigins,
			AllowMethods: []string{http.MethodPost},
		}),
		middlewares.NewKeycloakTokenAuth(opts.client, opts.resource, opts.role),
	)

	v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:  false,
			ExcludeResponseBody: true,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
	}))
	clientv1.RegisterHandlers(v1, opts.v1Handlers)

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
package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/gerladeno/chat-service/internal/errors"
)

func RequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if c.Request().Method == http.MethodOptions {
				return nil
			}
			status := v.Status
			if v.Error != nil {
				status = errors.GetServerErrorCode(v.Error)
			}
			logger.Info("request",
				zap.Duration("latency", v.Latency),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("path", v.RoutePath),
				zap.String("request_id", v.RequestID),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", status),
				zap.String("userId", MustUserID(c).String()),
				zap.Error(v.Error),
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
	})
}

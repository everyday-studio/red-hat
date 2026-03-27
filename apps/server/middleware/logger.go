package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/everyday-studio/redhat/kit/contexts"
)

func LoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{

		BeforeNextFunc: func(c echo.Context) {
			req := c.Request()
			requestID := contexts.GetRequestID(c.Request().Context())
			ctxLogger := logger.With(slog.String("request_id", requestID))
			reqCtx := contexts.WithLogger(req.Context(), ctxLogger)
			c.SetRequest(req.WithContext(reqCtx))
		},

		LogRequestID: true,
		LogStatus:    true,
		LogMethod:    true,
		LogURIPath:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		HandleError:  true,

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			baseLogger := logger.With(
				slog.String("request_id", v.RequestID),
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("path", v.URIPath),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("user_agent", v.UserAgent),
				slog.Float64("latency_ms", float64(v.Latency.Nanoseconds())/1e6),
			)

			if detailErr, ok := c.Get("detail_error").(error); ok {
				baseLogger = baseLogger.With(slog.String("detail_error", detailErr.Error()))
			} else if v.Error != nil {
				baseLogger = baseLogger.With(slog.String("err", v.Error.Error()))
			}

			switch {
			case v.Status >= 500:
				baseLogger.Error("SERVER_ERROR")
			case v.Status >= 400:
				baseLogger.Info("CLIENT_ERROR")
			default:
				baseLogger.Info("REQUEST_SUCCESS")
			}
			return nil
		},
	})
}

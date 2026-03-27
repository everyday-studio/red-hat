package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/everyday-studio/redhat/config"
	"github.com/everyday-studio/redhat/kit/contexts"
	"github.com/everyday-studio/redhat/kit/security"
)

func Setup(cfg *config.Config, logger *slog.Logger, e *echo.Echo) error {
	// Attach a unique request ID to each request
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, requestID string) {
			req := c.Request()
			req.Header.Set(echo.HeaderXRequestID, requestID)

			ctx := contexts.WithRequestID(req.Context(), requestID)
			c.SetRequest(req.WithContext(ctx))
		},
	}))

	// Structured request/response logging
	e.Use(LoggerMiddleware(logger))

	// Panic recovery
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))

	// Handler execution timeout
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Server-level timeouts
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 40 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	// JWT Authentication (RS256, Bearer token only)
	publicKey, err := security.ParseRSAPublicKeyFromBase64(cfg.Secure.JWT.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    publicKey,
		SigningMethod: "RS256",
		TokenLookup:  "header:" + echo.HeaderAuthorization + ":Bearer ",

		// Extract user_id and role from claims and store in echo context.
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Set("user", nil)
				return
			}
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("user_id", userID)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		},

		// Allow requests without a token to reach public endpoints.
		// AllowRoles middleware handles authorization on protected routes.
		ErrorHandler: func(c echo.Context, err error) error {
			if errors.Is(err, echojwt.ErrJWTMissing) {
				return nil
			}

			var statusCode int
			var errorMsg string

			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				statusCode = http.StatusUnauthorized
				errorMsg = "token has expired"
				if cfg.App.Debug {
					logger.Info("JWT token expired", "path", c.Path())
				}
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				statusCode = http.StatusUnauthorized
				errorMsg = "invalid token signature"
				logger.Warn("JWT invalid signature", "path", c.Path())
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				statusCode = http.StatusUnauthorized
				errorMsg = "token not valid yet"
			default:
				statusCode = http.StatusUnauthorized
				errorMsg = "invalid or malformed token"
				if cfg.App.Debug {
					logger.Warn("JWT validation failed", "error", err.Error(), "path", c.Path())
				}
			}

			return echo.NewHTTPError(statusCode, map[string]string{"error": errorMsg})
		},
		ContinueOnIgnoredError: true,
	}))

	return nil
}

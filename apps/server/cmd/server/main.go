package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	appconfig "github.com/everyday-studio/redhat/config"
	"github.com/everyday-studio/redhat/db"
	"github.com/everyday-studio/redhat/handlers"
	"github.com/everyday-studio/redhat/kit/security"
	"github.com/everyday-studio/redhat/middleware"
	"github.com/everyday-studio/redhat/repository/postgres"
	"github.com/everyday-studio/redhat/services"
)

func main() {
	app := fx.New(
		fx.Provide(
			NewConfig,
			NewLogger,
			NewDB,
			echo.New,
			NewAuthService,
		),
		fx.Provide(
			postgres.NewUserRepository,
			services.NewUserService,
		),
		fx.Invoke(
			middleware.Setup,
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
		),
		fx.Invoke(StartServer),
		fx.WithLogger(
			func(cfg *appconfig.Config, logger *slog.Logger) fxevent.Logger {
				if cfg.App.Env == "prod" || !cfg.App.Debug {
					return fxevent.NopLogger
				}
				return &fxevent.SlogLogger{Logger: logger}
			},
		),
	)

	app.Run()
}

func NewConfig() *appconfig.Config {
	envFlag := flag.String("env", "", "Environment (dev, prod)")
	flag.Parse()

	env := *envFlag
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = "dev"
	}

	validEnvs := map[string]bool{"dev": true, "prod": true}
	if !validEnvs[env] {
		log.Fatalf("Invalid environment: %s. Valid environments are: dev, prod", env)
	}

	cfg, err := appconfig.LoadConfig(env)
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	fmt.Printf("Config loaded for environment: %s\n", env)
	if env == "dev" {
		fmt.Printf("config: %+v\n", cfg)
	}

	return cfg
}

func NewLogger(cfg *appconfig.Config) *slog.Logger {
	logLevel := slog.LevelInfo
	switch strings.ToLower(cfg.App.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, time.Now().Format("2006-01-02T15:04:05.000Z07:00"))
			}
			return a
		},
	})
	return slog.New(handler)
}

func NewDB(lc fx.Lifecycle, cfg *appconfig.Config) *sql.DB {
	dbConn, err := db.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return dbConn.Close()
		},
	})

	return dbConn
}

// NewAuthService wires the RSA private key from config into AuthService.
func NewAuthService(cfg *appconfig.Config, userRepo *postgres.UserRepository) *services.AuthService {
	if cfg.Secure.JWT.PrivateKey == "" {
		log.Fatal("JWT private key is not configured (secure.jwt.private_key_base64)")
	}

	privateKey, err := security.ParseRSAPrivateKeyFromBase64(cfg.Secure.JWT.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to parse JWT private key: %v", err)
	}

	ttl := time.Duration(cfg.Secure.JWT.AccessExpirationMin) * time.Minute
	return services.NewAuthService(userRepo, privateKey, ttl, cfg.Steam)
}

func StartServer(lc fx.Lifecycle, e *echo.Echo, cfg *appconfig.Config) {
	if cfg.App.Env == "prod" || !cfg.App.Debug {
		e.HideBanner = true
		e.HidePort = true
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
					log.Fatal("Shutting down the server:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}

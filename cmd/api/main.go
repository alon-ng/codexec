// @title           Codim API
// @version         1.0
// @description     This is the Codim API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
// @description Authentication is performed via HTTP-only cookie named "auth_token" containing a JWT token. The cookie is automatically set on successful signup/login and cleared on logout.
package main

import (
	"codim/cmd/api/config"
	"codim/internal/api/auth"
	"codim/internal/api/v1"
	"codim/internal/db"
	"codim/internal/utils/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	log, err := initializeLogger(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize logger: %v", err)
	}

	log.Info("Logger initialized successfully")

	queries, pool := initializeQueries(cfg, log)
	defer pool.Close()

	log.Info("Database initialized successfully")

	authProvider := auth.NewProvider(
		cfg.API.PasswordSalt,
		[]byte(cfg.API.JwtSecret),
		cfg.API.JwtTTL,
		cfg.API.JwtRenewalThreshold,
	)

	router := api.NewRouter(queries, log, authProvider)

	log.Infof("Router initialized successfully")

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	log.Infof("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Info("Server started successfully")
}

func initializeLogger(cfg config.Config) (*logger.Logger, error) {
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func initializeQueries(cfg config.Config, log *logger.Logger) (*db.Queries, *pgxpool.Pool) {
	pool, err := db.NewPool(context.Background(), cfg.DB, log)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db.New(pool), pool
}

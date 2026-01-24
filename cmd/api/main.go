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

// @typeAlias json.RawMessage object
package main

import (
	"codim/cmd/api/config"
	"codim/pkg/api/auth"
	"codim/pkg/api/v1"
	"codim/pkg/api/v1/websocket"
	"codim/pkg/db"
	"codim/pkg/rabbitmq"
	"codim/pkg/utils/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

	queries, pool := initializeQueries(cfg, log)
	defer pool.Close()

	redisClient := initializeRedis(cfg, log)
	defer redisClient.Close()

	rmqClient := initializeRabbitMQ(cfg, log)
	defer rmqClient.Close()

	wsHub := websocket.NewHub(rmqClient, log, queries, pool)
	go wsHub.Run()
	go func() {
		if err := wsHub.ListenToRabbitMQ(context.Background(), "codexec.results"); err != nil {
			log.Errorf("Failed to listen to RabbitMQ: %v", err)
		}
	}()

	authProvider := auth.NewProvider(
		cfg.API.PasswordSalt,
		[]byte(cfg.API.JwtSecret),
		cfg.API.JwtTTL,
		cfg.API.JwtRenewalThreshold,
	)

	router := api.NewRouter(queries, pool, log, authProvider, redisClient, wsHub)

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

	log.Info("Logger initialized successfully")

	return log, nil
}

func initializeQueries(cfg config.Config, log *logger.Logger) (*db.Queries, *pgxpool.Pool) {
	pool, err := db.NewPool(context.Background(), cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Info("Database initialized successfully")

	return db.New(pool), pool
}

func initializeRedis(cfg config.Config, log *logger.Logger) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:       addr,
			Password:   cfg.Redis.Password,
			DB:         cfg.Redis.DB,
			MaxRetries: cfg.Redis.MaxRetries,
		},
	)

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	log.Info("Redis initialized successfully")

	return redisClient
}

func initializeRabbitMQ(cfg config.Config, log *logger.Logger) *rabbitmq.Client {
	rmqClient, err := rabbitmq.NewClient(cfg.RabbitMQ, log)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	return rmqClient
}

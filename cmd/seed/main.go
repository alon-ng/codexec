package main

import (
	"codim/cmd/seed/config"
	"codim/pkg/api/auth"
	"codim/pkg/db"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	queries, pool := initializeQueries(cfg)
	defer pool.Close()

	authProvider := auth.NewProvider(
		cfg.API.PasswordSalt,
		[]byte(cfg.API.JwtSecret),
		cfg.API.JwtTTL,
		cfg.API.JwtRenewalThreshold,
	)

	seed(context.Background(), queries, authProvider)
}

func initializeQueries(cfg config.Config) (*db.Queries, *pgxpool.Pool) {
	pool, err := db.NewPool(context.Background(), cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database initialized successfully")

	return db.New(pool), pool
}

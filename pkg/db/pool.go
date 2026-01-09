package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(cfg.MaxConns)
	dbConfig.MinConns = int32(cfg.MinConns)
	dbConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	dbConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	return pgxpool.NewWithConfig(ctx, dbConfig)
}

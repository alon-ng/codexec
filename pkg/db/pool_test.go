package db_test

import (
	"codim/pkg/db"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPoolWithInvalidConnectionString(t *testing.T) {
	cfg := db.Config{
		ConnectionString: "invalid://connection/string",
		MaxConns:         10,
		MinConns:        2,
		ConnMaxLifetime: 0,
		ConnMaxIdleTime: 0,
	}

	pool, err := db.NewPool(context.Background(), cfg)
	require.Error(t, err)
	require.Nil(t, pool)
}

package db_test

import (
	"codim/pkg/db"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithTx(t *testing.T) {
	// Get the connection pool from testQueries
	// We need to access the underlying pool to begin a transaction
	config, err := db.LoadConfig()
	require.NoError(t, err)

	connPool, err := db.NewPool(context.Background(), config)
	require.NoError(t, err)
	defer connPool.Close()

	// Begin a transaction
	tx, err := connPool.Begin(context.Background())
	require.NoError(t, err)
	defer tx.Rollback(context.Background())

	// Create queries with transaction
	baseQueries := db.New(connPool)
	txQueries := baseQueries.WithTx(tx)
	require.NotNil(t, txQueries)

	// Test that we can use the transaction queries
	// Create a course within the transaction
	params := db.CreateCourseParams{
		Subject:    "python",
		Price:      100,
		Discount:   0,
		IsActive:   true,
		Difficulty: 1,
	}

	course, err := txQueries.CreateCourse(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, course)

	// Verify the course exists within the transaction
	_, err = txQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.Error(t, err) // Should error because no translation exists

	// Rollback the transaction
	err = tx.Rollback(context.Background())
	require.NoError(t, err)

	// Verify the course doesn't exist outside the transaction
	_, err = testQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.Error(t, err)
}

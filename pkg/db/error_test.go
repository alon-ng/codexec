package db_test

import (
	"codim/pkg/db"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDuplicateKeyError(t *testing.T) {
	// Test with duplicate key error
	err := errors.New("duplicate key value violates unique constraint (SQLSTATE 23505)")
	require.True(t, db.IsDuplicateKeyError(err))

	// Test with non-duplicate key error
	err = errors.New("some other error")
	require.False(t, db.IsDuplicateKeyError(err))

	// Test with empty error
	err = errors.New("")
	require.False(t, db.IsDuplicateKeyError(err))
}

func TestIsDuplicateKeyErrorWithConstraint(t *testing.T) {
	// Test with matching constraint
	err := errors.New(`duplicate key value violates unique constraint "uq_users_email" (SQLSTATE 23505)`)
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email"))

	// Test with different constraint
	err = errors.New(`duplicate key value violates unique constraint "uq_courses_subject" (SQLSTATE 23505)`)
	require.False(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email"))

	// Test with non-duplicate key error
	err = errors.New("some other error")
	require.False(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email"))
}

package db_test

import (
	"codim/pkg/db"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExerciseTypeScan(t *testing.T) {
	var et db.ExerciseType

	// Test scanning from []byte
	err := et.Scan([]byte("quiz"))
	require.NoError(t, err)
	require.Equal(t, db.ExerciseTypeQuiz, et)

	// Test scanning from string
	err = et.Scan("code")
	require.NoError(t, err)
	require.Equal(t, db.ExerciseTypeCode, et)

	// Test scanning from unsupported type
	err = et.Scan(123)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported scan type")
}

func TestNullExerciseTypeScan(t *testing.T) {
	var net db.NullExerciseType

	// Test scanning nil value
	err := net.Scan(nil)
	require.NoError(t, err)
	require.False(t, net.Valid)
	require.Equal(t, db.ExerciseType(""), net.ExerciseType)

	// Test scanning valid value
	err = net.Scan("quiz")
	require.NoError(t, err)
	require.True(t, net.Valid)
	require.Equal(t, db.ExerciseTypeQuiz, net.ExerciseType)
}

func TestNullExerciseTypeValue(t *testing.T) {
	var net db.NullExerciseType

	// Test Value with invalid (nil) value
	net.Valid = false
	val, err := net.Value()
	require.NoError(t, err)
	require.Nil(t, val)

	// Test Value with valid value
	net.Valid = true
	net.ExerciseType = db.ExerciseTypeQuiz
	val, err = net.Value()
	require.NoError(t, err)
	require.Equal(t, driver.Value("quiz"), val)
}

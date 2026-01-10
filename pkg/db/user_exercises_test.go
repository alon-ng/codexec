package db_test

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUserExercise(t *testing.T) db.UserExercise {
	user := createRandomUser(t)
	exercise := createRandomExercise(t)

	submission := json.RawMessage(`{"answer": "test"}`)
	params := db.CreateUserExerciseParams{
		UserUuid:     user.Uuid,
		ExerciseUuid: exercise.Uuid,
		Submission:   submission,
		Attempts:     1,
		CompletedAt:  nil,
	}

	userExercise, err := testQueries.CreateUserExercise(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, userExercise)

	require.Equal(t, params.UserUuid, userExercise.UserUuid)
	require.Equal(t, params.ExerciseUuid, userExercise.ExerciseUuid)
	require.Equal(t, params.Attempts, userExercise.Attempts)
	require.Nil(t, userExercise.CompletedAt)

	require.NotZero(t, userExercise.Uuid)
	require.NotZero(t, userExercise.StartedAt)
	require.Nil(t, userExercise.LastAccessedAt)

	return userExercise
}

func assertUserExerciseEqual(t *testing.T, expectedUserExercise db.UserExercise, gotUserExercise db.UserExercise) {
	assert.NotNil(t, gotUserExercise)

	require.Equal(t, expectedUserExercise.Uuid, gotUserExercise.Uuid)
	require.Equal(t, expectedUserExercise.UserUuid, gotUserExercise.UserUuid)
	require.Equal(t, expectedUserExercise.ExerciseUuid, gotUserExercise.ExerciseUuid)
	require.Equal(t, expectedUserExercise.Attempts, gotUserExercise.Attempts)
	require.Equal(t, expectedUserExercise.CompletedAt, gotUserExercise.CompletedAt)

	require.NotZero(t, gotUserExercise.StartedAt)
}

func TestCreateUserExercise(t *testing.T) {
	createRandomUserExercise(t)
}

func TestCreateUserExerciseWithCompletedAt(t *testing.T) {
	user := createRandomUser(t)
	exercise := createRandomExercise(t)
	now := time.Now()

	submission := json.RawMessage(`{"answer": "completed"}`)
	params := db.CreateUserExerciseParams{
		UserUuid:     user.Uuid,
		ExerciseUuid: exercise.Uuid,
		Submission:   submission,
		Attempts:     3,
		CompletedAt:  &now,
	}

	userExercise, err := testQueries.CreateUserExercise(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, userExercise)

	require.Equal(t, params.UserUuid, userExercise.UserUuid)
	require.Equal(t, params.ExerciseUuid, userExercise.ExerciseUuid)
	require.Equal(t, params.Attempts, userExercise.Attempts)
	require.NotNil(t, userExercise.CompletedAt)
}

func TestGetUserExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	gotUserExercise, err := testQueries.GetUserExercise(context.Background(), userExercise.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotUserExercise)

	assertUserExerciseEqual(t, userExercise, gotUserExercise)
}

func TestGetUserExerciseByUserAndExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	gotUserExercise, err := testQueries.GetUserExerciseByUserAndExercise(context.Background(), db.GetUserExerciseByUserAndExerciseParams{
		UserUuid:     userExercise.UserUuid,
		ExerciseUuid: userExercise.ExerciseUuid,
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotUserExercise)

	assertUserExerciseEqual(t, userExercise, gotUserExercise)
}

func TestUpdateUserExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)
	now := time.Now().UTC()

	updatedSubmission := json.RawMessage(`{"answer": "updated"}`)
	updateParams := db.UpdateUserExerciseParams{
		Uuid:         userExercise.Uuid,
		UserUuid:     userExercise.UserUuid,
		ExerciseUuid: userExercise.ExerciseUuid,
		Submission:   updatedSubmission,
		Attempts:     5,
		CompletedAt:  &now,
	}

	updatedUserExercise, err := testQueries.UpdateUserExercise(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUserExercise)

	require.Equal(t, updateParams.Attempts, updatedUserExercise.Attempts)
	require.Equal(t, updateParams.CompletedAt, updatedUserExercise.CompletedAt)
	require.Equal(t, userExercise.Uuid, updatedUserExercise.Uuid)
	require.Equal(t, userExercise.UserUuid, updatedUserExercise.UserUuid)
	require.Equal(t, userExercise.ExerciseUuid, updatedUserExercise.ExerciseUuid)
}

func TestDeleteUserExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	err := testQueries.DeleteUserExercise(context.Background(), userExercise.Uuid)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = testQueries.GetUserExercise(context.Background(), userExercise.Uuid)
	require.Error(t, err)
}

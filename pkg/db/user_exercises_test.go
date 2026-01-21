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
	course := createRandomCourse(t)
	lesson := createRandomLesson(t, &course)
	exercise := createRandomExercise(t, &lesson)

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
	course := createRandomCourse(t)
	lesson := createRandomLesson(t, &course)
	exercise := createRandomExercise(t, &lesson)
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

	gotUserExercise, err := testQueries.GetUserExercise(context.Background(), db.GetUserExerciseParams{
		UserUuid:     userExercise.UserUuid,
		ExerciseUuid: userExercise.ExerciseUuid,
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotUserExercise)

	assertUserExerciseEqual(t, userExercise, gotUserExercise)
}

func TestGetUserExerciseByUserAndExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	gotUserExercise, err := testQueries.GetUserExercise(context.Background(), db.GetUserExerciseParams{
		UserUuid:     userExercise.UserUuid,
		ExerciseUuid: userExercise.ExerciseUuid,
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotUserExercise)

	assertUserExerciseEqual(t, userExercise, gotUserExercise)
}

func TestUpdateUserExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	updatedSubmission := json.RawMessage(`{"answer": "updated"}`)
	userUuid := userExercise.UserUuid
	exerciseUuid := userExercise.ExerciseUuid
	updateParams := db.UpdateUserExerciseSubmissionParams{
		UserUuid:     userUuid,
		ExerciseUuid: exerciseUuid,
		Submission:   &updatedSubmission,
		Type:         db.ExerciseTypeQuiz,
	}

	updatedUserExercise, err := testQueries.UpdateUserExerciseSubmission(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUserExercise)

	require.Equal(t, *updateParams.Submission, updatedUserExercise.Submission)
	require.Equal(t, userExercise.Uuid, updatedUserExercise.Uuid)
	require.Equal(t, userExercise.UserUuid, updatedUserExercise.UserUuid)
	require.Equal(t, userExercise.ExerciseUuid, updatedUserExercise.ExerciseUuid)
}

func TestUpdateUserExerciseSubmissionWithAttempts(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	updatedSubmission := json.RawMessage(`{"answer": "updated"}`)
	userUuid := userExercise.UserUuid
	exerciseUuid := userExercise.ExerciseUuid
	updateParams := db.UpdateUserExerciseSubmissionWithAttemptsParams{
		UserUuid:     userUuid,
		ExerciseUuid: exerciseUuid,
		Submission:   &updatedSubmission,
	}

	updatedUserExercise, err := testQueries.UpdateUserExerciseSubmissionWithAttempts(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUserExercise)

	require.Equal(t, userExercise.Attempts+1, updatedUserExercise.Attempts)
	require.Equal(t, *updateParams.Submission, updatedUserExercise.Submission)
	require.Equal(t, userExercise.Uuid, updatedUserExercise.Uuid)
	require.Equal(t, userExercise.UserUuid, updatedUserExercise.UserUuid)
	require.Equal(t, userExercise.ExerciseUuid, updatedUserExercise.ExerciseUuid)
}

func TestCompleteUserExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	completeParams := db.CompleteUserExerciseParams{
		UserUuid:     userExercise.UserUuid,
		ExerciseUuid: userExercise.ExerciseUuid,
	}

	completeUserExercise, err := testQueries.CompleteUserExercise(context.Background(), completeParams)
	require.NoError(t, err)
	require.NotEmpty(t, completeUserExercise)

	require.NotNil(t, completeUserExercise.CompletedAt)
	require.Equal(t, userExercise.Uuid, completeUserExercise.Uuid)
	require.Equal(t, userExercise.UserUuid, completeUserExercise.UserUuid)
	require.Equal(t, userExercise.ExerciseUuid, completeUserExercise.ExerciseUuid)
}

func TestResetUserExercise(t *testing.T) {
	userExercise := createRandomUserExercise(t)

	resetUserExercise, err := testQueries.ResetUserExercise(context.Background(), db.ResetUserExerciseParams{
		UserUuid:     userExercise.UserUuid,
		ExerciseUuid: userExercise.ExerciseUuid,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resetUserExercise)

	require.Equal(t, userExercise.Uuid, resetUserExercise.Uuid)
	require.Equal(t, userExercise.UserUuid, resetUserExercise.UserUuid)
	require.Equal(t, userExercise.ExerciseUuid, resetUserExercise.ExerciseUuid)
	require.Equal(t, json.RawMessage(`{}`), resetUserExercise.Submission)
	require.Equal(t, int32(0), resetUserExercise.Attempts)
	require.Nil(t, resetUserExercise.CompletedAt)
}

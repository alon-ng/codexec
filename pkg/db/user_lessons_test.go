package db_test

import (
	"codim/pkg/db"
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUserLesson(t *testing.T) db.UserLesson {
	user := createRandomUser(t)
	lesson := createRandomLesson(t)

	params := db.CreateUserLessonParams{
		UserUuid:    user.Uuid,
		LessonUuid:  lesson.Uuid,
		CompletedAt: nil,
	}

	userLesson, err := testQueries.CreateUserLesson(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, userLesson)

	require.Equal(t, params.UserUuid, userLesson.UserUuid)
	require.Equal(t, params.LessonUuid, userLesson.LessonUuid)
	require.Nil(t, userLesson.CompletedAt)

	require.NotZero(t, userLesson.Uuid)
	require.NotZero(t, userLesson.StartedAt)
	require.Nil(t, userLesson.LastAccessedAt)

	return userLesson
}

func assertUserLessonEqual(t *testing.T, expectedUserLesson db.UserLesson, gotUserLesson db.UserLesson) {
	assert.NotNil(t, gotUserLesson)

	require.Equal(t, expectedUserLesson.Uuid, gotUserLesson.Uuid)
	require.Equal(t, expectedUserLesson.UserUuid, gotUserLesson.UserUuid)
	require.Equal(t, expectedUserLesson.LessonUuid, gotUserLesson.LessonUuid)
	require.Equal(t, expectedUserLesson.CompletedAt, gotUserLesson.CompletedAt)

	require.NotZero(t, gotUserLesson.StartedAt)
}

func TestCreateUserLesson(t *testing.T) {
	createRandomUserLesson(t)
}

func TestCreateUserLessonWithCompletedAt(t *testing.T) {
	user := createRandomUser(t)
	lesson := createRandomLesson(t)
	now := time.Now()

	params := db.CreateUserLessonParams{
		UserUuid:    user.Uuid,
		LessonUuid:  lesson.Uuid,
		CompletedAt: &now,
	}

	userLesson, err := testQueries.CreateUserLesson(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, userLesson)

	require.Equal(t, params.UserUuid, userLesson.UserUuid)
	require.Equal(t, params.LessonUuid, userLesson.LessonUuid)
	require.NotNil(t, userLesson.CompletedAt)
}

func TestGetUserLesson(t *testing.T) {
	userLesson := createRandomUserLesson(t)

	gotUserLesson, err := testQueries.GetUserLesson(context.Background(), userLesson.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotUserLesson)

	assertUserLessonEqual(t, userLesson, gotUserLesson)
}

func TestGetUserLessonByUserAndLesson(t *testing.T) {
	userLesson := createRandomUserLesson(t)

	gotUserLesson, err := testQueries.GetUserLessonByUserAndLesson(context.Background(), db.GetUserLessonByUserAndLessonParams{
		UserUuid:   userLesson.UserUuid,
		LessonUuid: userLesson.LessonUuid,
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotUserLesson)

	assertUserLessonEqual(t, userLesson, gotUserLesson)
}

func TestUpdateUserLesson(t *testing.T) {
	userLesson := createRandomUserLesson(t)
	now := time.Now().UTC()

	updateParams := db.UpdateUserLessonParams{
		Uuid:        userLesson.Uuid,
		UserUuid:    userLesson.UserUuid,
		LessonUuid:  userLesson.LessonUuid,
		CompletedAt: &now,
	}

	updatedUserLesson, err := testQueries.UpdateUserLesson(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUserLesson)

	require.Equal(t, updateParams.CompletedAt, updatedUserLesson.CompletedAt)
	require.Equal(t, userLesson.Uuid, updatedUserLesson.Uuid)
	require.Equal(t, userLesson.UserUuid, updatedUserLesson.UserUuid)
	require.Equal(t, userLesson.LessonUuid, updatedUserLesson.LessonUuid)
}

func TestDeleteUserLesson(t *testing.T) {
	userLesson := createRandomUserLesson(t)

	err := testQueries.DeleteUserLesson(context.Background(), userLesson.Uuid)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = testQueries.GetUserLesson(context.Background(), userLesson.Uuid)
	require.Error(t, err)
}

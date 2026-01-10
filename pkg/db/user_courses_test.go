package db_test

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomLessonWithCourse(t *testing.T, course db.Course) db.Lesson {
	rnd := getRandomInt()
	params := db.CreateLessonParams{
		CourseUuid:  course.Uuid,
		Name:        fmt.Sprintf("Test Lesson %d", rnd),
		Description: fmt.Sprintf("Test Description %d", rnd),
		OrderIndex:  1,
		IsPublic:    true,
	}

	lesson, err := testQueries.CreateLesson(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, lesson)

	return lesson
}

func createRandomExerciseWithLesson(t *testing.T, lesson db.Lesson) db.Exercise {
	rnd := getRandomInt()
	testData := json.RawMessage(fmt.Sprintf(`{"answer": "Test Answer", "question": "Test Question %d"}`, rnd))
	params := db.CreateExerciseParams{
		LessonUuid:  lesson.Uuid,
		Name:        fmt.Sprintf("Test Exercise %d", rnd),
		Description: fmt.Sprintf("Test Description %d", rnd),
		OrderIndex:  1,
		Reward:      10,
		Type:        db.ExerciseTypeQuiz,
		Data:        testData,
	}

	exercise, err := testQueries.CreateExercise(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, exercise)

	return exercise
}

func createRandomUserCourse(t *testing.T) db.UserCourse {
	user := createRandomUser(t)
	course := createRandomCourse(t)

	params := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: nil,
	}

	userCourse, err := testQueries.CreateUserCourse(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, userCourse)

	require.Equal(t, params.UserUuid, userCourse.UserUuid)
	require.Equal(t, params.CourseUuid, userCourse.CourseUuid)
	require.Nil(t, userCourse.CompletedAt)

	require.NotZero(t, userCourse.Uuid)
	require.NotZero(t, userCourse.StartedAt)
	require.Nil(t, userCourse.LastAccessedAt)

	return userCourse
}

func assertUserCourseEqual(t *testing.T, expectedUserCourse db.UserCourse, gotUserCourse db.UserCourse) {
	assert.NotNil(t, gotUserCourse)

	require.Equal(t, expectedUserCourse.Uuid, gotUserCourse.Uuid)
	require.Equal(t, expectedUserCourse.UserUuid, gotUserCourse.UserUuid)
	require.Equal(t, expectedUserCourse.CourseUuid, gotUserCourse.CourseUuid)
	require.Equal(t, expectedUserCourse.CompletedAt, gotUserCourse.CompletedAt)

	require.NotZero(t, gotUserCourse.StartedAt)
}

func TestCreateUserCourse(t *testing.T) {
	createRandomUserCourse(t)
}

func TestCreateUserCourseWithCompletedAt(t *testing.T) {
	user := createRandomUser(t)
	course := createRandomCourse(t)
	now := time.Now()

	params := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: &now,
	}

	userCourse, err := testQueries.CreateUserCourse(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, userCourse)

	require.Equal(t, params.UserUuid, userCourse.UserUuid)
	require.Equal(t, params.CourseUuid, userCourse.CourseUuid)
	require.NotNil(t, userCourse.CompletedAt)
}

func TestUpdateUserCourse(t *testing.T) {
	userCourse := createRandomUserCourse(t)
	now := time.Now().UTC()

	updateParams := db.UpdateUserCourseParams{
		Uuid:        userCourse.Uuid,
		UserUuid:    userCourse.UserUuid,
		CourseUuid:  userCourse.CourseUuid,
		CompletedAt: &now,
	}

	updatedUserCourse, err := testQueries.UpdateUserCourse(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUserCourse)

	require.Equal(t, updateParams.CompletedAt, updatedUserCourse.CompletedAt)
	require.Equal(t, userCourse.Uuid, updatedUserCourse.Uuid)
	require.Equal(t, userCourse.UserUuid, updatedUserCourse.UserUuid)
	require.Equal(t, userCourse.CourseUuid, updatedUserCourse.CourseUuid)
}

func TestDeleteUserCourse(t *testing.T) {
	userCourse := createRandomUserCourse(t)

	err := testQueries.DeleteUserCourse(context.Background(), userCourse.Uuid)
	require.NoError(t, err)
}

func TestGetUserCourseFull(t *testing.T) {
	user := createRandomUser(t)
	course := createRandomCourse(t)

	// Create user course
	userCourseParams := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: nil,
	}
	_, err := testQueries.CreateUserCourse(context.Background(), userCourseParams)
	require.NoError(t, err)

	// Create lesson
	lesson := createRandomLessonWithCourse(t, course)

	// Create user lesson
	userLessonParams := db.CreateUserLessonParams{
		UserUuid:    user.Uuid,
		LessonUuid:  lesson.Uuid,
		CompletedAt: nil,
	}
	_, err = testQueries.CreateUserLesson(context.Background(), userLessonParams)
	require.NoError(t, err)

	// Create exercise
	exercise := createRandomExerciseWithLesson(t, lesson)

	// Create user exercise
	submission := json.RawMessage(`{"answer": "test"}`)
	userExerciseParams := db.CreateUserExerciseParams{
		UserUuid:     user.Uuid,
		ExerciseUuid: exercise.Uuid,
		Submission:   submission,
		Attempts:     1,
		CompletedAt:  nil,
	}
	_, err = testQueries.CreateUserExercise(context.Background(), userExerciseParams)
	require.NoError(t, err)

	// Test GetUserCourseFull
	result, err := testQueries.GetUserCourseFull(context.Background(), user.Uuid, course.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, course.Uuid, result.CourseUuid)
	require.NotNil(t, result.StartedAt)
	require.Equal(t, false, result.IsCompleted)
	require.Nil(t, result.CompletedAt)

	require.Len(t, result.Lessons, 1)
	require.Equal(t, lesson.Uuid, result.Lessons[0].LessonUuid)
	require.NotNil(t, result.Lessons[0].StartedAt)
	require.Equal(t, false, result.Lessons[0].IsCompleted)

	require.Len(t, result.Lessons[0].Exercises, 1)
	require.Equal(t, exercise.Uuid, result.Lessons[0].Exercises[0].ExerciseUuid)
	require.NotNil(t, result.Lessons[0].Exercises[0].StartedAt)
	require.Equal(t, false, result.Lessons[0].Exercises[0].IsCompleted)
}

func TestGetUserCourseFullWithCompletedItems(t *testing.T) {
	user := createRandomUser(t)
	course := createRandomCourse(t)
	now := time.Now()

	// Create completed user course
	userCourseParams := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: &now,
	}
	_, err := testQueries.CreateUserCourse(context.Background(), userCourseParams)
	require.NoError(t, err)

	// Create lesson
	lesson := createRandomLessonWithCourse(t, course)

	// Create completed user lesson
	userLessonParams := db.CreateUserLessonParams{
		UserUuid:    user.Uuid,
		LessonUuid:  lesson.Uuid,
		CompletedAt: &now,
	}
	_, err = testQueries.CreateUserLesson(context.Background(), userLessonParams)
	require.NoError(t, err)

	// Create exercise
	exercise := createRandomExerciseWithLesson(t, lesson)

	// Create completed user exercise
	submission := json.RawMessage(`{"answer": "test"}`)
	userExerciseParams := db.CreateUserExerciseParams{
		UserUuid:     user.Uuid,
		ExerciseUuid: exercise.Uuid,
		Submission:   submission,
		Attempts:     1,
		CompletedAt:  &now,
	}
	_, err = testQueries.CreateUserExercise(context.Background(), userExerciseParams)
	require.NoError(t, err)

	// Test GetUserCourseFull
	result, err := testQueries.GetUserCourseFull(context.Background(), user.Uuid, course.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, true, result.IsCompleted)
	require.NotNil(t, result.CompletedAt)

	require.Len(t, result.Lessons, 1)
	require.Equal(t, true, result.Lessons[0].IsCompleted)
	require.NotNil(t, result.Lessons[0].CompletedAt)

	require.Len(t, result.Lessons[0].Exercises, 1)
	require.Equal(t, true, result.Lessons[0].Exercises[0].IsCompleted)
	require.NotNil(t, result.Lessons[0].Exercises[0].CompletedAt)
}

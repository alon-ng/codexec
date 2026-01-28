package db_test

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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
		UserUuid:    &userCourse.UserUuid,
		CourseUuid:  &userCourse.CourseUuid,
		CompletedAt: &now,
	}

	updatedUserCourse, err := testQueries.UpdateUserCourse(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUserCourse)

	require.Equal(t, updateParams.Uuid, updatedUserCourse.Uuid)
	require.Equal(t, *updateParams.UserUuid, updatedUserCourse.UserUuid)
	require.Equal(t, *updateParams.CourseUuid, updatedUserCourse.CourseUuid)
	require.Equal(t, updateParams.CompletedAt, updatedUserCourse.CompletedAt)
}

func TestDeleteUserCourse(t *testing.T) {
	userCourse := createRandomUserCourse(t)

	err := testQueries.DeleteUserCourse(context.Background(), userCourse.Uuid)
	require.NoError(t, err)
}

func TestGetUserCourseFull(t *testing.T) {
	user := createRandomUser(t)
	course := createRandomCourse(t)
	lesson := createRandomLesson(t, &course)
	exercise := createRandomExercise(t, &lesson)

	// Create user course
	userCourseParams := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: nil,
	}
	_, err := testQueries.CreateUserCourse(context.Background(), userCourseParams)
	require.NoError(t, err)

	// Create user lesson
	userLessonParams := db.CreateUserLessonParams{
		UserUuid:    user.Uuid,
		LessonUuid:  lesson.Uuid,
		CompletedAt: nil,
	}
	_, err = testQueries.CreateUserLesson(context.Background(), userLessonParams)
	require.NoError(t, err)

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
	lesson := createRandomLesson(t, &course)
	exercise := createRandomExercise(t, &lesson)
	now := time.Now()

	// Create completed user course
	userCourseParams := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: &now,
	}
	_, err := testQueries.CreateUserCourse(context.Background(), userCourseParams)
	require.NoError(t, err)

	// Create completed user lesson
	userLessonParams := db.CreateUserLessonParams{
		UserUuid:    user.Uuid,
		LessonUuid:  lesson.Uuid,
		CompletedAt: &now,
	}
	_, err = testQueries.CreateUserLesson(context.Background(), userLessonParams)
	require.NoError(t, err)

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

func TestGetUserCourseFullWithEmptyResult(t *testing.T) {
	user := createRandomUser(t)
	nonExistentUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	result, err := testQueries.GetUserCourseFull(context.Background(), user.Uuid, nonExistentUUID)
	require.NoError(t, err)
	require.Equal(t, uuid.Nil, result.CourseUuid)
	require.Empty(t, result.Lessons)
}

func TestInitUserCourse(t *testing.T) {
	user := createRandomUser(t)
	course := createRandomCourse(t)
	lesson := createRandomLesson(t, &course)
	_ = createRandomExercise(t, &lesson)

	result, err := testQueries.InitUserCourse(context.Background(), db.InitUserCourseParams{
		UserUuid:   user.Uuid,
		CourseUuid: course.Uuid,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.Uuid, result.UserUuid)
	require.Equal(t, course.Uuid, result.CourseUuid)
	require.NotZero(t, result.Uuid)
	require.NotZero(t, result.StartedAt)
	require.Nil(t, result.CompletedAt)
}

func TestListUserCoursesWithProgress(t *testing.T) {
	userCourse := createRandomUserCourse(t)

	params := db.ListUserCoursesWithProgressParams{
		UserUuid:   userCourse.UserUuid,
		Language:   "en",
		Limit:      10,
		Offset:     0,
		CourseUuid: nil,
		Subject:    nil,
		IsActive:   nil,
	}

	courses, err := testQueries.ListUserCoursesWithProgress(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(courses), 1)

	var foundCourse *db.ListUserCoursesWithProgressRow
	for i := range courses {
		if courses[i].CourseUuid == userCourse.CourseUuid {
			foundCourse = &courses[i]
			break
		}
	}

	require.NotNil(t, foundCourse)
	courseWithProgress := foundCourse.ToUserCourseWithProgress()
	require.Equal(t, userCourse.CourseUuid, courseWithProgress.Uuid)
	require.NotZero(t, courseWithProgress.UserCourseStartedAt)
}

func TestListUserCoursesWithProgressWithFilters(t *testing.T) {
	user := createRandomUser(t)
	course := createRandomCourse(t)
	userCourseParams := db.CreateUserCourseParams{
		UserUuid:    user.Uuid,
		CourseUuid:  course.Uuid,
		CompletedAt: nil,
	}
	_, err := testQueries.CreateUserCourse(context.Background(), userCourseParams)
	require.NoError(t, err)

	subject := course.Subject
	isActive := course.IsActive
	params := db.ListUserCoursesWithProgressParams{
		UserUuid:   user.Uuid,
		Language:   "en",
		Limit:      10,
		Offset:     0,
		CourseUuid: nil,
		Subject:    &subject,
		IsActive:   &isActive,
	}

	courses, err := testQueries.ListUserCoursesWithProgress(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(courses), 1)

	for _, c := range courses {
		require.Equal(t, subject, c.CourseSubject)
		require.Equal(t, isActive, c.CourseIsActive)
	}
}

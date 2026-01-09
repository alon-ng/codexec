package db_test

import (
	"codim/pkg/db"
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomLesson(t *testing.T) db.Lesson {
	course := createRandomCourse(t)

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

	require.Equal(t, params.CourseUuid, lesson.CourseUuid)
	require.Equal(t, params.Name, lesson.Name)
	require.Equal(t, params.Description, lesson.Description)
	require.Equal(t, params.OrderIndex, lesson.OrderIndex)
	require.Equal(t, params.IsPublic, lesson.IsPublic)

	require.NotZero(t, lesson.Uuid)
	require.NotZero(t, lesson.CreatedAt)
	require.NotZero(t, lesson.ModifiedAt)
	require.Nil(t, lesson.DeletedAt)

	return lesson
}

func assertLessonEqual(t *testing.T, expectedLesson db.Lesson, gotLesson db.Lesson) {
	assert.NotNil(t, gotLesson)

	require.Equal(t, expectedLesson.Uuid, gotLesson.Uuid)
	require.Equal(t, expectedLesson.CourseUuid, gotLesson.CourseUuid)
	require.Equal(t, expectedLesson.Name, gotLesson.Name)
	require.Equal(t, expectedLesson.Description, gotLesson.Description)
	require.Equal(t, expectedLesson.OrderIndex, gotLesson.OrderIndex)
	require.Equal(t, expectedLesson.IsPublic, gotLesson.IsPublic)

	require.NotZero(t, gotLesson.CreatedAt)
	require.NotZero(t, gotLesson.ModifiedAt)
	require.Nil(t, gotLesson.DeletedAt)
}

func TestCreateLesson(t *testing.T) {
	createRandomLesson(t)
}

func TestGetLesson(t *testing.T) {
	lesson := createRandomLesson(t)

	gotLesson, err := testQueries.GetLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotLesson)

	assertLessonEqual(t, lesson, gotLesson)
}

func TestUpdateLesson(t *testing.T) {
	lesson := createRandomLesson(t)
	course := createRandomCourse(t)

	rnd := getRandomInt()
	updateParams := db.UpdateLessonParams{
		Uuid:        lesson.Uuid,
		CourseUuid:  course.Uuid,
		Name:        fmt.Sprintf("Updated Test Lesson %d", rnd),
		Description: fmt.Sprintf("Updated Test Description %d", rnd),
		OrderIndex:  2,
		IsPublic:    false,
	}

	updatedLesson, err := testQueries.UpdateLesson(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedLesson)

	require.Equal(t, updateParams.CourseUuid, updatedLesson.CourseUuid)
	require.Equal(t, updateParams.Name, updatedLesson.Name)
	require.Equal(t, updateParams.Description, updatedLesson.Description)
	require.Equal(t, updateParams.OrderIndex, updatedLesson.OrderIndex)
	require.Equal(t, updateParams.IsPublic, updatedLesson.IsPublic)

	require.NotZero(t, updatedLesson.ModifiedAt)
	require.Nil(t, updatedLesson.DeletedAt)
}

func TestDeleteLesson(t *testing.T) {
	lesson := createRandomLesson(t)

	err := testQueries.DeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err := testQueries.GetLesson(context.Background(), lesson.Uuid)
	require.Error(t, err)
	require.Empty(t, gotLesson)
}

func TestHardDeleteLesson(t *testing.T) {
	lesson := createRandomLesson(t)

	err := testQueries.HardDeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err := testQueries.GetLesson(context.Background(), lesson.Uuid)
	require.Error(t, err)
	require.Empty(t, gotLesson)
}

func TestUndeleteLesson(t *testing.T) {
	lesson := createRandomLesson(t)

	err := testQueries.DeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err := testQueries.GetLesson(context.Background(), lesson.Uuid)
	require.Error(t, err)
	require.Empty(t, gotLesson)

	err = testQueries.UndeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err = testQueries.GetLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotLesson)

	assertLessonEqual(t, lesson, gotLesson)
}

func TestCreateLessonConflict(t *testing.T) {
	lesson := createRandomLesson(t)

	params := db.CreateLessonParams{
		CourseUuid:  lesson.CourseUuid,
		Name:        lesson.Name,
		Description: lesson.Description,
		OrderIndex:  lesson.OrderIndex,
		IsPublic:    lesson.IsPublic,
	}

	_, err := testQueries.CreateLesson(context.Background(), params)
	require.Error(t, err)
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "lessons_name_key"))
}

func TestCountLessons(t *testing.T) {
	initialCount, err := testQueries.CountLessons(context.Background())
	require.NoError(t, err)

	lesson1 := createRandomLesson(t)
	_ = createRandomLesson(t)

	count, err := testQueries.CountLessons(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+2, count)

	err = testQueries.DeleteLesson(context.Background(), lesson1.Uuid)
	require.NoError(t, err)

	count, err = testQueries.CountLessons(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
}

func TestCountLessonsByCourse(t *testing.T) {
	course := createRandomCourse(t)

	count, err := testQueries.CountLessonsByCourse(context.Background(), course.Uuid)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	lesson1 := createRandomLesson(t)
	lesson2 := createRandomLesson(t)

	count, err = testQueries.CountLessonsByCourse(context.Background(), lesson1.CourseUuid)
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(1))

	if lesson1.CourseUuid == lesson2.CourseUuid {
		count, err = testQueries.CountLessonsByCourse(context.Background(), lesson1.CourseUuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, int64(2))
	}
}

func TestListLessons(t *testing.T) {
	lesson1 := createRandomLesson(t)
	lesson2 := createRandomLesson(t)

	params := db.ListLessonsParams{
		Limit:      10,
		Offset:     0,
		CourseUuid: nil,
	}

	lessons, err := testQueries.ListLessons(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(lessons), 2)

	var foundLesson1, foundLesson2 bool
	for _, lesson := range lessons {
		if lesson.Uuid == lesson1.Uuid {
			foundLesson1 = true
			assertLessonEqual(t, lesson1, lesson)
		}
		if lesson.Uuid == lesson2.Uuid {
			foundLesson2 = true
			assertLessonEqual(t, lesson2, lesson)
		}
	}
	require.True(t, foundLesson1)
	require.True(t, foundLesson2)
}

func TestListLessonsWithCourseFilter(t *testing.T) {
	course := createRandomCourse(t)

	rnd := getRandomInt()
	lesson1Params := db.CreateLessonParams{
		CourseUuid:  course.Uuid,
		Name:        fmt.Sprintf("Test Lesson 1 %d", rnd),
		Description: fmt.Sprintf("Test Description 1 %d", rnd),
		OrderIndex:  1,
		IsPublic:    true,
	}
	lesson1, err := testQueries.CreateLesson(context.Background(), lesson1Params)
	require.NoError(t, err)

	lesson2Params := db.CreateLessonParams{
		CourseUuid:  course.Uuid,
		Name:        fmt.Sprintf("Test Lesson 2 %d", rnd),
		Description: fmt.Sprintf("Test Description 2 %d", rnd),
		OrderIndex:  2,
		IsPublic:    false,
	}
	lesson2, err := testQueries.CreateLesson(context.Background(), lesson2Params)
	require.NoError(t, err)

	params := db.ListLessonsParams{
		Limit:      10,
		Offset:     0,
		CourseUuid: &course.Uuid,
	}

	lessons, err := testQueries.ListLessons(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(lessons), 2)

	var foundLesson1, foundLesson2 bool
	for _, lesson := range lessons {
		require.Equal(t, course.Uuid, lesson.CourseUuid)
		if lesson.Uuid == lesson1.Uuid {
			foundLesson1 = true
			assertLessonEqual(t, lesson1, lesson)
		}
		if lesson.Uuid == lesson2.Uuid {
			foundLesson2 = true
			assertLessonEqual(t, lesson2, lesson)
		}
	}
	require.True(t, foundLesson1)
	require.True(t, foundLesson2)
}

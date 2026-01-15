package db_test

import (
	"codim/pkg/db"
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomLesson(t *testing.T, course *db.CourseWithTranslation) db.LessonWithTranslation {
	if course == nil {
		c := createRandomCourse(t)
		course = &c
	}

	rnd := getRandomInt()
	params := db.CreateLessonParams{
		CourseUuid: course.Uuid,
		OrderIndex: 1,
		IsPublic:   true,
	}

	lesson, err := testQueries.CreateLesson(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, lesson)

	require.Equal(t, params.CourseUuid, lesson.CourseUuid)
	require.Equal(t, params.OrderIndex, lesson.OrderIndex)
	require.Equal(t, params.IsPublic, lesson.IsPublic)

	require.NotZero(t, lesson.Uuid)
	require.NotZero(t, lesson.CreatedAt)
	require.NotZero(t, lesson.ModifiedAt)
	require.Nil(t, lesson.DeletedAt)

	translationParams := db.CreateLessonTranslationParams{
		LessonUuid:  lesson.Uuid,
		Language:    "en",
		Name:        fmt.Sprintf("Test Lesson %d", rnd),
		Description: fmt.Sprintf("Test Description %d", rnd),
	}
	lessonTranslation, err := testQueries.CreateLessonTranslation(context.Background(), translationParams)
	require.NoError(t, err)
	require.NotEmpty(t, lessonTranslation)

	require.Equal(t, translationParams.LessonUuid, lessonTranslation.LessonUuid)
	require.Equal(t, translationParams.Language, lessonTranslation.Language)
	require.Equal(t, translationParams.Name, lessonTranslation.Name)
	require.Equal(t, translationParams.Description, lessonTranslation.Description)

	return db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: lessonTranslation,
	}
}

func assertLessonWithTranslationEqual(t *testing.T, expectedLesson db.LessonWithTranslation, gotLesson db.LessonWithTranslation) {
	getLessonRow := db.GetLessonRow{
		Uuid:        gotLesson.Uuid,
		CreatedAt:   gotLesson.CreatedAt,
		ModifiedAt:  gotLesson.ModifiedAt,
		DeletedAt:   gotLesson.DeletedAt,
		CourseUuid:  gotLesson.CourseUuid,
		OrderIndex:  gotLesson.OrderIndex,
		IsPublic:    gotLesson.IsPublic,
		LessonUuid:  gotLesson.Translation.LessonUuid,
		Language:    gotLesson.Translation.Language,
		Name:        gotLesson.Translation.Name,
		Description: gotLesson.Translation.Description,
	}
	assertLessonEqual(t, expectedLesson, getLessonRow)
}

func assertLessonListEqual(t *testing.T, expectedLesson db.LessonWithTranslation, gotLesson db.ListLessonsRow) {
	getLessonRow := db.GetLessonRow{
		Uuid:        gotLesson.Uuid,
		CreatedAt:   gotLesson.CreatedAt,
		ModifiedAt:  gotLesson.ModifiedAt,
		DeletedAt:   gotLesson.DeletedAt,
		CourseUuid:  gotLesson.CourseUuid,
		OrderIndex:  gotLesson.OrderIndex,
		IsPublic:    gotLesson.IsPublic,
		LessonUuid:  gotLesson.LessonUuid,
		Language:    gotLesson.Language,
		Name:        gotLesson.Name,
		Description: gotLesson.Description,
	}
	assertLessonEqual(t, expectedLesson, getLessonRow)
}

func assertLessonEqual(t *testing.T, expectedLesson db.LessonWithTranslation, gotLesson db.GetLessonRow) {
	assert.NotNil(t, gotLesson)

	require.Equal(t, expectedLesson.Uuid, gotLesson.Uuid)
	require.Equal(t, expectedLesson.CourseUuid, gotLesson.CourseUuid)
	require.Equal(t, expectedLesson.OrderIndex, gotLesson.OrderIndex)
	require.Equal(t, expectedLesson.IsPublic, gotLesson.IsPublic)

	require.Equal(t, expectedLesson.Translation.Name, gotLesson.Name)
	require.Equal(t, expectedLesson.Translation.Description, gotLesson.Description)
	require.Equal(t, expectedLesson.Translation.Language, gotLesson.Language)
	require.Equal(t, expectedLesson.Translation.LessonUuid, gotLesson.LessonUuid)

	require.NotZero(t, gotLesson.CreatedAt)
	require.NotZero(t, gotLesson.ModifiedAt)
	require.Nil(t, gotLesson.DeletedAt)
}

func TestCreateLesson(t *testing.T) {
	createRandomLesson(t, nil)
}

func TestGetLesson(t *testing.T) {
	lesson := createRandomLesson(t, nil)

	gotLesson, err := testQueries.GetLesson(context.Background(), db.GetLessonParams{
		Uuid:     lesson.Lesson.Uuid,
		Language: "en",
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotLesson)

	assertLessonEqual(t, lesson, gotLesson)
}

func TestUpdateLesson(t *testing.T) {
	lesson := createRandomLesson(t, nil)
	course := createRandomCourse(t)

	rnd := getRandomInt()
	updateParams := db.UpdateLessonParams{
		Uuid:       lesson.Uuid,
		CourseUuid: course.Uuid,
		OrderIndex: 2,
		IsPublic:   false,
	}

	updatedLesson, err := testQueries.UpdateLesson(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedLesson)

	require.Equal(t, updateParams.CourseUuid, updatedLesson.CourseUuid)
	require.Equal(t, updateParams.OrderIndex, updatedLesson.OrderIndex)
	require.Equal(t, updateParams.IsPublic, updatedLesson.IsPublic)

	require.NotZero(t, updatedLesson.ModifiedAt)
	require.Nil(t, updatedLesson.DeletedAt)

	updateTranslationParams := db.UpdateLessonTranslationParams{
		Uuid:        lesson.Uuid,
		Language:    "en",
		Name:        fmt.Sprintf("Updated Test Lesson %d", rnd),
		Description: fmt.Sprintf("Updated Test Description %d", rnd),
	}
	updateTranslation, err := testQueries.UpdateLessonTranslation(context.Background(), updateTranslationParams)
	require.NoError(t, err)
	require.NotEmpty(t, updateTranslation)

	require.Equal(t, updateTranslationParams.Language, updateTranslation.Language)
	require.Equal(t, updateTranslationParams.Name, updateTranslation.Name)
	require.Equal(t, updateTranslationParams.Description, updateTranslation.Description)
}

func TestDeleteLesson(t *testing.T) {
	lesson := createRandomLesson(t, nil)

	err := testQueries.DeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err := testQueries.GetLesson(context.Background(), db.GetLessonParams{
		Uuid:     lesson.Lesson.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotLesson)
}

func TestHardDeleteLesson(t *testing.T) {
	lesson := createRandomLesson(t, nil)

	err := testQueries.HardDeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err := testQueries.GetLesson(context.Background(), db.GetLessonParams{
		Uuid:     lesson.Lesson.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotLesson)
}

func TestUndeleteLesson(t *testing.T) {
	lesson := createRandomLesson(t, nil)

	err := testQueries.DeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err := testQueries.GetLesson(context.Background(), db.GetLessonParams{
		Uuid:     lesson.Lesson.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotLesson)

	err = testQueries.UndeleteLesson(context.Background(), lesson.Uuid)
	require.NoError(t, err)

	gotLesson, err = testQueries.GetLesson(context.Background(), db.GetLessonParams{
		Uuid:     lesson.Lesson.Uuid,
		Language: "en",
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotLesson)

	assertLessonEqual(t, lesson, gotLesson)
}

func TestCountLessons(t *testing.T) {
	initialCount, err := testQueries.CountLessons(context.Background())
	require.NoError(t, err)

	lesson1 := createRandomLesson(t, nil)
	_ = createRandomLesson(t, nil)

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

	lesson1 := createRandomLesson(t, nil)
	lesson2 := createRandomLesson(t, nil)

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
	lesson1 := createRandomLesson(t, nil)
	lesson2 := createRandomLesson(t, nil)

	params := db.ListLessonsParams{
		Limit:      10,
		Offset:     0,
		Language:   "en",
		CourseUuid: nil,
	}

	lessons, err := testQueries.ListLessons(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(lessons), 2)

	var foundLesson1, foundLesson2 bool
	for _, lesson := range lessons {
		if lesson.Uuid == lesson1.Uuid {
			foundLesson1 = true
			assertLessonListEqual(t, lesson1, lesson)
		}
		if lesson.Uuid == lesson2.Uuid {
			foundLesson2 = true
			assertLessonListEqual(t, lesson2, lesson)
		}
	}
	require.True(t, foundLesson1)
	require.True(t, foundLesson2)
}

func TestListLessonsWithCourseFilter(t *testing.T) {
	course := createRandomCourse(t)

	lesson1 := createRandomLesson(t, &course)
	lesson2 := createRandomLesson(t, &course)

	params := db.ListLessonsParams{
		Limit:      10,
		Offset:     0,
		Language:   "en",
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
			assertLessonListEqual(t, lesson1, lesson)
		}
		if lesson.Uuid == lesson2.Uuid {
			foundLesson2 = true
			assertLessonListEqual(t, lesson2, lesson)
		}
	}
	require.True(t, foundLesson1)
	require.True(t, foundLesson2)
}

func TestCreateLessonTranslationWithConflict(t *testing.T) {
	lesson := createRandomLesson(t, nil)
	_, err := testQueries.CreateLessonTranslation(context.Background(), db.CreateLessonTranslationParams{
		LessonUuid:  lesson.Uuid,
		Language:    "en",
		Name:        "Test Lesson",
		Description: "Test Description",
	})
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_lesson_translations_lesson_language"))
}

package db_test

import (
	"codim/pkg/db"
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomCourse(t *testing.T) db.CourseWithTranslation {
	rnd := getRandomInt()
	params := db.CreateCourseParams{
		Subject:    "python",
		Price:      100,
		Discount:   0,
		IsActive:   true,
		Difficulty: 1,
	}

	course, err := testQueries.CreateCourse(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, course)

	require.Equal(t, params.Subject, course.Subject)
	require.Equal(t, params.Price, course.Price)
	require.Equal(t, params.Discount, course.Discount)
	require.Equal(t, params.IsActive, course.IsActive)
	require.Equal(t, params.Difficulty, course.Difficulty)

	require.NotZero(t, course.Uuid)
	require.NotZero(t, course.CreatedAt)
	require.NotZero(t, course.ModifiedAt)
	require.Nil(t, course.DeletedAt)

	courseTranslation, err := testQueries.CreateCourseTranslation(context.Background(), db.CreateCourseTranslationParams{
		CourseUuid:  course.Uuid,
		Language:    "en",
		Name:        fmt.Sprintf("Test Course %d", rnd),
		Description: fmt.Sprintf("Test Description %d", rnd),
		Bullets:     "Test Bullets 1\nTest Bullets 2\nTest Bullets 3",
	})
	require.NoError(t, err)
	require.NotEmpty(t, courseTranslation)

	require.Equal(t, course.Uuid, courseTranslation.CourseUuid)
	require.Equal(t, "en", courseTranslation.Language)
	require.Equal(t, fmt.Sprintf("Test Course %d", rnd), courseTranslation.Name)
	require.Equal(t, fmt.Sprintf("Test Description %d", rnd), courseTranslation.Description)
	require.Equal(t, "Test Bullets 1\nTest Bullets 2\nTest Bullets 3", courseTranslation.Bullets)

	return db.CourseWithTranslation{
		Course:      course,
		Translation: courseTranslation,
	}
}

func assertCourseWithTranslationEqual(t *testing.T, expectedCourse db.CourseWithTranslation, gotCourse db.CourseWithTranslation) {
	getCourseRow := db.GetCourseRow{
		Uuid:        gotCourse.Uuid,
		CreatedAt:   gotCourse.CreatedAt,
		ModifiedAt:  gotCourse.ModifiedAt,
		DeletedAt:   gotCourse.DeletedAt,
		Subject:     gotCourse.Subject,
		Price:       gotCourse.Price,
		Discount:    gotCourse.Discount,
		IsActive:    gotCourse.IsActive,
		Difficulty:  gotCourse.Difficulty,
		Language:    gotCourse.Translation.Language,
		Name:        gotCourse.Translation.Name,
		Description: gotCourse.Translation.Description,
		Bullets:     gotCourse.Translation.Bullets,
	}
	assertCourseEqual(t, expectedCourse, getCourseRow)
}

func assertCourseListEqual(t *testing.T, expectedCourse db.CourseWithTranslation, gotCourse db.ListCoursesRow) {
	getCourseRow := db.GetCourseRow{
		Uuid:        gotCourse.Uuid,
		CreatedAt:   gotCourse.CreatedAt,
		ModifiedAt:  gotCourse.ModifiedAt,
		DeletedAt:   gotCourse.DeletedAt,
		Subject:     gotCourse.Subject,
		Price:       gotCourse.Price,
		Discount:    gotCourse.Discount,
		IsActive:    gotCourse.IsActive,
		Difficulty:  gotCourse.Difficulty,
		Language:    gotCourse.Language,
		Name:        gotCourse.Name,
		Description: gotCourse.Description,
		Bullets:     gotCourse.Bullets,
	}
	assertCourseEqual(t, expectedCourse, getCourseRow)
}

func assertCourseEqual(t *testing.T, expectedCourse db.CourseWithTranslation, gotCourse db.GetCourseRow) {
	assert.NotNil(t, gotCourse)

	require.Equal(t, expectedCourse.Uuid, gotCourse.Uuid)
	require.Equal(t, expectedCourse.Subject, gotCourse.Subject)
	require.Equal(t, expectedCourse.Price, gotCourse.Price)
	require.Equal(t, expectedCourse.Discount, gotCourse.Discount)
	require.Equal(t, expectedCourse.IsActive, gotCourse.IsActive)
	require.Equal(t, expectedCourse.Difficulty, gotCourse.Difficulty)
	require.Equal(t, expectedCourse.Translation.Name, gotCourse.Name)
	require.Equal(t, expectedCourse.Translation.Description, gotCourse.Description)
	require.Equal(t, expectedCourse.Translation.Bullets, gotCourse.Bullets)
	require.Equal(t, expectedCourse.Translation.Language, gotCourse.Language)
	require.Equal(t, expectedCourse.Translation.CourseUuid, gotCourse.Uuid)

	require.NotZero(t, gotCourse.CreatedAt)
	require.NotZero(t, gotCourse.ModifiedAt)
	require.Nil(t, gotCourse.DeletedAt)
}

func TestCreateCourse(t *testing.T) {
	createRandomCourse(t)
}

func TestGetCourse(t *testing.T) {
	course := createRandomCourse(t)

	gotCourse, err := testQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotCourse)

	assertCourseEqual(t, course, gotCourse)
}

func TestUpdateCourse(t *testing.T) {
	course := createRandomCourse(t)

	rnd := getRandomInt()
	subject := "javascript"
	price := int16(150)
	discount := int16(5)
	isActive := false
	difficulty := int16(2)
	updateParams := db.UpdateCourseParams{
		Uuid:       course.Uuid,
		Subject:    &subject,
		Price:      &price,
		Discount:   &discount,
		IsActive:   &isActive,
		Difficulty: &difficulty,
	}

	updatedCourse, err := testQueries.UpdateCourse(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedCourse)

	require.Equal(t, *updateParams.Subject, updatedCourse.Subject)
	require.Equal(t, *updateParams.Price, updatedCourse.Price)
	require.Equal(t, *updateParams.Discount, updatedCourse.Discount)
	require.Equal(t, *updateParams.IsActive, updatedCourse.IsActive)
	require.Equal(t, *updateParams.Difficulty, updatedCourse.Difficulty)

	require.NotZero(t, updatedCourse.ModifiedAt)
	require.Nil(t, updatedCourse.DeletedAt)

	language := "en"
	name := fmt.Sprintf("Updated Test Course %d", rnd)
	description := fmt.Sprintf("Updated Test Description %d", rnd)
	bullets := "Updated Test Bullets 1\nUpdated Test Bullets 2\nUpdated Test Bullets 3"
	updateTranslationParams := db.UpdateCourseTranslationParams{
		Uuid:        course.Uuid,
		Language:    language,
		Name:        &name,
		Description: &description,
		Bullets:     &bullets,
	}
	updateTranslation, err := testQueries.UpdateCourseTranslation(context.Background(), updateTranslationParams)
	require.NoError(t, err)
	require.NotEmpty(t, updateTranslation)

	require.Equal(t, updateTranslationParams.Language, updateTranslation.Language)
	require.Equal(t, *updateTranslationParams.Name, updateTranslation.Name)
	require.Equal(t, *updateTranslationParams.Description, updateTranslation.Description)
	require.Equal(t, *updateTranslationParams.Bullets, updateTranslation.Bullets)
}

func TestDeleteCourse(t *testing.T) {
	course := createRandomCourse(t)

	err := testQueries.DeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err := testQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotCourse)
}

func TestHardDeleteCourse(t *testing.T) {
	course := createRandomCourse(t)

	err := testQueries.HardDeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err := testQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotCourse)
}

func TestUndeleteCourse(t *testing.T) {
	course := createRandomCourse(t)

	err := testQueries.DeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err := testQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotCourse)

	err = testQueries.UndeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err = testQueries.GetCourse(context.Background(), db.GetCourseParams{
		Uuid:     course.Uuid,
		Language: "en",
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotCourse)

	assertCourseEqual(t, course, gotCourse)
}

func TestGetCourseFull(t *testing.T) {
	course := createRandomCourse(t)
	lesson1 := createRandomLesson(t, &course)
	lesson2 := createRandomLesson(t, &course)
	exercise1 := createRandomExercise(t, &lesson1)
	exercise2 := createRandomExercise(t, &lesson1)
	exercise3 := createRandomExercise(t, &lesson2)

	courseFull, err := testQueries.GetCourseFull(context.Background(), course.Uuid, "en")
	require.NoError(t, err)
	require.NotEmpty(t, courseFull)

	assertCourseWithTranslationEqual(t, course, courseFull.CourseWithTranslation)

	require.Len(t, courseFull.Lessons, 2)

	var foundLesson1, foundLesson2 *db.LessonFull
	for i := range courseFull.Lessons {
		if courseFull.Lessons[i].Uuid == lesson1.Uuid {
			foundLesson1 = &courseFull.Lessons[i]
		}
		if courseFull.Lessons[i].Uuid == lesson2.Uuid {
			foundLesson2 = &courseFull.Lessons[i]
		}
	}

	require.NotNil(t, foundLesson1, "Lesson 1 should be found")
	require.NotNil(t, foundLesson2, "Lesson 2 should be found")

	assertLessonWithTranslationEqual(t, lesson1, foundLesson1.LessonWithTranslation)
	require.Len(t, foundLesson1.Exercises, 2)

	var foundExercise1, foundExercise2 *db.ExerciseWithTranslation
	for i := range foundLesson1.Exercises {
		if foundLesson1.Exercises[i].Uuid == exercise1.Uuid {
			foundExercise1 = &foundLesson1.Exercises[i]
		}
		if foundLesson1.Exercises[i].Uuid == exercise2.Uuid {
			foundExercise2 = &foundLesson1.Exercises[i]
		}
	}

	require.NotNil(t, foundExercise1, "Exercise 1 should be found in lesson 1")
	require.NotNil(t, foundExercise2, "Exercise 2 should be found in lesson 1")
	assertExerciseWithTranslationEqual(t, exercise1, *foundExercise1)
	assertExerciseWithTranslationEqual(t, exercise2, *foundExercise2)

	assertLessonWithTranslationEqual(t, lesson2, foundLesson2.LessonWithTranslation)
	require.Len(t, foundLesson2.Exercises, 1)

	require.Equal(t, exercise3.Uuid, foundLesson2.Exercises[0].Uuid)
	assertExerciseWithTranslationEqual(t, exercise3, foundLesson2.Exercises[0])
}

func TestCountCourses(t *testing.T) {
	initialCount, err := testQueries.CountCourses(context.Background())
	require.NoError(t, err)

	course1 := createRandomCourse(t)
	_ = createRandomCourse(t)

	count, err := testQueries.CountCourses(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+2, count)

	err = testQueries.DeleteCourse(context.Background(), course1.Uuid)
	require.NoError(t, err)

	count, err = testQueries.CountCourses(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
}

func TestListCourses(t *testing.T) {
	course1 := createRandomCourse(t)
	course2 := createRandomCourse(t)

	params := db.ListCoursesParams{
		Limit:    10,
		Offset:   0,
		Language: "en",
		Subject:  nil,
		IsActive: nil,
	}

	courses, err := testQueries.ListCourses(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(courses), 2)

	var foundCourse1, foundCourse2 bool
	for _, course := range courses {
		if course.Uuid == course1.Uuid {
			foundCourse1 = true
			assertCourseListEqual(t, course1, course)
		}
		if course.Uuid == course2.Uuid {
			foundCourse2 = true
			assertCourseListEqual(t, course2, course)
		}
	}
	require.True(t, foundCourse1)
	require.True(t, foundCourse2)
}

func TestListCoursesWithFilters(t *testing.T) {
	course1 := createRandomCourse(t)

	course2Params := db.CreateCourseParams{
		Subject:    "javascript",
		Price:      100,
		Discount:   0,
		IsActive:   false,
		Difficulty: 1,
	}
	_, err := testQueries.CreateCourse(context.Background(), course2Params)
	require.NoError(t, err)

	subject := "python"
	isActive := true
	params := db.ListCoursesParams{
		Limit:    10,
		Offset:   0,
		Language: "en",
		Subject:  &subject,
		IsActive: &isActive,
	}

	courses, err := testQueries.ListCourses(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(courses), 1)

	var foundCourse1 bool
	for _, course := range courses {
		require.Equal(t, "python", course.Subject)
		require.True(t, course.IsActive)
		if course.Uuid == course1.Uuid {
			foundCourse1 = true
		}
	}
	require.True(t, foundCourse1)
}

func TestCreateCourseTranslationWithConflict(t *testing.T) {
	course := createRandomCourse(t)
	_, err := testQueries.CreateCourseTranslation(context.Background(), db.CreateCourseTranslationParams{
		CourseUuid:  course.Uuid,
		Language:    "en",
		Name:        "Test Course",
		Description: "Test Description",
		Bullets:     "Test Bullets",
	})

	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_course_translations_course_language"))
}

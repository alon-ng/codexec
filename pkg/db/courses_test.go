package db_test

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomCourse(t *testing.T) db.Course {
	rnd := getRandomInt()
	params := db.CreateCourseParams{
		Name:        fmt.Sprintf("Test Course %d", rnd),
		Description: fmt.Sprintf("Test Description %d", rnd),
		Subject:     "Python",
		Price:       100,
		Discount:    0,
		IsActive:    true,
		Difficulty:  1,
		Bullets:     "Test Bullets 1\nTest Bullets 2\nTest Bullets 3",
	}

	course, err := testQueries.CreateCourse(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, course)

	require.Equal(t, params.Name, course.Name)
	require.Equal(t, params.Description, course.Description)
	require.Equal(t, params.Subject, course.Subject)
	require.Equal(t, params.Price, course.Price)
	require.Equal(t, params.Discount, course.Discount)
	require.Equal(t, params.IsActive, course.IsActive)
	require.Equal(t, params.Difficulty, course.Difficulty)
	require.Equal(t, params.Bullets, course.Bullets)

	require.NotZero(t, course.Uuid)
	require.NotZero(t, course.CreatedAt)
	require.NotZero(t, course.ModifiedAt)
	require.Nil(t, course.DeletedAt)

	return course
}

func assertCourseEqual(t *testing.T, expectedCourse db.Course, gotCourse db.Course) {
	assert.NotNil(t, gotCourse)

	require.Equal(t, expectedCourse.Uuid, gotCourse.Uuid)
	require.Equal(t, expectedCourse.Name, gotCourse.Name)
	require.Equal(t, expectedCourse.Description, gotCourse.Description)
	require.Equal(t, expectedCourse.Subject, gotCourse.Subject)
	require.Equal(t, expectedCourse.Price, gotCourse.Price)
	require.Equal(t, expectedCourse.Discount, gotCourse.Discount)
	require.Equal(t, expectedCourse.IsActive, gotCourse.IsActive)
	require.Equal(t, expectedCourse.Difficulty, gotCourse.Difficulty)
	require.Equal(t, expectedCourse.Bullets, gotCourse.Bullets)

	require.NotZero(t, gotCourse.CreatedAt)
	require.NotZero(t, gotCourse.ModifiedAt)
	require.Nil(t, gotCourse.DeletedAt)
}

func TestCreateCourse(t *testing.T) {
	createRandomCourse(t)
}

func TestGetCourse(t *testing.T) {
	course := createRandomCourse(t)

	gotCourse, err := testQueries.GetCourse(context.Background(), course.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotCourse)

	assertCourseEqual(t, course, gotCourse)
}

func TestGetCourseByName(t *testing.T) {
	course := createRandomCourse(t)

	gotCourse, err := testQueries.GetCourseByName(context.Background(), course.Name)
	require.NoError(t, err)
	require.NotEmpty(t, gotCourse)

	assertCourseEqual(t, course, gotCourse)
}

func TestUpdateCourse(t *testing.T) {
	course := createRandomCourse(t)

	rnd := getRandomInt()
	updateParams := db.UpdateCourseParams{
		Uuid:        course.Uuid,
		Name:        fmt.Sprintf("Updated Test Course %d", rnd),
		Description: fmt.Sprintf("Updated Test Description %d", rnd),
		Subject:     "TypeScript",
		Price:       150,
		Discount:    5,
		IsActive:    false,
		Difficulty:  2,
		Bullets:     "Updated Test Bullets 1\nUpdated Test Bullets 2\nUpdated Test Bullets 3",
	}

	updatedCourse, err := testQueries.UpdateCourse(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedCourse)

	require.Equal(t, updateParams.Name, updatedCourse.Name)
	require.Equal(t, updateParams.Description, updatedCourse.Description)
	require.Equal(t, updateParams.Subject, updatedCourse.Subject)
	require.Equal(t, updateParams.Price, updatedCourse.Price)
	require.Equal(t, updateParams.Discount, updatedCourse.Discount)
	require.Equal(t, updateParams.IsActive, updatedCourse.IsActive)
	require.Equal(t, updateParams.Difficulty, updatedCourse.Difficulty)
	require.Equal(t, updateParams.Bullets, updatedCourse.Bullets)

	require.NotZero(t, updatedCourse.ModifiedAt)
	require.Nil(t, updatedCourse.DeletedAt)
}

func TestDeleteCourse(t *testing.T) {
	course := createRandomCourse(t)

	err := testQueries.DeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err := testQueries.GetCourse(context.Background(), course.Uuid)
	require.Error(t, err)
	require.Empty(t, gotCourse)
}

func TestHardDeleteCourse(t *testing.T) {
	course := createRandomCourse(t)

	err := testQueries.HardDeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err := testQueries.GetCourse(context.Background(), course.Uuid)
	require.Error(t, err)
	require.Empty(t, gotCourse)
}

func TestUndeleteCourse(t *testing.T) {
	course := createRandomCourse(t)

	err := testQueries.DeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err := testQueries.GetCourse(context.Background(), course.Uuid)
	require.Error(t, err)
	require.Empty(t, gotCourse)

	err = testQueries.UndeleteCourse(context.Background(), course.Uuid)
	require.NoError(t, err)

	gotCourse, err = testQueries.GetCourse(context.Background(), course.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotCourse)

	assertCourseEqual(t, course, gotCourse)
}

func TestCreateCourseConflict(t *testing.T) {
	course := createRandomCourse(t)

	params := db.CreateCourseParams{
		Name:        course.Name,
		Description: course.Description,
		Subject:     course.Subject,
		Price:       course.Price,
		Discount:    course.Discount,
		IsActive:    course.IsActive,
		Difficulty:  course.Difficulty,
		Bullets:     course.Bullets,
	}

	_, err := testQueries.CreateCourse(context.Background(), params)
	require.Error(t, err)
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "courses_name_key"))
}

func TestGetCourseFull(t *testing.T) {
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

	exercise1Data := json.RawMessage(fmt.Sprintf(`{"answer": "Answer 1", "question": "Question 1 %d"}`, rnd))
	exercise1Params := db.CreateExerciseParams{
		LessonUuid:  lesson1.Uuid,
		Name:        fmt.Sprintf("Test Exercise 1 %d", rnd),
		Description: fmt.Sprintf("Test Exercise Description 1 %d", rnd),
		OrderIndex:  1,
		Reward:      10,
		Type:        db.ExerciseTypeQuiz,
		Data:        exercise1Data,
	}
	exercise1, err := testQueries.CreateExercise(context.Background(), exercise1Params)
	require.NoError(t, err)

	exercise2Data := json.RawMessage(fmt.Sprintf(`{"answer": "Answer 2", "question": "Question 2 %d"}`, rnd))
	exercise2Params := db.CreateExerciseParams{
		LessonUuid:  lesson1.Uuid,
		Name:        fmt.Sprintf("Test Exercise 2 %d", rnd),
		Description: fmt.Sprintf("Test Exercise Description 2 %d", rnd),
		OrderIndex:  2,
		Reward:      20,
		Type:        db.ExerciseTypeCode,
		Data:        exercise2Data,
	}
	exercise2, err := testQueries.CreateExercise(context.Background(), exercise2Params)
	require.NoError(t, err)

	exercise3Data := json.RawMessage(fmt.Sprintf(`{"answer": "Answer 3", "question": "Question 3 %d"}`, rnd))
	exercise3Params := db.CreateExerciseParams{
		LessonUuid:  lesson2.Uuid,
		Name:        fmt.Sprintf("Test Exercise 3 %d", rnd),
		Description: fmt.Sprintf("Test Exercise Description 3 %d", rnd),
		OrderIndex:  1,
		Reward:      15,
		Type:        db.ExerciseTypeQuiz,
		Data:        exercise3Data,
	}
	exercise3, err := testQueries.CreateExercise(context.Background(), exercise3Params)
	require.NoError(t, err)

	courseFull, err := testQueries.GetCourseFull(context.Background(), course.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, courseFull)

	assertCourseEqual(t, course, courseFull.Course)

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

	assertLessonEqual(t, lesson1, foundLesson1.Lesson)
	require.Len(t, foundLesson1.Exercises, 2)

	var foundExercise1, foundExercise2 *db.Exercise
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
	assertExerciseEqual(t, exercise1, *foundExercise1)
	assertExerciseEqual(t, exercise2, *foundExercise2)

	assertLessonEqual(t, lesson2, foundLesson2.Lesson)
	require.Len(t, foundLesson2.Exercises, 1)

	require.Equal(t, exercise3.Uuid, foundLesson2.Exercises[0].Uuid)
	assertExerciseEqual(t, exercise3, foundLesson2.Exercises[0])
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
			assertCourseEqual(t, course1, course)
		}
		if course.Uuid == course2.Uuid {
			foundCourse2 = true
			assertCourseEqual(t, course2, course)
		}
	}
	require.True(t, foundCourse1)
	require.True(t, foundCourse2)
}

func TestListCoursesWithFilters(t *testing.T) {
	course1 := createRandomCourse(t)

	rnd := getRandomInt()
	course2Params := db.CreateCourseParams{
		Name:        fmt.Sprintf("Test Course %d", rnd),
		Description: fmt.Sprintf("Test Description %d", rnd),
		Subject:     "TypeScript",
		Price:       100,
		Discount:    0,
		IsActive:    false,
		Difficulty:  1,
		Bullets:     "Test Bullets",
	}
	_, err := testQueries.CreateCourse(context.Background(), course2Params)
	require.NoError(t, err)

	subject := "Python"
	isActive := true
	params := db.ListCoursesParams{
		Limit:    10,
		Offset:   0,
		Subject:  &subject,
		IsActive: &isActive,
	}

	courses, err := testQueries.ListCourses(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(courses), 1)

	var foundCourse1 bool
	for _, course := range courses {
		require.Equal(t, "Python", course.Subject)
		require.True(t, course.IsActive)
		if course.Uuid == course1.Uuid {
			foundCourse1 = true
		}
	}
	require.True(t, foundCourse1)
}

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

func createRandomExercise(t *testing.T) db.Exercise {
	lesson := createRandomLesson(t)

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

	require.Equal(t, params.LessonUuid, exercise.LessonUuid)
	require.Equal(t, params.Name, exercise.Name)
	require.Equal(t, params.Description, exercise.Description)
	require.Equal(t, params.OrderIndex, exercise.OrderIndex)
	require.Equal(t, params.Reward, exercise.Reward)
	require.Equal(t, params.Type, exercise.Type)
	require.Equal(t, params.Data, exercise.Data)

	require.NotZero(t, exercise.Uuid)
	require.NotZero(t, exercise.CreatedAt)
	require.NotZero(t, exercise.ModifiedAt)
	require.Nil(t, exercise.DeletedAt)

	return exercise
}

func assertExerciseEqual(t *testing.T, expectedExercise db.Exercise, gotExercise db.Exercise) {
	assert.NotNil(t, gotExercise)

	require.Equal(t, expectedExercise.Uuid, gotExercise.Uuid)
	require.Equal(t, expectedExercise.LessonUuid, gotExercise.LessonUuid)
	require.Equal(t, expectedExercise.Name, gotExercise.Name)
	require.Equal(t, expectedExercise.Description, gotExercise.Description)
	require.Equal(t, expectedExercise.OrderIndex, gotExercise.OrderIndex)
	require.Equal(t, expectedExercise.Reward, gotExercise.Reward)
	require.Equal(t, expectedExercise.Type, gotExercise.Type)
	require.Equal(t, expectedExercise.Data, gotExercise.Data)

	require.NotZero(t, gotExercise.CreatedAt)
	require.NotZero(t, gotExercise.ModifiedAt)
	require.Nil(t, gotExercise.DeletedAt)
}

func TestCreateExercise(t *testing.T) {
	createRandomExercise(t)
}

func TestGetExercise(t *testing.T) {
	exercise := createRandomExercise(t)

	gotExercise, err := testQueries.GetExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotExercise)

	assertExerciseEqual(t, exercise, gotExercise)
}

func TestUpdateExercise(t *testing.T) {
	exercise := createRandomExercise(t)
	lesson := createRandomLesson(t)

	rnd := getRandomInt()
	updatedData := json.RawMessage(`{"answer": "Updated Answer", "question": "Updated Question"}`)
	updateParams := db.UpdateExerciseParams{
		Uuid:        exercise.Uuid,
		LessonUuid:  lesson.Uuid,
		Name:        fmt.Sprintf("Updated Test Exercise %d", rnd),
		Description: fmt.Sprintf("Updated Test Description %d", rnd),
		OrderIndex:  2,
		Reward:      20,
		Type:        db.ExerciseTypeCode,
		Data:        updatedData,
	}

	updatedExercise, err := testQueries.UpdateExercise(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedExercise)

	require.Equal(t, updateParams.LessonUuid, updatedExercise.LessonUuid)
	require.Equal(t, updateParams.Name, updatedExercise.Name)
	require.Equal(t, updateParams.Description, updatedExercise.Description)
	require.Equal(t, updateParams.OrderIndex, updatedExercise.OrderIndex)
	require.Equal(t, updateParams.Reward, updatedExercise.Reward)
	require.Equal(t, updateParams.Type, updatedExercise.Type)
	require.Equal(t, updateParams.Data, updatedExercise.Data)

	require.NotZero(t, updatedExercise.ModifiedAt)
	require.Nil(t, updatedExercise.DeletedAt)
}

func TestDeleteExercise(t *testing.T) {
	exercise := createRandomExercise(t)

	err := testQueries.DeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err := testQueries.GetExercise(context.Background(), exercise.Uuid)
	require.Error(t, err)
	require.Empty(t, gotExercise)
}

func TestHardDeleteExercise(t *testing.T) {
	exercise := createRandomExercise(t)

	err := testQueries.HardDeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err := testQueries.GetExercise(context.Background(), exercise.Uuid)
	require.Error(t, err)
	require.Empty(t, gotExercise)
}

func TestUndeleteExercise(t *testing.T) {
	exercise := createRandomExercise(t)

	err := testQueries.DeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err := testQueries.GetExercise(context.Background(), exercise.Uuid)
	require.Error(t, err)
	require.Empty(t, gotExercise)

	err = testQueries.UndeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err = testQueries.GetExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotExercise)

	assertExerciseEqual(t, exercise, gotExercise)
}

func TestCreateExerciseConflict(t *testing.T) {
	exercise := createRandomExercise(t)

	params := db.CreateExerciseParams{
		LessonUuid:  exercise.LessonUuid,
		Name:        exercise.Name,
		Description: exercise.Description,
		OrderIndex:  exercise.OrderIndex,
		Reward:      exercise.Reward,
		Type:        db.ExerciseTypeQuiz,
		Data:        exercise.Data,
	}

	_, err := testQueries.CreateExercise(context.Background(), params)
	require.Error(t, err)
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "exercises_name_key"))
}

func TestCountExercises(t *testing.T) {
	initialCount, err := testQueries.CountExercises(context.Background())
	require.NoError(t, err)

	exercise1 := createRandomExercise(t)
	_ = createRandomExercise(t)

	count, err := testQueries.CountExercises(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+2, count)

	err = testQueries.DeleteExercise(context.Background(), exercise1.Uuid)
	require.NoError(t, err)

	count, err = testQueries.CountExercises(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
}

func TestListExercises(t *testing.T) {
	exercise1 := createRandomExercise(t)
	exercise2 := createRandomExercise(t)

	params := db.ListExercisesParams{
		Limit:     10,
		Offset:    0,
		LessonUuid: nil,
	}

	exercises, err := testQueries.ListExercises(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(exercises), 2)

	var foundExercise1, foundExercise2 bool
	for _, exercise := range exercises {
		if exercise.Uuid == exercise1.Uuid {
			foundExercise1 = true
			assertExerciseEqual(t, exercise1, exercise)
		}
		if exercise.Uuid == exercise2.Uuid {
			foundExercise2 = true
			assertExerciseEqual(t, exercise2, exercise)
		}
	}
	require.True(t, foundExercise1)
	require.True(t, foundExercise2)
}

func TestListExercisesWithLessonFilter(t *testing.T) {
	lesson := createRandomLesson(t)

	rnd := getRandomInt()
	exercise1Data := json.RawMessage(fmt.Sprintf(`{"answer": "Answer 1", "question": "Question 1 %d"}`, rnd))
	exercise1Params := db.CreateExerciseParams{
		LessonUuid:  lesson.Uuid,
		Name:        fmt.Sprintf("Test Exercise 1 %d", rnd),
		Description: fmt.Sprintf("Test Description 1 %d", rnd),
		OrderIndex:  1,
		Reward:      10,
		Type:        db.ExerciseTypeQuiz,
		Data:        exercise1Data,
	}
	exercise1, err := testQueries.CreateExercise(context.Background(), exercise1Params)
	require.NoError(t, err)

	exercise2Data := json.RawMessage(fmt.Sprintf(`{"answer": "Answer 2", "question": "Question 2 %d"}`, rnd))
	exercise2Params := db.CreateExerciseParams{
		LessonUuid:  lesson.Uuid,
		Name:        fmt.Sprintf("Test Exercise 2 %d", rnd),
		Description: fmt.Sprintf("Test Description 2 %d", rnd),
		OrderIndex:  2,
		Reward:      20,
		Type:        db.ExerciseTypeCode,
		Data:        exercise2Data,
	}
	exercise2, err := testQueries.CreateExercise(context.Background(), exercise2Params)
	require.NoError(t, err)

	params := db.ListExercisesParams{
		Limit:     10,
		Offset:    0,
		LessonUuid: &lesson.Uuid,
	}

	exercises, err := testQueries.ListExercises(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(exercises), 2)

	var foundExercise1, foundExercise2 bool
	for _, exercise := range exercises {
		require.Equal(t, lesson.Uuid, exercise.LessonUuid)
		if exercise.Uuid == exercise1.Uuid {
			foundExercise1 = true
			assertExerciseEqual(t, exercise1, exercise)
		}
		if exercise.Uuid == exercise2.Uuid {
			foundExercise2 = true
			assertExerciseEqual(t, exercise2, exercise)
		}
	}
	require.True(t, foundExercise1)
	require.True(t, foundExercise2)
}

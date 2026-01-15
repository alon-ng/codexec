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

func createRandomExercise(t *testing.T, lesson *db.LessonWithTranslation) db.ExerciseWithTranslation {
	if lesson == nil {
		l := createRandomLesson(t, nil)
		lesson = &l
	}

	rnd := getRandomInt()
	testData := json.RawMessage(fmt.Sprintf(`{"answer": "Test Answer", "question": "Test Question %d"}`, rnd))
	params := db.CreateExerciseParams{
		LessonUuid: lesson.Uuid,
		OrderIndex: 1,
		Reward:     10,
		Type:       db.ExerciseTypeQuiz,
		Data:       testData,
	}

	exercise, err := testQueries.CreateExercise(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, exercise)

	require.Equal(t, params.LessonUuid, exercise.LessonUuid)
	require.Equal(t, params.OrderIndex, exercise.OrderIndex)
	require.Equal(t, params.Reward, exercise.Reward)
	require.Equal(t, params.Type, exercise.Type)
	require.Equal(t, params.Data, exercise.Data)

	require.NotZero(t, exercise.Uuid)
	require.NotZero(t, exercise.CreatedAt)
	require.NotZero(t, exercise.ModifiedAt)
	require.Nil(t, exercise.DeletedAt)

	translationParams := db.CreateExerciseTranslationParams{
		ExerciseUuid: exercise.Uuid,
		Language:     "en",
		Name:         fmt.Sprintf("Test Exercise %d", rnd),
		Description:  fmt.Sprintf("Test Description %d", rnd),
	}

	exerciseTranslation, err := testQueries.CreateExerciseTranslation(context.Background(), translationParams)
	require.NoError(t, err)
	require.NotEmpty(t, exerciseTranslation)

	require.Equal(t, translationParams.ExerciseUuid, exerciseTranslation.ExerciseUuid)
	require.Equal(t, translationParams.Language, exerciseTranslation.Language)
	require.Equal(t, translationParams.Name, exerciseTranslation.Name)
	require.Equal(t, translationParams.Description, exerciseTranslation.Description)

	return db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: exerciseTranslation,
	}
}

func assertExerciseWithTranslationEqual(t *testing.T, expectedExercise db.ExerciseWithTranslation, gotExercise db.ExerciseWithTranslation) {
	getExerciseRow := db.GetExerciseRow{
		Uuid:         gotExercise.Uuid,
		CreatedAt:    gotExercise.CreatedAt,
		ModifiedAt:   gotExercise.ModifiedAt,
		DeletedAt:    gotExercise.DeletedAt,
		LessonUuid:   gotExercise.LessonUuid,
		OrderIndex:   gotExercise.OrderIndex,
		Reward:       gotExercise.Reward,
		Type:         gotExercise.Type,
		Data:         gotExercise.Data,
		ExerciseUuid: gotExercise.Translation.ExerciseUuid,
		Language:     gotExercise.Translation.Language,
		Name:         gotExercise.Translation.Name,
		Description:  gotExercise.Translation.Description,
	}
	assertExerciseEqual(t, expectedExercise, getExerciseRow)
}

func assertExerciseListEqual(t *testing.T, expectedExercise db.ExerciseWithTranslation, gotExercise db.ListExercisesRow) {
	getExerciseRow := db.GetExerciseRow{
		Uuid:         gotExercise.Uuid,
		CreatedAt:    gotExercise.CreatedAt,
		ModifiedAt:   gotExercise.ModifiedAt,
		DeletedAt:    gotExercise.DeletedAt,
		LessonUuid:   gotExercise.LessonUuid,
		OrderIndex:   gotExercise.OrderIndex,
		Reward:       gotExercise.Reward,
		Type:         gotExercise.Type,
		Data:         gotExercise.Data,
		ExerciseUuid: gotExercise.ExerciseUuid,
		Language:     gotExercise.Language,
		Name:         gotExercise.Name,
		Description:  gotExercise.Description,
	}
	assertExerciseEqual(t, expectedExercise, getExerciseRow)
}

func assertExerciseEqual(t *testing.T, expectedExercise db.ExerciseWithTranslation, gotExercise db.GetExerciseRow) {
	assert.NotNil(t, gotExercise)

	require.Equal(t, expectedExercise.Uuid, gotExercise.Uuid)
	require.Equal(t, expectedExercise.LessonUuid, gotExercise.LessonUuid)
	require.Equal(t, expectedExercise.OrderIndex, gotExercise.OrderIndex)
	require.Equal(t, expectedExercise.Reward, gotExercise.Reward)
	require.Equal(t, expectedExercise.Type, gotExercise.Type)
	require.Equal(t, expectedExercise.Data, gotExercise.Data)
	require.Equal(t, expectedExercise.Translation.Name, gotExercise.Name)
	require.Equal(t, expectedExercise.Translation.Description, gotExercise.Description)
	require.Equal(t, expectedExercise.Translation.Language, gotExercise.Language)
	require.Equal(t, expectedExercise.Translation.ExerciseUuid, gotExercise.ExerciseUuid)

	require.NotZero(t, gotExercise.CreatedAt)
	require.NotZero(t, gotExercise.ModifiedAt)
	require.Nil(t, gotExercise.DeletedAt)
}

func TestCreateExercise(t *testing.T) {
	createRandomExercise(t, nil)
}

func TestGetExercise(t *testing.T) {
	exercise := createRandomExercise(t, nil)

	gotExercise, err := testQueries.GetExercise(context.Background(), db.GetExerciseParams{
		Uuid:     exercise.Uuid,
		Language: "en",
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotExercise)

	assertExerciseEqual(t, exercise, gotExercise)
}

func TestUpdateExercise(t *testing.T) {
	lesson := createRandomLesson(t, nil)
	exercise := createRandomExercise(t, &lesson)

	rnd := getRandomInt()
	updatedData := json.RawMessage(`{"answer": "Updated Answer", "question": "Updated Question"}`)
	orderIndex := int16(2)
	reward := int16(20)
	type_ := db.ExerciseTypeQuiz
	updateParams := db.UpdateExerciseParams{
		Uuid:       exercise.Uuid,
		OrderIndex: &orderIndex,
		Reward:     &reward,
		Type:       &type_,
		Data:       &updatedData,
	}

	updatedExercise, err := testQueries.UpdateExercise(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedExercise)

	require.Equal(t, *updateParams.OrderIndex, updatedExercise.OrderIndex)
	require.Equal(t, *updateParams.Reward, updatedExercise.Reward)
	require.Equal(t, *updateParams.Type, updatedExercise.Type)
	require.Equal(t, *updateParams.Data, updatedExercise.Data)

	require.NotZero(t, updatedExercise.ModifiedAt)
	require.Nil(t, updatedExercise.DeletedAt)

	language := "en"
	name := fmt.Sprintf("Updated Test Exercise %d", rnd)
	description := fmt.Sprintf("Updated Test Description %d", rnd)
	updateTranslationParams := db.UpdateExerciseTranslationParams{
		Uuid:        exercise.Uuid,
		Language:    language,
		Name:        &name,
		Description: &description,
	}
	updateTranslation, err := testQueries.UpdateExerciseTranslation(context.Background(), updateTranslationParams)
	require.NoError(t, err)
	require.NotEmpty(t, updateTranslation)

	require.Equal(t, updateTranslationParams.Language, updateTranslation.Language)
	require.Equal(t, *updateTranslationParams.Name, updateTranslation.Name)
	require.Equal(t, *updateTranslationParams.Description, updateTranslation.Description)
}

func TestDeleteExercise(t *testing.T) {
	exercise := createRandomExercise(t, nil)

	err := testQueries.DeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err := testQueries.GetExercise(context.Background(), db.GetExerciseParams{
		Uuid:     exercise.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotExercise)
}

func TestHardDeleteExercise(t *testing.T) {
	exercise := createRandomExercise(t, nil)

	err := testQueries.HardDeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err := testQueries.GetExercise(context.Background(), db.GetExerciseParams{
		Uuid:     exercise.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotExercise)
}

func TestUndeleteExercise(t *testing.T) {
	exercise := createRandomExercise(t, nil)

	err := testQueries.DeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err := testQueries.GetExercise(context.Background(), db.GetExerciseParams{
		Uuid:     exercise.Uuid,
		Language: "en",
	})
	require.Error(t, err)
	require.Empty(t, gotExercise)

	err = testQueries.UndeleteExercise(context.Background(), exercise.Uuid)
	require.NoError(t, err)

	gotExercise, err = testQueries.GetExercise(context.Background(), db.GetExerciseParams{
		Uuid:     exercise.Uuid,
		Language: "en",
	})
	require.NoError(t, err)
	require.NotEmpty(t, gotExercise)

	assertExerciseEqual(t, exercise, gotExercise)
}

func TestCountExercises(t *testing.T) {
	initialCount, err := testQueries.CountExercises(context.Background())
	require.NoError(t, err)

	exercise1 := createRandomExercise(t, nil)
	_ = createRandomExercise(t, nil)

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
	exercise1 := createRandomExercise(t, nil)
	exercise2 := createRandomExercise(t, nil)

	params := db.ListExercisesParams{
		Limit:      10,
		Offset:     0,
		Language:   "en",
		LessonUuid: nil,
	}

	exercises, err := testQueries.ListExercises(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(exercises), 2)

	var foundExercise1, foundExercise2 bool
	for _, exercise := range exercises {
		if exercise.Uuid == exercise1.Uuid {
			foundExercise1 = true
			assertExerciseListEqual(t, exercise1, exercise)
		}
		if exercise.Uuid == exercise2.Uuid {
			foundExercise2 = true
			assertExerciseListEqual(t, exercise2, exercise)
		}
	}
	require.True(t, foundExercise1)
	require.True(t, foundExercise2)
}

func TestListExercisesWithLessonFilter(t *testing.T) {
	lesson := createRandomLesson(t, nil)
	exercise1 := createRandomExercise(t, &lesson)
	exercise2 := createRandomExercise(t, &lesson)

	params := db.ListExercisesParams{
		Limit:      10,
		Offset:     0,
		Language:   "en",
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
			assertExerciseListEqual(t, exercise1, exercise)
		}
		if exercise.Uuid == exercise2.Uuid {
			foundExercise2 = true
			assertExerciseListEqual(t, exercise2, exercise)
		}
	}
	require.True(t, foundExercise1)
	require.True(t, foundExercise2)
}

func TestCreateExerciseTranslationWithConflict(t *testing.T) {
	exercise := createRandomExercise(t, nil)
	_, err := testQueries.CreateExerciseTranslation(context.Background(), db.CreateExerciseTranslationParams{
		ExerciseUuid: exercise.Uuid,
		Language:     "en",
		Name:         "Test Exercise",
		Description:  "Test Description",
	})
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_exercise_translations_exercise_language"))
}

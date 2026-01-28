package models

import (
	"codim/pkg/db"
	"codim/pkg/fs"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ExerciseCodeData = fs.Entry
type ExerciseQuizData = map[string]interface{}

type ExerciseTranslationCodeData struct {
	Instructions string `json:"instructions" binding:"required" example:"<p>Hello! Start writing your exercise instructions here...</p>"`
}

type ExerciseTranslationQuizDataQuestion struct {
	Question string            `json:"question" binding:"required" example:"What is the capital of France?"`
	Answers  map[string]string `json:"answers" binding:"required" example:"{\"1\": \"Paris\", \"2\": \"London\", \"3\": \"Berlin\", \"4\": \"Madrid\"}"`
}

type ExerciseTranslationQuizData = map[string]ExerciseTranslationQuizDataQuestion

type Exercise struct {
	Uuid       uuid.UUID         `json:"uuid" binding:"required"`
	CreatedAt  time.Time         `json:"created_at" binding:"required"`
	ModifiedAt time.Time         `json:"modified_at" binding:"required"`
	DeletedAt  *time.Time        `json:"deleted_at,omitempty"`
	LessonUuid uuid.UUID         `json:"lesson_uuid" binding:"required"`
	OrderIndex int16             `json:"order_index" binding:"required" example:"1"`
	Reward     int16             `json:"reward" binding:"required" example:"10"`
	Type       db.ExerciseType   `json:"type" binding:"required" example:"quiz"`
	CodeData   *ExerciseCodeData `json:"code_data,omitempty"`
	QuizData   *ExerciseQuizData `json:"quiz_data,omitempty"`
}

type ExerciseTranslation struct {
	Uuid         uuid.UUID                    `json:"uuid" binding:"required"`
	ExerciseUuid uuid.UUID                    `json:"exercise_uuid" binding:"required"`
	Language     string                       `json:"language" binding:"required" example:"en"`
	Name         string                       `json:"name" binding:"required" example:"Hello World"`
	Description  string                       `json:"description" binding:"required" example:"Print Hello World"`
	CodeData     *ExerciseTranslationCodeData `json:"code_data,omitempty"`
	QuizData     *ExerciseTranslationQuizData `json:"quiz_data,omitempty"`
}

type ExerciseWithTranslation struct {
	Exercise
	Translation ExerciseTranslation `json:"translation" binding:"required"`
}

func ToExercise(d db.Exercise) (Exercise, error) {
	var codeData *ExerciseCodeData
	var quizData *ExerciseQuizData
	if d.CodeData != nil {
		err := json.Unmarshal(*d.CodeData, &codeData)
		if err != nil {
			return Exercise{}, err
		}
	}
	if d.QuizData != nil {
		err := json.Unmarshal(*d.QuizData, &quizData)
		if err != nil {
			return Exercise{}, err
		}
	}

	return Exercise{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		LessonUuid: d.LessonUuid,
		OrderIndex: d.OrderIndex,
		Reward:     d.Reward,
		Type:       d.Type,
		CodeData:   codeData,
		QuizData:   quizData,
	}, nil
}

func ToExerciseTranslation(d db.ExerciseTranslation) (ExerciseTranslation, error) {
	var codeData *ExerciseTranslationCodeData
	var quizData *ExerciseTranslationQuizData

	if d.CodeData != nil {
		err := json.Unmarshal(*d.CodeData, &codeData)
		if err != nil {
			return ExerciseTranslation{}, err
		}
	}
	if d.QuizData != nil {
		err := json.Unmarshal(*d.QuizData, &quizData)
		if err != nil {
			return ExerciseTranslation{}, err
		}
	}

	return ExerciseTranslation{
		Uuid:         d.Uuid,
		ExerciseUuid: d.ExerciseUuid,
		Language:     d.Language,
		Name:         d.Name,
		Description:  d.Description,
		CodeData:     codeData,
		QuizData:     quizData,
	}, nil
}

func ToExerciseWithTranslation(d db.ExerciseWithTranslation) (ExerciseWithTranslation, error) {
	exercise, err := ToExercise(d.Exercise)
	if err != nil {
		return ExerciseWithTranslation{}, err
	}

	translation, err := ToExerciseTranslation(d.Translation)
	if err != nil {
		return ExerciseWithTranslation{}, err
	}

	return ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	}, nil
}

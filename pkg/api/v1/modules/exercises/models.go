package exercises

import (
	"codim/pkg/db"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateExerciseRequest struct {
	LessonUuid  uuid.UUID              `json:"lesson_uuid" binding:"required"`
	Type        db.ExerciseType        `json:"type" binding:"required" example:"quiz"`
	OrderIndex  int16                  `json:"order_index" binding:"required" example:"1"`
	Reward      int16                  `json:"reward" binding:"required" example:"10"`
	Data        map[string]interface{} `json:"data" binding:"required"`
	Language    string                 `json:"language" binding:"required" example:"en"`
	Name        string                 `json:"name" binding:"required" example:"Hello World"`
	Description string                 `json:"description" binding:"required" example:"Print Hello World"`
}

type UpdateExerciseRequest struct {
	Uuid        uuid.UUID               `json:"uuid" binding:"required"`
	Language    string                  `json:"language" binding:"required" example:"en"`
	LessonUuid  *uuid.UUID              `json:"lesson_uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrderIndex  *int16                  `json:"order_index" example:"1"`
	Reward      *int16                  `json:"reward" example:"10"`
	Type        *db.ExerciseType        `json:"type" example:"quiz"`
	Data        *map[string]interface{} `json:"data"`
	Name        *string                 `json:"name" example:"Hello World"`
	Description *string                 `json:"description" example:"Print Hello World"`
}

type ListExercisesRequest struct {
	Limit      int32      `json:"limit" form:"limit,default=10" example:"10"`
	Offset     int32      `json:"offset" form:"offset,default=0" example:"0"`
	LessonUuid *uuid.UUID `json:"lesson_uuid" form:"lesson_uuid"`
	Language   string     `json:"language" form:"language,default=en" example:"en"`
}

type AddExerciseTranslationRequest struct {
	ExerciseUuid uuid.UUID `json:"exercise_uuid" binding:"required"`
	Language     string    `json:"language" binding:"required" example:"es"`
	Name         string    `json:"name" binding:"required" example:"Hola Mundo"`
	Description  string    `json:"description" binding:"required" example:"Imprime Hola Mundo"`
}

// Response types
type Exercise struct {
	Uuid       uuid.UUID       `json:"uuid" binding:"required"`
	CreatedAt  time.Time       `json:"created_at" binding:"required"`
	ModifiedAt time.Time       `json:"modified_at" binding:"required"`
	DeletedAt  *time.Time      `json:"deleted_at,omitempty"`
	LessonUuid uuid.UUID       `json:"lesson_uuid" binding:"required"`
	OrderIndex int16           `json:"order_index" binding:"required" example:"1"`
	Reward     int16           `json:"reward" binding:"required" example:"10"`
	Type        db.ExerciseType `json:"type" binding:"required" example:"quiz"`
	Data        json.RawMessage `json:"data" binding:"required"`
}

type ExerciseTranslation struct {
	Uuid         uuid.UUID `json:"uuid" binding:"required"`
	ExerciseUuid uuid.UUID `json:"exercise_uuid" binding:"required"`
	Language     string    `json:"language" binding:"required" example:"en"`
	Name         string    `json:"name" binding:"required" example:"Hello World"`
	Description  string    `json:"description" binding:"required" example:"Print Hello World"`
}

type ExerciseWithTranslation struct {
	Exercise
	Translation ExerciseTranslation `json:"translation" binding:"required"`
}

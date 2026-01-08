package exercises

import "github.com/google/uuid"

type CreateExerciseRequest struct {
	LessonUuid  uuid.UUID              `json:"lesson_uuid" binding:"required"`
	Name        string                 `json:"name" binding:"required" example:"Hello World"`
	Description string                 `json:"description" binding:"required" example:"Print Hello World"`
	OrderIndex  int16                  `json:"order_index" binding:"required" example:"1"`
	Reward      int16                  `json:"reward" binding:"required" example:"10"`
	Data        map[string]interface{} `json:"data" binding:"required"`
}

type UpdateRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
	UpdateExerciseRequest
}

type UpdateExerciseRequest struct {
	LessonUuid  uuid.UUID              `json:"lesson_uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string                 `json:"name" example:"Hello World"`
	Description string                 `json:"description" example:"Print Hello World"`
	OrderIndex  int16                  `json:"order_index" example:"1"`
	Reward      int16                  `json:"reward" example:"10"`
	Data        map[string]interface{} `json:"data"`
}

type ListExercisesRequest struct {
	Limit      int32      `json:"limit" form:"limit,default=10" example:"10"`
	Offset     int32      `json:"offset" form:"offset,default=0" example:"0"`
	LessonUuid *uuid.UUID `json:"lesson_uuid" form:"lesson_uuid"`
}

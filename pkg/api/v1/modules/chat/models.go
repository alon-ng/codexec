package chat

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	Uuid         uuid.UUID `json:"uuid" binding:"required"`
	Ts           time.Time `json:"ts" binding:"required"`
	ExerciseUuid uuid.UUID `json:"exercise_uuid" binding:"required"`
	UserUuid     uuid.UUID `json:"user_uuid" binding:"required"`
	Role         string    `json:"role" binding:"required"`
	Content      string    `json:"content" binding:"required"`
}

type ListChatMessagesRequest struct {
	Limit  int32 `json:"limit" form:"limit,default=10" example:"10" query:"limit"`
	Offset int32 `json:"offset" form:"offset,default=0" example:"0" query:"offset"`
}

type SendChatMessageRequest struct {
	Content              string `json:"content" binding:"required" example:"Hello, how are you?"`
	Code                 string `json:"code" binding:"required" example:"print('Hello, world!')"`
	ExerciseInstructions string `json:"exercise_instructions" binding:"required" example:"Write a function that prints 'Hello, world!'"`
}

package models

import (
	"codim/pkg/db"
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

func ToChatMessage(d db.ChatMessage) ChatMessage {
	return ChatMessage{
		Uuid:         d.Uuid,
		Ts:           d.Ts,
		ExerciseUuid: d.ExerciseUuid,
		UserUuid:     d.UserUuid,
		Role:         d.Role,
		Content:      d.Content,
	}
}

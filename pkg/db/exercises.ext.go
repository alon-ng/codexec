package db

import (
	"time"

	"github.com/google/uuid"
)

type ExerciseResponse struct {
	Uuid        uuid.UUID              `json:"uuid"`
	CreatedAt   time.Time              `json:"created_at"`
	ModifiedAt  time.Time              `json:"modified_at"`
	DeletedAt   *time.Time             `json:"deleted_at"`
	LessonUuid  uuid.UUID              `json:"lesson_uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	OrderIndex  int16                  `json:"order_index"`
	Reward      int16                  `json:"reward"`
	Type        ExerciseType           `json:"type"`
	Data        map[string]interface{} `json:"data"`
}

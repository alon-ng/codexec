package models

import (
	"codim/pkg/db"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid       uuid.UUID  `json:"uuid" binding:"required"`
	CreatedAt  time.Time  `json:"created_at" binding:"required"`
	ModifiedAt time.Time  `json:"modified_at" binding:"required"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	FirstName  string     `json:"first_name" binding:"required" example:"John"`
	LastName   string     `json:"last_name" binding:"required" example:"Doe"`
	Email      string     `json:"email" binding:"required" example:"john.doe@example.com"`
	IsVerified bool       `json:"is_verified" example:"false"`
	Streak     int32      `json:"streak" example:"0"`
	Score      int32      `json:"score" example:"0"`
	IsAdmin    bool       `json:"is_admin" example:"false"`
}

func ToUser(d db.User) User {
	return User{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		FirstName:  d.FirstName,
		LastName:   d.LastName,
		Email:      d.Email,
		IsVerified: d.IsVerified,
		Streak:     d.Streak,
		Score:      d.Score,
		IsAdmin:    d.IsAdmin,
	}
}

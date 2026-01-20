package users

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName  string `json:"first_name" binding:"required" example:"John"`
	LastName   string `json:"last_name" binding:"required" example:"Doe"`
	Email      string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password   string `json:"password" binding:"required" example:"password"`
	IsVerified bool   `json:"is_verified" example:"false"`
	IsAdmin    bool   `json:"is_admin" example:"false"`
}

type UpdateUserRequest struct {
	Uuid      uuid.UUID `json:"uuid" binding:"required"`
	FirstName *string   `json:"first_name" example:"John"`
	LastName  *string   `json:"last_name" example:"Doe"`
}

// Response types
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

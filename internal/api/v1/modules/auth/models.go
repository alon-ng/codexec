package auth

import "codim/internal/db"

type SignupRequest struct {
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
	Email     string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password  string `json:"password" binding:"required" example:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required" example:"password"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  db.User `json:"user"`
}

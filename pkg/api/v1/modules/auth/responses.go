package auth

import "codim/pkg/db"

type AuthResponse struct {
	Token string  `json:"token"`
	User  db.User `json:"user"`
}

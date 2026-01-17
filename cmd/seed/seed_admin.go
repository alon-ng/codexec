package main

import (
	"codim/pkg/api/auth"
	"codim/pkg/db"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

func seedAdmin(ctx context.Context, queries *db.Queries, authProvider *auth.Provider) db.User {
	hashedPassword := authProvider.HashPassword("password")

	u, err := queries.GetUserByEmail(ctx, "admin@example.com")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			admin := db.CreateUserParams{
				FirstName:    "Admin",
				LastName:     "Admin",
				Email:        "admin@example.com",
				PasswordHash: hashedPassword,
			}
			u, err := queries.CreateUser(ctx, admin)
			if err != nil {
				log.Fatalf("Failed to create admin user: %v", err)
			}

			log.Println("Admin user created successfully")
			return u
		} else {
			log.Fatalf("Failed to get admin user: %v", err)
		}
	} else {
		log.Println("Admin user already exists, skipping...")
	}

	return u
}

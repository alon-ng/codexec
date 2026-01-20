package main

import (
	"codim/pkg/api/v1/modules/progress"
	"codim/pkg/db"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func seedUserCourse(ctx context.Context, queries *db.Queries, pool *pgxpool.Pool, userUuid uuid.UUID, courseUuid uuid.UUID) {
	log.Println("Seeding user course...")

	svc := progress.NewService(queries, pool)

	err := svc.InitUserCourse(ctx, userUuid, courseUuid)
	if err != nil {
		log.Fatalf("Failed to init user course: %v", err)
	}

	log.Println("User course seeded successfully!")
}

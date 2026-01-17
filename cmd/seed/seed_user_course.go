package main

import (
	"codim/pkg/db"
	"context"
	"log"

	"github.com/google/uuid"
)

func seedUserCourse(ctx context.Context, queries *db.Queries, userUuid uuid.UUID, courseUuid uuid.UUID) {
	log.Println("Seeding user course...")

	_, err := queries.CreateUserCourse(ctx, db.CreateUserCourseParams{
		UserUuid:   userUuid,
		CourseUuid: courseUuid,
	})

	if err != nil {
		log.Fatalf("Failed to create user course: %v", err)
	}

	log.Println("User course seeded successfully!")
}

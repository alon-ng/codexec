package main

import (
	"codim/pkg/api/auth"
	"codim/pkg/db"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func seed(ctx context.Context, queries *db.Queries, pool *pgxpool.Pool, authProvider *auth.Provider) {
	user := seedAdmin(ctx, queries, authProvider)
	course := seedCourse(ctx, queries)
	seedUserCourse(ctx, queries, pool, user.Uuid, course.Uuid)
}

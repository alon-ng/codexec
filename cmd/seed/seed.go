package main

import (
	"codim/pkg/api/auth"
	"codim/pkg/db"
	"context"
)

func seed(ctx context.Context, queries *db.Queries, authProvider *auth.Provider) {
	user := seedAdmin(ctx, queries, authProvider)
	course := seedCourse(ctx, queries)
	seedUserCourse(ctx, queries, user.Uuid, course.Uuid)
}

package main

import (
	"codim/pkg/api/auth"
	"codim/pkg/db"
	"context"
)

func seed(ctx context.Context, queries *db.Queries, authProvider *auth.Provider) {
	seedAdmin(ctx, queries, authProvider)
	seedCourse(ctx, queries)
}

package me

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	q *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) Me(ctx context.Context, meUUID uuid.UUID) (db.User, *e.APIError) {
	u, err := s.q.GetUser(ctx, meUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, e.NewAPIError(err, ErrMeFailed)
		}

		return db.User{}, e.NewAPIError(err, ErrMeFailed)
	}

	return u, nil
}

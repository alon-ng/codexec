package me

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewService(q *db.Queries, p *pgxpool.Pool) *Service {
	return &Service{q: q, p: p}
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

func (s *Service) ListUserCoursesWithProgress(ctx context.Context, meUUID uuid.UUID, req ListUserCoursesWithProgressRequest) ([]db.UserCourseWithProgress, *e.APIError) {
	userCourses, err := s.q.ListUserCoursesWithProgress(ctx, db.ListUserCoursesWithProgressParams{
		UserUuid: meUUID,
		Language: req.Language,
		Limit:    req.Limit,
		Offset:   req.Offset,
		Subject:  req.Subject,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, e.NewAPIError(err, ErrGetUserCoursesWithProgressFailed)
	}

	userCoursesWithProgress := make([]db.UserCourseWithProgress, len(userCourses))
	for i, userCourse := range userCourses {
		userCoursesWithProgress[i] = userCourse.ToUserCourseWithProgress()
	}

	return userCoursesWithProgress, nil
}

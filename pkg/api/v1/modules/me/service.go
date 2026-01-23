package me

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/progress"
	"codim/pkg/api/v1/modules/users"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q           *db.Queries
	p           *pgxpool.Pool
	progressSvc *progress.Service
}

func NewService(q *db.Queries, p *pgxpool.Pool, progressSvc *progress.Service) *Service {
	return &Service{q: q, p: p, progressSvc: progressSvc}
}

// Conversion functions
func toUser(d db.User) users.User {
	return users.User{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		FirstName:  d.FirstName,
		LastName:   d.LastName,
		Email:      d.Email,
		IsVerified: d.IsVerified,
		Streak:     d.Streak,
		Score:      d.Score,
		IsAdmin:    d.IsAdmin,
	}
}

func (s *Service) Me(ctx context.Context, meUUID uuid.UUID) (users.User, *e.APIError) {
	u, err := s.q.GetUser(ctx, meUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return users.User{}, e.NewAPIError(err, ErrMeFailed)
		}

		return users.User{}, e.NewAPIError(err, ErrMeFailed)
	}

	return toUser(u), nil
}

func (s *Service) ListUserCoursesWithProgress(ctx context.Context, meUUID uuid.UUID, req ListUserCoursesWithProgressRequest) ([]UserCourseWithProgress, *e.APIError) {
	return s.progressSvc.ListUserCoursesWithProgress(ctx, meUUID, req)
}

func (s *Service) GetUserCourseFull(ctx context.Context, meUUID uuid.UUID, courseUUID uuid.UUID) (UserCourseFull, *e.APIError) {
	return s.progressSvc.GetUserCourseFull(ctx, meUUID, courseUUID)
}

func (s *Service) GetUserExercise(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID) (UserExercise, *e.APIError) {
	return s.progressSvc.GetUserExercise(ctx, meUUID, exerciseUUID)
}

func (s *Service) SaveUserExerciseSubmission(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID, req SaveUserExerciseSubmissionRequest) *e.APIError {
	return s.progressSvc.SaveUserExerciseSubmission(ctx, meUUID, exerciseUUID, req)
}

func (s *Service) RunUserExerciseCodeSubmission(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID, req RunUserExerciseCodeSubmissionRequest) *e.APIError {
	exercise, err := s.q.GetExercise(ctx, db.GetExerciseParams{
		Uuid:     exerciseUUID,
		Language: "en",
	})

	if err != nil {
		return e.NewAPIError(err, ErrGetExerciseFailed)
	}

	if exercise.Type != db.ExerciseTypeCode {
		return e.NewAPIError(errors.New("exercise is not a code exercise"), ErrGetExerciseFailed)
	}

	return nil
}

package exercises

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"context"
	"encoding/json"
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

func (s *Service) Create(ctx context.Context, req CreateExerciseRequest) (db.Exercise, *e.APIError) {
	data, err := json.Marshal(req.Data)
	if err != nil {
		return db.Exercise{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	exercise, err := s.q.CreateExercise(ctx, db.CreateExerciseParams{
		LessonUuid:  req.LessonUuid,
		Name:        req.Name,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
		Reward:      req.Reward,
		Data:        data,
	})

	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "exercises_name_key") {
			return db.Exercise{}, e.NewAPIError(err, ErrExerciseNameAlreadyExists)
		}
		return db.Exercise{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	return exercise, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateExerciseRequest) (db.Exercise, *e.APIError) {
	data, err := json.Marshal(req.Data)
	if err != nil {
		return db.Exercise{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	exercise, err := s.q.UpdateExercise(ctx, db.UpdateExerciseParams{
		Uuid:        id,
		LessonUuid:  req.LessonUuid,
		Name:        req.Name,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
		Reward:      req.Reward,
		Data:        data,
	})

	if err != nil {
		return db.Exercise{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	return exercise, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.DeleteExercise(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrExerciseDeleteFailed)
	}
	return nil
}

func (s *Service) Restore(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.UndeleteExercise(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrExerciseRestoreFailed)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (db.Exercise, *e.APIError) {
	exercise, err := s.q.GetExercise(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Exercise{}, e.NewAPIError(err, ErrExerciseNotFound)
		}
		return db.Exercise{}, e.NewAPIError(err, ErrExerciseGetFailed)
	}
	return exercise, nil
}

func (s *Service) List(ctx context.Context, req ListExercisesRequest) ([]db.Exercise, *e.APIError) {
	var exercises []db.Exercise
	var err error

	exercises, err = s.q.ListExercises(ctx, db.ListExercisesParams{
		LessonUuid: req.LessonUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrExerciseListFailed)
	}

	return exercises, nil
}

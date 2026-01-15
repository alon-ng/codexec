package exercises

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"context"
	"encoding/json"
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

func (s *Service) Create(ctx context.Context, req CreateExerciseRequest) (db.ExerciseWithTranslation, *e.APIError) {
	data, err := json.Marshal(req.Data)
	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	tx, err := s.p.Begin(ctx)
	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	exercise, err := qtx.CreateExercise(ctx, db.CreateExerciseParams{
		LessonUuid: req.LessonUuid,
		OrderIndex: req.OrderIndex,
		Reward:     req.Reward,
		Type:       req.Type,
		Data:       data,
	})

	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	translation, err := qtx.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
		ExerciseUuid: exercise.Uuid,
		Language:     req.Language,
		Name:         req.Name,
		Description:  req.Description,
	})
	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	return db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	}, nil
}

func (s *Service) Update(ctx context.Context, req UpdateExerciseRequest) (db.ExerciseWithTranslation, *e.APIError) {
	var data *json.RawMessage
	if req.Data != nil {
		d, err := json.Marshal(req.Data)
		if err != nil {
			return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
		}
		data = (*json.RawMessage)(&d)
	}

	tx, err := s.p.Begin(ctx)
	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	exercise, err := qtx.UpdateExercise(ctx, db.UpdateExerciseParams{
		Uuid:       req.Uuid,
		OrderIndex: req.OrderIndex,
		Reward:     req.Reward,
		Type:       req.Type,
		Data:       data,
	})

	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	translation, err := qtx.UpdateExerciseTranslation(ctx, db.UpdateExerciseTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	return db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	}, nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (db.ExerciseWithTranslation, *e.APIError) {
	exercise, err := s.q.GetExercise(ctx, db.GetExerciseParams{
		Uuid:     id,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseNotFound)
		}
		return db.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseGetFailed)
	}
	return exercise.ToExerciseWithTranslation(), nil
}

func (s *Service) List(ctx context.Context, req ListExercisesRequest) ([]db.ExerciseWithTranslation, *e.APIError) {
	exercises, err := s.q.ListExercises(ctx, db.ListExercisesParams{
		LessonUuid: req.LessonUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
		Language:   req.Language,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrExerciseListFailed)
	}

	exercisesWithTranslation := make([]db.ExerciseWithTranslation, len(exercises))
	for i, exercise := range exercises {
		exercisesWithTranslation[i] = exercise.ToExerciseWithTranslation()
	}

	return exercisesWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddExerciseTranslationRequest) (db.ExerciseTranslation, *e.APIError) {
	translation, err := s.q.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
		ExerciseUuid: req.ExerciseUuid,
		Language:     req.Language,
		Name:         req.Name,
		Description:  req.Description,
	})

	if err != nil {
		return db.ExerciseTranslation{}, e.NewAPIError(err, ErrExerciseAddTranslationFailed)
	}

	return translation, nil
}

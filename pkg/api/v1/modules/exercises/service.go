package exercises

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/models"
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

func (s *Service) Create(ctx context.Context, req CreateExerciseRequest) (models.ExerciseWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	exercise, err := qtx.CreateExercise(ctx, db.CreateExerciseParams{
		LessonUuid: req.LessonUuid,
		OrderIndex: req.OrderIndex,
		Reward:     req.Reward,
		Type:       req.Type,
	})

	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	translation, err := qtx.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
		ExerciseUuid: exercise.Uuid,
		Language:     req.Language,
		Name:         req.Name,
		Description:  req.Description,
	})
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	exerciseWithTranslation, err := models.ToExerciseWithTranslation(db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	})
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	return exerciseWithTranslation, nil
}

func (s *Service) Update(ctx context.Context, req UpdateExerciseRequest) (models.ExerciseWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	exercise, err := qtx.UpdateExercise(ctx, db.UpdateExerciseParams{
		Uuid:       req.Uuid,
		OrderIndex: req.OrderIndex,
		Reward:     req.Reward,
		Type:       req.Type,
	})

	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	translation, err := qtx.UpdateExerciseTranslation(ctx, db.UpdateExerciseTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	exerciseWithTranslation, err := models.ToExerciseWithTranslation(db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	})
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	return exerciseWithTranslation, nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (models.ExerciseWithTranslation, *e.APIError) {
	exercise, err := s.q.GetExercise(ctx, db.GetExerciseParams{
		Uuid:     id,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseNotFound)
		}
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseGetFailed)
	}
	exerciseWithTranslation, err := models.ToExerciseWithTranslation(exercise.ToExerciseWithTranslation())
	if err != nil {
		return models.ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseGetFailed)
	}
	return exerciseWithTranslation, nil
}

func (s *Service) List(ctx context.Context, req ListExercisesRequest) ([]models.ExerciseWithTranslation, *e.APIError) {
	exercises, err := s.q.ListExercises(ctx, db.ListExercisesParams{
		LessonUuid: req.LessonUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
		Language:   req.Language,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrExerciseListFailed)
	}

	exercisesWithTranslation := make([]models.ExerciseWithTranslation, len(exercises))
	for i, exercise := range exercises {
		exerciseWithTranslation, err := models.ToExerciseWithTranslation(exercise.ToExerciseWithTranslation())
		if err != nil {
			return nil, e.NewAPIError(err, ErrExerciseListFailed)
		}
		exercisesWithTranslation[i] = exerciseWithTranslation
	}

	return exercisesWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddExerciseTranslationRequest) (models.ExerciseTranslation, *e.APIError) {
	translation, err := s.q.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
		ExerciseUuid: req.ExerciseUuid,
		Language:     req.Language,
		Name:         req.Name,
		Description:  req.Description,
	})

	if err != nil {
		return models.ExerciseTranslation{}, e.NewAPIError(err, ErrExerciseAddTranslationFailed)
	}

	exerciseTranslation, err := models.ToExerciseTranslation(translation)
	if err != nil {
		return models.ExerciseTranslation{}, e.NewAPIError(err, ErrExerciseAddTranslationFailed)
	}
	return exerciseTranslation, nil
}

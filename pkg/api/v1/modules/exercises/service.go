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

// Conversion functions
func ToExercise(d db.Exercise) (Exercise, error) {
	var codeData *ExerciseCodeData
	var quizData *ExerciseQuizData
	if d.CodeData != nil {
		err := json.Unmarshal(*d.CodeData, &codeData)
		if err != nil {
			return Exercise{}, err
		}
	}
	if d.QuizData != nil {
		err := json.Unmarshal(*d.QuizData, &quizData)
		if err != nil {
			return Exercise{}, err
		}
	}

	return Exercise{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		LessonUuid: d.LessonUuid,
		OrderIndex: d.OrderIndex,
		Reward:     d.Reward,
		Type:       d.Type,
		CodeData:   codeData,
		QuizData:   quizData,
	}, nil
}

func ToExerciseTranslation(d db.ExerciseTranslation) (ExerciseTranslation, error) {
	var codeData *ExerciseTranslationCodeData
	var quizData *ExerciseTranslationQuizData

	if d.CodeData != nil {
		err := json.Unmarshal(*d.CodeData, &codeData)
		if err != nil {
			return ExerciseTranslation{}, err
		}
	}
	if d.QuizData != nil {
		err := json.Unmarshal(*d.QuizData, &quizData)
		if err != nil {
			return ExerciseTranslation{}, err
		}
	}

	return ExerciseTranslation{
		Uuid:         d.Uuid,
		ExerciseUuid: d.ExerciseUuid,
		Language:     d.Language,
		Name:         d.Name,
		Description:  d.Description,
		CodeData:     codeData,
		QuizData:     quizData,
	}, nil
}

func ToExerciseWithTranslation(d db.ExerciseWithTranslation) (ExerciseWithTranslation, error) {
	exercise, err := ToExercise(d.Exercise)
	if err != nil {
		return ExerciseWithTranslation{}, err
	}

	translation, err := ToExerciseTranslation(d.Translation)
	if err != nil {
		return ExerciseWithTranslation{}, err
	}

	return ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	}, nil
}

func (s *Service) Create(ctx context.Context, req CreateExerciseRequest) (ExerciseWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
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
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	translation, err := qtx.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
		ExerciseUuid: exercise.Uuid,
		Language:     req.Language,
		Name:         req.Name,
		Description:  req.Description,
	})
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	exerciseWithTranslation, err := ToExerciseWithTranslation(db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	})
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseCreationFailed)
	}

	return exerciseWithTranslation, nil
}

func (s *Service) Update(ctx context.Context, req UpdateExerciseRequest) (ExerciseWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
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
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	translation, err := qtx.UpdateExerciseTranslation(ctx, db.UpdateExerciseTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
	}

	exerciseWithTranslation, err := ToExerciseWithTranslation(db.ExerciseWithTranslation{
		Exercise:    exercise,
		Translation: translation,
	})
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseUpdateFailed)
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

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (ExerciseWithTranslation, *e.APIError) {
	exercise, err := s.q.GetExercise(ctx, db.GetExerciseParams{
		Uuid:     id,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseNotFound)
		}
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseGetFailed)
	}
	exerciseWithTranslation, err := ToExerciseWithTranslation(exercise.ToExerciseWithTranslation())
	if err != nil {
		return ExerciseWithTranslation{}, e.NewAPIError(err, ErrExerciseGetFailed)
	}
	return exerciseWithTranslation, nil
}

func (s *Service) List(ctx context.Context, req ListExercisesRequest) ([]ExerciseWithTranslation, *e.APIError) {
	exercises, err := s.q.ListExercises(ctx, db.ListExercisesParams{
		LessonUuid: req.LessonUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
		Language:   req.Language,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrExerciseListFailed)
	}

	exercisesWithTranslation := make([]ExerciseWithTranslation, len(exercises))
	for i, exercise := range exercises {
		exerciseWithTranslation, err := ToExerciseWithTranslation(exercise.ToExerciseWithTranslation())
		if err != nil {
			return nil, e.NewAPIError(err, ErrExerciseListFailed)
		}
		exercisesWithTranslation[i] = exerciseWithTranslation
	}

	return exercisesWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddExerciseTranslationRequest) (ExerciseTranslation, *e.APIError) {
	translation, err := s.q.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
		ExerciseUuid: req.ExerciseUuid,
		Language:     req.Language,
		Name:         req.Name,
		Description:  req.Description,
	})

	if err != nil {
		return ExerciseTranslation{}, e.NewAPIError(err, ErrExerciseAddTranslationFailed)
	}

	exerciseTranslation, err := ToExerciseTranslation(translation)
	if err != nil {
		return ExerciseTranslation{}, e.NewAPIError(err, ErrExerciseAddTranslationFailed)
	}
	return exerciseTranslation, nil
}

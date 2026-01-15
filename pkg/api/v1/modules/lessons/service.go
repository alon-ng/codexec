package lessons

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

func (s *Service) Create(ctx context.Context, req CreateLessonRequest) (db.LessonWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	lesson, err := qtx.CreateLesson(ctx, db.CreateLessonParams{
		CourseUuid: req.CourseUuid,
		OrderIndex: req.OrderIndex,
		IsPublic:   req.IsPublic,
	})

	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	translation, err := qtx.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
		LessonUuid:  lesson.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	return db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	}, nil
}

func (s *Service) Update(ctx context.Context, req UpdateLessonRequest) (db.LessonWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	lesson, err := qtx.UpdateLesson(ctx, db.UpdateLessonParams{
		Uuid:       req.Uuid,
		OrderIndex: req.OrderIndex,
		IsPublic:   req.IsPublic,
	})

	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	translation, err := qtx.UpdateLessonTranslation(ctx, db.UpdateLessonTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	return db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.DeleteLesson(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrLessonDeleteFailed)
	}
	return nil
}

func (s *Service) Restore(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.UndeleteLesson(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrLessonRestoreFailed)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (db.LessonWithTranslation, *e.APIError) {
	lesson, err := s.q.GetLesson(ctx, db.GetLessonParams{
		Uuid:     id,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonNotFound)
		}
		return db.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonGetFailed)
	}
	return lesson.ToLessonWithTranslation(), nil
}

func (s *Service) List(ctx context.Context, req ListLessonsRequest) ([]db.LessonWithTranslation, *e.APIError) {
	lessons, err := s.q.ListLessons(ctx, db.ListLessonsParams{
		CourseUuid: req.CourseUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrLessonListFailed)
	}

	lessonsWithTranslation := make([]db.LessonWithTranslation, len(lessons))
	for i, lesson := range lessons {
		lessonsWithTranslation[i] = lesson.ToLessonWithTranslation()
	}

	return lessonsWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddLessonTranslationRequest) (db.LessonTranslation, *e.APIError) {
	translation, err := s.q.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
		LessonUuid:  req.LessonUuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return db.LessonTranslation{}, e.NewAPIError(err, ErrLessonAddTranslationFailed)
	}

	return translation, nil
}

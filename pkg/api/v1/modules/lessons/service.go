package lessons

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

func (s *Service) Create(ctx context.Context, req CreateLessonRequest) (models.LessonWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	lesson, err := qtx.CreateLesson(ctx, db.CreateLessonParams{
		CourseUuid: req.CourseUuid,
		OrderIndex: req.OrderIndex,
		IsPublic:   req.IsPublic,
	})

	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	translation, err := qtx.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
		LessonUuid:  lesson.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
	})
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	lessonWithTranslation, err := models.ToLessonWithTranslation(db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	})
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}
	return lessonWithTranslation, nil
}

func (s *Service) Update(ctx context.Context, req UpdateLessonRequest) (models.LessonWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	lesson, err := qtx.UpdateLesson(ctx, db.UpdateLessonParams{
		Uuid:       req.Uuid,
		OrderIndex: req.OrderIndex,
		IsPublic:   req.IsPublic,
	})

	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	translation, err := qtx.UpdateLessonTranslation(ctx, db.UpdateLessonTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
	})
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	lessonWithTranslation, err := models.ToLessonWithTranslation(db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	})
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}
	return lessonWithTranslation, nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (models.LessonWithTranslation, *e.APIError) {
	lesson, err := s.q.GetLesson(ctx, db.GetLessonParams{
		Uuid:     id,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonNotFound)
		}
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonGetFailed)
	}
	lessonWithTranslation, err := models.ToLessonWithTranslation(lesson.ToLessonWithTranslation())
	if err != nil {
		return models.LessonWithTranslation{}, e.NewAPIError(err, ErrLessonGetFailed)
	}
	return lessonWithTranslation, nil
}

func (s *Service) List(ctx context.Context, req ListLessonsRequest) ([]models.LessonWithTranslation, *e.APIError) {
	lessons, err := s.q.ListLessons(ctx, db.ListLessonsParams{
		CourseUuid: req.CourseUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrLessonListFailed)
	}

	lessonsWithTranslation := make([]models.LessonWithTranslation, len(lessons))
	for i, lesson := range lessons {
		lessonWithTranslation, err := models.ToLessonWithTranslation(lesson.ToLessonWithTranslation())
		if err != nil {
			return nil, e.NewAPIError(err, ErrLessonListFailed)
		}
		lessonsWithTranslation[i] = lessonWithTranslation
	}

	return lessonsWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddLessonTranslationRequest) (models.LessonTranslation, *e.APIError) {
	translation, err := s.q.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
		LessonUuid:  req.LessonUuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
	})

	if err != nil {
		return models.LessonTranslation{}, e.NewAPIError(err, ErrLessonAddTranslationFailed)
	}

	lessonTranslation, err := models.ToLessonTranslation(translation)
	if err != nil {
		return models.LessonTranslation{}, e.NewAPIError(err, ErrLessonAddTranslationFailed)
	}
	return lessonTranslation, nil
}

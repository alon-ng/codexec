package lessons

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

func (s *Service) Create(ctx context.Context, req CreateLessonRequest) (db.Lesson, *e.APIError) {
	lesson, err := s.q.CreateLesson(ctx, db.CreateLessonParams{
		CourseUuid:  req.CourseUuid,
		Name:        req.Name,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
		IsPublic:    req.IsPublic,
	})

	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "lessons_name_key") {
			return db.Lesson{}, e.NewAPIError(err, ErrLessonNameAlreadyExists)
		}
		return db.Lesson{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	return lesson, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateLessonRequest) (db.Lesson, *e.APIError) {
	lesson, err := s.q.UpdateLesson(ctx, db.UpdateLessonParams{
		Uuid:        id,
		CourseUuid:  req.CourseUuid,
		Name:        req.Name,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
		IsPublic:    req.IsPublic,
	})

	if err != nil {
		return db.Lesson{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	return lesson, nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID) (db.Lesson, *e.APIError) {
	lesson, err := s.q.GetLesson(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Lesson{}, e.NewAPIError(err, ErrLessonNotFound)
		}
		return db.Lesson{}, e.NewAPIError(err, ErrLessonGetFailed)
	}
	return lesson, nil
}

func (s *Service) List(ctx context.Context, req ListLessonsRequest) ([]db.Lesson, *e.APIError) {
	var lessons []db.Lesson
	var err error

	lessons, err = s.q.ListLessons(ctx, db.ListLessonsParams{
		CourseUuid: req.CourseUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrLessonListFailed)
	}

	return lessons, nil
}

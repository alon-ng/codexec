package courses

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

func (s *Service) Create(ctx context.Context, req CreateCourseRequest) (db.Course, *e.APIError) {
	course, err := s.q.CreateCourse(ctx, db.CreateCourseParams{
		Name:        req.Name,
		Description: req.Description,
		Subject:     req.Subject,
		Price:       req.Price,
		Discount:    req.Discount,
		IsActive:    req.IsActive,
		Difficulty:  req.Difficulty,
		Bullets:     req.Bullets,
	})

	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "courses_name_key") {
			return db.Course{}, e.NewAPIError(err, ErrCourseNameAlreadyExists)
		}
		return db.Course{}, e.NewAPIError(err, ErrCourseCreationFailed)
	}

	return course, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateCourseRequest) (db.Course, *e.APIError) {
	course, err := s.q.UpdateCourse(ctx, db.UpdateCourseParams{
		Uuid:        id,
		Name:        req.Name,
		Description: req.Description,
		Subject:     req.Subject,
		Price:       req.Price,
		Discount:    req.Discount,
		IsActive:    req.IsActive,
		Difficulty:  req.Difficulty,
		Bullets:     req.Bullets,
	})

	if err != nil {
		return db.Course{}, e.NewAPIError(err, ErrCourseUpdateFailed)
	}

	return course, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.DeleteCourse(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrCourseDeleteFailed)
	}
	return nil
}

func (s *Service) Restore(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.UndeleteCourse(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrCourseRestoreFailed)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (db.CourseFull, *e.APIError) {
	course, err := s.q.GetCourseFull(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.CourseFull{}, e.NewAPIError(err, ErrCourseNotFound)
		}
		return db.CourseFull{}, e.NewAPIError(err, ErrCourseGetFailed)
	}
	return course, nil
}

func (s *Service) List(ctx context.Context, req ListCoursesRequest) ([]db.Course, *e.APIError) {
	var courses []db.Course
	var err error

	courses, err = s.q.ListCourses(ctx, db.ListCoursesParams{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Subject:  req.Subject,
		IsActive: req.IsActive,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrCourseListFailed)
	}

	return courses, nil
}

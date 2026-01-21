package courses

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/lessons"
	"codim/pkg/db"
	"context"
	"errors"
	"fmt"

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
func ToCourse(d db.Course) (Course, error) {
	return Course{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		Subject:    d.Subject,
		Price:      d.Price,
		Discount:   d.Discount,
		IsActive:   d.IsActive,
		Difficulty: d.Difficulty,
	}, nil
}

func ToCourseTranslation(d db.CourseTranslation) (CourseTranslation, error) {
	return CourseTranslation{
		Uuid:        d.Uuid,
		CourseUuid:  d.CourseUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
		Bullets:     d.Bullets,
	}, nil
}

func ToCourseWithTranslation(d db.CourseWithTranslation) (CourseWithTranslation, error) {
	course, err := ToCourse(d.Course)
	if err != nil {
		return CourseWithTranslation{}, err
	}

	translation, err := ToCourseTranslation(d.Translation)
	if err != nil {
		return CourseWithTranslation{}, err
	}

	return CourseWithTranslation{
		Course:      course,
		Translation: translation,
	}, nil
}

func ToCourseFull(d db.CourseFull) (CourseFull, error) {
	courseWithTranslation, err := ToCourseWithTranslation(d.CourseWithTranslation)
	if err != nil {
		return CourseFull{}, err
	}

	lessonsList := make([]lessons.LessonFull, len(d.Lessons))
	for i, lesson := range d.Lessons {
		lessonFull, err := lessons.ToLessonFull(lesson)
		if err != nil {
			return CourseFull{}, err
		}
		lessonsList[i] = lessonFull
	}
	return CourseFull{
		CourseWithTranslation: courseWithTranslation,
		Lessons:               lessonsList,
	}, nil
}

func (s *Service) Create(ctx context.Context, req CreateCourseRequest) (CourseWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseCreationFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	course, err := qtx.CreateCourse(ctx, db.CreateCourseParams{
		Subject:    req.Subject,
		Price:      req.Price,
		Discount:   req.Discount,
		IsActive:   req.IsActive,
		Difficulty: req.Difficulty,
	})

	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseCreationFailed)
	}

	translation, err := qtx.CreateCourseTranslation(ctx, db.CreateCourseTranslationParams{
		CourseUuid:  course.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
		Bullets:     req.Bullets,
	})

	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseCreationFailed)
	}

	courseWithTranslation, err := ToCourseWithTranslation(db.CourseWithTranslation{
		Course:      course,
		Translation: translation,
	})
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseCreationFailed)
	}
	return courseWithTranslation, nil
}

func (s *Service) Update(ctx context.Context, req UpdateCourseRequest) (CourseWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseUpdateFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	course, err := qtx.UpdateCourse(ctx, db.UpdateCourseParams{
		Uuid:       req.Uuid,
		Subject:    req.Subject,
		Price:      req.Price,
		Discount:   req.Discount,
		IsActive:   req.IsActive,
		Difficulty: req.Difficulty,
	})

	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseUpdateFailed)
	}

	translation, err := qtx.UpdateCourseTranslation(ctx, db.UpdateCourseTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
		Bullets:     req.Bullets,
	})
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseUpdateFailed)
	}

	courseWithTranslation, err := ToCourseWithTranslation(db.CourseWithTranslation{
		Course:      course,
		Translation: translation,
	})
	if err != nil {
		return CourseWithTranslation{}, e.NewAPIError(err, ErrCourseUpdateFailed)
	}
	return courseWithTranslation, nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (CourseFull, *e.APIError) {
	course, err := s.q.GetCourseFull(ctx, id, language)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CourseFull{}, e.NewAPIError(err, ErrCourseNotFound)
		}
		return CourseFull{}, e.NewAPIError(err, ErrCourseGetFailed)
	}
	courseFull, err := ToCourseFull(course)
	if err != nil {
		return CourseFull{}, e.NewAPIError(err, ErrCourseGetFailed)
	}
	return courseFull, nil
}

func (s *Service) List(ctx context.Context, req ListCoursesRequest) ([]CourseWithTranslation, *e.APIError) {
	courses, err := s.q.ListCourses(ctx, db.ListCoursesParams{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Subject:  req.Subject,
		IsActive: req.IsActive,
		Language: req.Language,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrCourseListFailed)
	}

	coursesWithTranslation := make([]CourseWithTranslation, len(courses))
	for i, course := range courses {
		courseWithTranslation, err := ToCourseWithTranslation(course.ToCourseWithTranslation())
		if err != nil {
			return nil, e.NewAPIError(err, ErrCourseListFailed)
		}
		coursesWithTranslation[i] = courseWithTranslation
	}

	return coursesWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddCourseTranslationRequest) (CourseTranslation, *e.APIError) {
	translation, err := s.q.CreateCourseTranslation(ctx, db.CreateCourseTranslationParams{
		CourseUuid:  req.CourseUuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
		Bullets:     req.Bullets,
	})

	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "uq_course_translations_course_language") {
			return CourseTranslation{}, e.NewAPIError(err, fmt.Sprintf(ErrCourseTranslationAlreadyExists, req.Language))
		}
		return CourseTranslation{}, e.NewAPIError(err, ErrCourseAddTranslationFailed)
	}

	courseTranslation, err := ToCourseTranslation(translation)
	if err != nil {
		return CourseTranslation{}, e.NewAPIError(err, ErrCourseAddTranslationFailed)
	}
	return courseTranslation, nil
}

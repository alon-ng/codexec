package courses

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/exercises"
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
func toCourse(d db.Course) Course {
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
	}
}

func toCourseTranslation(d db.CourseTranslation) CourseTranslation {
	return CourseTranslation{
		Uuid:        d.Uuid,
		CourseUuid:  d.CourseUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
		Bullets:     d.Bullets,
	}
}

func toCourseWithTranslation(d db.CourseWithTranslation) CourseWithTranslation {
	return CourseWithTranslation{
		Course:      toCourse(d.Course),
		Translation: toCourseTranslation(d.Translation),
	}
}

func toCourseFull(d db.CourseFull) CourseFull {
	lessons := make([]lessons.LessonFull, len(d.Lessons))
	for i, lesson := range d.Lessons {
		lessons[i] = toLessonFull(lesson)
	}
	return CourseFull{
		CourseWithTranslation: toCourseWithTranslation(d.CourseWithTranslation),
		Lessons:                lessons,
	}
}

func toLessonFull(d db.LessonFull) lessons.LessonFull {
	exerciseList := make([]exercises.ExerciseWithTranslation, len(d.Exercises))
	for i, exercise := range d.Exercises {
		exerciseList[i] = toExerciseWithTranslation(exercise)
	}
	return lessons.LessonFull{
		LessonWithTranslation: toLessonWithTranslation(d.LessonWithTranslation),
		Exercises:              exerciseList,
	}
}

func toLessonWithTranslation(d db.LessonWithTranslation) lessons.LessonWithTranslation {
	return lessons.LessonWithTranslation{
		Lesson:      toLesson(d.Lesson),
		Translation: toLessonTranslation(d.Translation),
	}
}

func toLesson(d db.Lesson) lessons.Lesson {
	return lessons.Lesson{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		CourseUuid: d.CourseUuid,
		OrderIndex: d.OrderIndex,
		IsPublic:   d.IsPublic,
	}
}

func toLessonTranslation(d db.LessonTranslation) lessons.LessonTranslation {
	return lessons.LessonTranslation{
		Uuid:        d.Uuid,
		LessonUuid:  d.LessonUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
	}
}

func toExerciseWithTranslation(d db.ExerciseWithTranslation) exercises.ExerciseWithTranslation {
	return exercises.ExerciseWithTranslation{
		Exercise:    toExercise(d.Exercise),
		Translation: toExerciseTranslation(d.Translation),
	}
}

func toExercise(d db.Exercise) exercises.Exercise {
	return exercises.Exercise{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		LessonUuid: d.LessonUuid,
		OrderIndex: d.OrderIndex,
		Reward:     d.Reward,
		Type:       d.Type,
		Data:       d.Data,
	}
}

func toExerciseTranslation(d db.ExerciseTranslation) exercises.ExerciseTranslation {
	return exercises.ExerciseTranslation{
		Uuid:         d.Uuid,
		ExerciseUuid: d.ExerciseUuid,
		Language:     d.Language,
		Name:         d.Name,
		Description:  d.Description,
	}
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

	return toCourseWithTranslation(db.CourseWithTranslation{
		Course:      course,
		Translation: translation,
	}), nil
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

	return toCourseWithTranslation(db.CourseWithTranslation{
		Course:      course,
		Translation: translation,
	}), nil
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
	return toCourseFull(course), nil
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
		coursesWithTranslation[i] = toCourseWithTranslation(course.ToCourseWithTranslation())
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

	return toCourseTranslation(translation), nil
}

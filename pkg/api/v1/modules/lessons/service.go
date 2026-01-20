package lessons

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/exercises"
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

// Conversion functions
func toLesson(d db.Lesson) Lesson {
	return Lesson{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		CourseUuid: d.CourseUuid,
		OrderIndex: d.OrderIndex,
		IsPublic:   d.IsPublic,
	}
}

func toLessonTranslation(d db.LessonTranslation) LessonTranslation {
	return LessonTranslation{
		Uuid:        d.Uuid,
		LessonUuid:  d.LessonUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
	}
}

func toLessonWithTranslation(d db.LessonWithTranslation) LessonWithTranslation {
	return LessonWithTranslation{
		Lesson:      toLesson(d.Lesson),
		Translation: toLessonTranslation(d.Translation),
	}
}

func toLessonFull(d db.LessonFull) LessonFull {
	exerciseList := make([]exercises.ExerciseWithTranslation, len(d.Exercises))
	for i, exercise := range d.Exercises {
		exerciseList[i] = toExerciseWithTranslation(exercise)
	}
	return LessonFull{
		LessonWithTranslation: toLessonWithTranslation(d.LessonWithTranslation),
		Exercises:             exerciseList,
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

func (s *Service) Create(ctx context.Context, req CreateLessonRequest) (LessonWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	lesson, err := qtx.CreateLesson(ctx, db.CreateLessonParams{
		CourseUuid: req.CourseUuid,
		OrderIndex: req.OrderIndex,
		IsPublic:   req.IsPublic,
	})

	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	translation, err := qtx.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
		LessonUuid:  lesson.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonTranslationCreationFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonCreationFailed)
	}

	return toLessonWithTranslation(db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	}), nil
}

func (s *Service) Update(ctx context.Context, req UpdateLessonRequest) (LessonWithTranslation, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	lesson, err := qtx.UpdateLesson(ctx, db.UpdateLessonParams{
		Uuid:       req.Uuid,
		OrderIndex: req.OrderIndex,
		IsPublic:   req.IsPublic,
	})

	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	translation, err := qtx.UpdateLessonTranslation(ctx, db.UpdateLessonTranslationParams{
		Uuid:        req.Uuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonTranslationUpdateFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonUpdateFailed)
	}

	return toLessonWithTranslation(db.LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	}), nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID, language string) (LessonWithTranslation, *e.APIError) {
	lesson, err := s.q.GetLesson(ctx, db.GetLessonParams{
		Uuid:     id,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonNotFound)
		}
		return LessonWithTranslation{}, e.NewAPIError(err, ErrLessonGetFailed)
	}
	return toLessonWithTranslation(lesson.ToLessonWithTranslation()), nil
}

func (s *Service) List(ctx context.Context, req ListLessonsRequest) ([]LessonWithTranslation, *e.APIError) {
	lessons, err := s.q.ListLessons(ctx, db.ListLessonsParams{
		CourseUuid: req.CourseUuid,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrLessonListFailed)
	}

	lessonsWithTranslation := make([]LessonWithTranslation, len(lessons))
	for i, lesson := range lessons {
		lessonsWithTranslation[i] = toLessonWithTranslation(lesson.ToLessonWithTranslation())
	}

	return lessonsWithTranslation, nil
}

func (s *Service) AddTranslation(ctx context.Context, req AddLessonTranslationRequest) (LessonTranslation, *e.APIError) {
	translation, err := s.q.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
		LessonUuid:  req.LessonUuid,
		Language:    req.Language,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return LessonTranslation{}, e.NewAPIError(err, ErrLessonAddTranslationFailed)
	}

	return toLessonTranslation(translation), nil
}

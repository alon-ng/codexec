package me

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/courses"
	"codim/pkg/api/v1/modules/users"
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
func toUser(d db.User) users.User {
	return users.User{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		FirstName:  d.FirstName,
		LastName:   d.LastName,
		Email:      d.Email,
		IsVerified: d.IsVerified,
		Streak:     d.Streak,
		Score:      d.Score,
		IsAdmin:    d.IsAdmin,
	}
}

func toUserCourseWithProgress(d db.UserCourseWithProgress) UserCourseWithProgress {
	return UserCourseWithProgress{
		CourseWithTranslation: toCourseWithTranslation(d.CourseWithTranslation),
		UserCourseStartedAt:   d.UserCourseStartedAt,
		UserCourseLastAccessedAt: d.UserCourseLastAccessedAt,
		UserCourseCompletedAt:    d.UserCourseCompletedAt,
		TotalExercises:           d.TotalExercises,
		CompletedExercises:       d.CompletedExercises,
		NextLessonUuid:           d.NextLessonUuid,
		NextLessonName:           d.NextLessonName,
		NextExerciseUuid:         d.NextExerciseUuid,
		NextExerciseName:         d.NextExerciseName,
	}
}

func toCourseWithTranslation(d db.CourseWithTranslation) courses.CourseWithTranslation {
	return courses.CourseWithTranslation{
		Course:      toCourse(d.Course),
		Translation: toCourseTranslation(d.Translation),
	}
}

func toCourse(d db.Course) courses.Course {
	return courses.Course{
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

func toCourseTranslation(d db.CourseTranslation) courses.CourseTranslation {
	return courses.CourseTranslation{
		Uuid:        d.Uuid,
		CourseUuid:  d.CourseUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
		Bullets:     d.Bullets,
	}
}

func toUserCourseFull(d db.UserCourseFull) UserCourseFull {
	lessons := make([]UserLessonStatus, len(d.Lessons))
	for i, lesson := range d.Lessons {
		lessons[i] = toUserLessonStatus(lesson)
	}
	return UserCourseFull{
		CourseUuid:     d.CourseUuid,
		StartedAt:      d.StartedAt,
		LastAccessedAt: d.LastAccessedAt,
		IsCompleted:    d.IsCompleted,
		CompletedAt:    d.CompletedAt,
		Lessons:        lessons,
	}
}

func toUserLessonStatus(d db.UserLessonStatus) UserLessonStatus {
	exercises := make([]UserExerciseStatus, len(d.Exercises))
	for i, exercise := range d.Exercises {
		exercises[i] = toUserExerciseStatus(exercise)
	}
	return UserLessonStatus{
		LessonUuid:     d.LessonUuid,
		StartedAt:      d.StartedAt,
		LastAccessedAt: d.LastAccessedAt,
		IsCompleted:    d.IsCompleted,
		CompletedAt:    d.CompletedAt,
		Exercises:      exercises,
	}
}

func toUserExerciseStatus(d db.UserExerciseStatus) UserExerciseStatus {
	return UserExerciseStatus{
		ExerciseUuid:   d.ExerciseUuid,
		StartedAt:      d.StartedAt,
		LastAccessedAt: d.LastAccessedAt,
		IsCompleted:    d.IsCompleted,
		CompletedAt:    d.CompletedAt,
	}
}

func toUserExercise(d db.UserExercise) UserExercise {
	return UserExercise{
		Uuid:           d.Uuid,
		StartedAt:      d.StartedAt,
		LastAccessedAt: d.LastAccessedAt,
		UserUuid:       d.UserUuid,
		ExerciseUuid:   d.ExerciseUuid,
		Submission:     d.Submission,
		Attempts:       d.Attempts,
		CompletedAt:    d.CompletedAt,
	}
}

func (s *Service) Me(ctx context.Context, meUUID uuid.UUID) (users.User, *e.APIError) {
	u, err := s.q.GetUser(ctx, meUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return users.User{}, e.NewAPIError(err, ErrMeFailed)
		}

		return users.User{}, e.NewAPIError(err, ErrMeFailed)
	}

	return toUser(u), nil
}

func (s *Service) ListUserCoursesWithProgress(ctx context.Context, meUUID uuid.UUID, req ListUserCoursesWithProgressRequest) ([]UserCourseWithProgress, *e.APIError) {
	userCourses, err := s.q.ListUserCoursesWithProgress(ctx, db.ListUserCoursesWithProgressParams{
		UserUuid: meUUID,
		Language: req.Language,
		Limit:    req.Limit,
		Offset:   req.Offset,
		Subject:  req.Subject,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, e.NewAPIError(err, ErrGetUserCoursesWithProgressFailed)
	}

	userCoursesWithProgress := make([]UserCourseWithProgress, len(userCourses))
	for i, userCourse := range userCourses {
		userCoursesWithProgress[i] = toUserCourseWithProgress(userCourse.ToUserCourseWithProgress())
	}

	return userCoursesWithProgress, nil
}

func (s *Service) GetUserCourseFull(ctx context.Context, meUUID uuid.UUID, courseUUID uuid.UUID) (UserCourseFull, *e.APIError) {
	userCourse, err := s.q.GetUserCourseFull(ctx, meUUID, courseUUID)
	if err != nil {
		return UserCourseFull{}, e.NewAPIError(err, ErrGetUserCourseFullFailed)
	}

	return toUserCourseFull(userCourse), nil
}

func (s *Service) GetUserExercise(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID) (UserExercise, *e.APIError) {
	userExercise, err := s.q.GetUserExercise(ctx, db.GetUserExerciseParams{
		UserUuid:     meUUID,
		ExerciseUuid: exerciseUUID,
	})
	if err != nil {
		return UserExercise{}, e.NewAPIError(err, ErrGetUserExerciseFailed)
	}

	return toUserExercise(userExercise), nil
}

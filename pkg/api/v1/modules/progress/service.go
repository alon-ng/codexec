package progress

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/courses"
	"codim/pkg/db"
	"context"
	"encoding/json"
	"errors"

	v "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewService(q *db.Queries, p *pgxpool.Pool) *Service {
	return &Service{q: q, p: p}
}

func (s *Service) InitUserCourse(ctx context.Context, userUuid uuid.UUID, courseUuid uuid.UUID) *e.APIError {
	_, err := s.q.InitUserCourse(ctx, db.InitUserCourseParams{
		UserUuid:   userUuid,
		CourseUuid: courseUuid,
	})
	if err != nil {
		return e.NewAPIError(err, ErrInitUserCourseFailed)
	}

	return nil
}

func (s *Service) CompleteUserExercise(ctx context.Context, userUuid uuid.UUID, exerciseUuid uuid.UUID) (*uuid.UUID, *uuid.UUID, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return nil, nil, e.NewAPIError(err, ErrCompleteUserExerciseFailed)
	}

	qtx := s.q.WithTx(tx)

	defer tx.Rollback(ctx)

	_, err = qtx.CompleteUserExercise(ctx, db.CompleteUserExerciseParams{
		UserUuid:     userUuid,
		ExerciseUuid: exerciseUuid,
	})
	if err != nil {
		return nil, nil, e.NewAPIError(err, ErrCompleteUserExerciseFailed)
	}

	err = qtx.SyncProgressAfterExercise(ctx, db.SyncProgressAfterExerciseParams{
		UserUuid:     userUuid,
		ExerciseUuid: exerciseUuid,
	})
	if err != nil {
		return nil, nil, e.NewAPIError(err, ErrCompleteUserExerciseFailed)
	}

	courseAndLesson, err := qtx.GetExerciseLessonCourse(ctx, exerciseUuid)
	if err != nil {
		return nil, nil, e.NewAPIError(err, ErrCompleteUserExerciseFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, nil, e.NewAPIError(err, ErrCompleteUserExerciseFailed)
	}

	userCourses, err := s.q.ListUserCoursesWithProgress(ctx, db.ListUserCoursesWithProgressParams{
		UserUuid:   userUuid,
		Language:   "en",
		Limit:      1,
		CourseUuid: &courseAndLesson.CourseUuid,
	})
	if err != nil {
		return nil, nil, e.NewAPIError(err, ErrGetUserCoursesWithProgressFailed)
	}

	if len(userCourses) == 0 {
		return nil, nil, e.NewAPIError(errors.New("user course not found"), ErrCompleteUserExerciseFailed)
	}
	userCourse := userCourses[0]

	return userCourse.NextLessonUuid, userCourse.NextExerciseUuid, nil
}

// Conversion functions
func toUserCourseWithProgress(d db.UserCourseWithProgress) (UserCourseWithProgress, error) {
	courseWithTranslation, err := courses.ToCourseWithTranslation(d.CourseWithTranslation)
	if err != nil {
		return UserCourseWithProgress{}, err
	}

	return UserCourseWithProgress{
		CourseWithTranslation:    courseWithTranslation,
		UserCourseStartedAt:      d.UserCourseStartedAt,
		UserCourseLastAccessedAt: d.UserCourseLastAccessedAt,
		UserCourseCompletedAt:    d.UserCourseCompletedAt,
		TotalExercises:           d.TotalExercises,
		CompletedExercises:       d.CompletedExercises,
		NextLessonUuid:           d.NextLessonUuid,
		NextLessonName:           d.NextLessonName,
		NextExerciseUuid:         d.NextExerciseUuid,
		NextExerciseName:         d.NextExerciseName,
	}, nil
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
		userCourseWithProgress, err := toUserCourseWithProgress(userCourse.ToUserCourseWithProgress())
		if err != nil {
			return nil, e.NewAPIError(err, ErrGetUserCoursesWithProgressFailed)
		}
		userCoursesWithProgress[i] = userCourseWithProgress
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

func (s *Service) SaveUserExerciseSubmission(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID, submission SaveUserExerciseSubmissionRequest) *e.APIError {
	submissionRaw, err := s.ValidateSubmission(submission)
	if err != nil {
		return e.NewAPIError(err, ErrValidateSubmissionFailed)
	}

	_, err = s.q.UpdateUserExerciseSubmission(ctx, db.UpdateUserExerciseSubmissionParams{
		UserUuid:     meUUID,
		ExerciseUuid: exerciseUUID,
		Submission:   submissionRaw,
		Type:         submission.Type,
	})
	if err != nil {
		return e.NewAPIError(err, ErrSaveUserExerciseSubmissionFailed)
	}

	return nil
}

func (s *Service) ValidateSubmission(submission SaveUserExerciseSubmissionRequest) (*json.RawMessage, error) {
	validate := v.New()

	switch submission.Type {
	case db.ExerciseTypeCode:
		var userExerciseSubmissionCode UserExerciseSubmissionCode
		err := json.Unmarshal(submission.Submission, &userExerciseSubmissionCode)
		if err != nil {
			return nil, err
		}

		err = validate.Struct(userExerciseSubmissionCode)
		if err != nil {
			return nil, err
		}
		return &submission.Submission, nil
	case db.ExerciseTypeQuiz:
		return &submission.Submission, nil
	}

	return nil, errors.New("invalid submission type")
}

package progress

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/models"
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

func (s *Service) ListUserCoursesWithProgress(ctx context.Context, meUUID uuid.UUID, req ListUserCoursesWithProgressRequest) ([]models.UserCourseWithProgress, *e.APIError) {
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

	userCoursesWithProgress := make([]models.UserCourseWithProgress, len(userCourses))
	for i, userCourse := range userCourses {
		userCourseWithProgress, err := models.ToUserCourseWithProgress(userCourse.ToUserCourseWithProgress())
		if err != nil {
			return nil, e.NewAPIError(err, ErrGetUserCoursesWithProgressFailed)
		}
		userCoursesWithProgress[i] = userCourseWithProgress
	}

	return userCoursesWithProgress, nil
}

func (s *Service) GetUserCourseFull(ctx context.Context, meUUID uuid.UUID, courseUUID uuid.UUID) (models.UserCourseFull, *e.APIError) {
	userCourse, err := s.q.GetUserCourseFull(ctx, meUUID, courseUUID)
	if err != nil {
		return models.UserCourseFull{}, e.NewAPIError(err, ErrGetUserCourseFullFailed)
	}

	return models.ToUserCourseFull(userCourse), nil
}

func (s *Service) GetUserExercise(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID) (models.UserExercise, *e.APIError) {
	userExercise, err := s.q.GetUserExercise(ctx, db.GetUserExerciseParams{
		UserUuid:     meUUID,
		ExerciseUuid: exerciseUUID,
	})
	if err != nil {
		return models.UserExercise{}, e.NewAPIError(err, ErrGetUserExerciseFailed)
	}

	return models.ToUserExercise(userExercise), nil
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
		var userExerciseSubmissionCode models.UserExerciseSubmissionCode
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

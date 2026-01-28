package models

import (
	"codim/pkg/db"
	execmodels "codim/pkg/executors/drivers/models"
	"codim/pkg/fs"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserExerciseStatus struct {
	ExerciseUuid   uuid.UUID  `json:"exercise_uuid" binding:"required"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	IsCompleted    bool       `json:"is_completed" example:"false"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type UserLessonStatus struct {
	LessonUuid     uuid.UUID            `json:"lesson_uuid" binding:"required"`
	StartedAt      *time.Time           `json:"started_at,omitempty"`
	LastAccessedAt *time.Time           `json:"last_accessed_at,omitempty"`
	IsCompleted    bool                 `json:"is_completed" example:"false"`
	CompletedAt    *time.Time           `json:"completed_at,omitempty"`
	Exercises      []UserExerciseStatus `json:"exercises"`
}

type UserCourseFull struct {
	CourseUuid     uuid.UUID          `json:"course_uuid" binding:"required"`
	StartedAt      *time.Time         `json:"started_at,omitempty"`
	LastAccessedAt *time.Time         `json:"last_accessed_at,omitempty"`
	IsCompleted    bool               `json:"is_completed" example:"false"`
	CompletedAt    *time.Time         `json:"completed_at,omitempty"`
	Lessons        []UserLessonStatus `json:"lessons"`
}

type UserCourseWithProgress struct {
	CourseWithTranslation
	UserCourseStartedAt      time.Time  `json:"user_course_started_at" binding:"required"`
	UserCourseLastAccessedAt *time.Time `json:"user_course_last_accessed_at,omitempty"`
	UserCourseCompletedAt    *time.Time `json:"user_course_completed_at,omitempty"`
	TotalExercises           int32      `json:"total_exercises" binding:"required" example:"10"`
	CompletedExercises       int32      `json:"completed_exercises" binding:"required" example:"5"`
	NextLessonUuid           *uuid.UUID `json:"next_lesson_uuid,omitempty"`
	NextLessonName           *string    `json:"next_lesson_name,omitempty"`
	NextExerciseUuid         *uuid.UUID `json:"next_exercise_uuid,omitempty"`
	NextExerciseName         *string    `json:"next_exercise_name,omitempty"`
}

type UserExercise struct {
	Uuid           uuid.UUID       `json:"uuid" binding:"required"`
	StartedAt      time.Time       `json:"started_at" binding:"required"`
	LastAccessedAt *time.Time      `json:"last_accessed_at,omitempty"`
	UserUuid       uuid.UUID       `json:"user_uuid" binding:"required"`
	ExerciseUuid   uuid.UUID       `json:"exercise_uuid" binding:"required"`
	Submission     json.RawMessage `json:"submission" binding:"required"`
	Attempts       int32           `json:"attempts" binding:"required" example:"0"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty"`
}

type UserExerciseSubmissionCode = fs.Entry

type UserExerciseSubmissionQuiz struct {
	Answers map[string]string `json:"answers"`
	Results map[string]bool   `json:"results,omitempty"`
}

type UserExerciseSubmissionResponse struct {
	execmodels.ExecuteResponse
	Passed           bool       `json:"passed" binding:"required" example:"false"`
	NextLessonUuid   *uuid.UUID `json:"next_lesson_uuid,omitempty"`
	NextExerciseUuid *uuid.UUID `json:"next_exercise_uuid,omitempty"`
	Reward           int32      `json:"reward" binding:"required" example:"10"`
}

func ToUserExerciseStatus(d db.UserExerciseStatus) UserExerciseStatus {
	return UserExerciseStatus{
		ExerciseUuid:   d.ExerciseUuid,
		StartedAt:      d.StartedAt,
		LastAccessedAt: d.LastAccessedAt,
		IsCompleted:    d.IsCompleted,
		CompletedAt:    d.CompletedAt,
	}
}

func ToUserLessonStatus(d db.UserLessonStatus) UserLessonStatus {
	exercises := make([]UserExerciseStatus, len(d.Exercises))
	for i, exercise := range d.Exercises {
		exercises[i] = ToUserExerciseStatus(exercise)
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

func ToUserCourseFull(d db.UserCourseFull) UserCourseFull {
	lessons := make([]UserLessonStatus, len(d.Lessons))
	for i, lesson := range d.Lessons {
		lessons[i] = ToUserLessonStatus(lesson)
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

func ToUserCourseWithProgress(d db.UserCourseWithProgress) (UserCourseWithProgress, error) {
	courseWithTranslation, err := ToCourseWithTranslation(d.CourseWithTranslation)
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

func ToUserExercise(d db.UserExercise) UserExercise {
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

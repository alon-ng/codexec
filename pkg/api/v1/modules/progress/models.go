package progress

import (
	"codim/pkg/api/v1/modules/courses"
	"codim/pkg/db"
	"codim/pkg/executors/drivers/models"
	"codim/pkg/fs"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ListUserCoursesWithProgressRequest struct {
	Subject  *string `json:"subject" form:"subject" example:"Programming" query:"subject"`
	IsActive *bool   `json:"is_active" form:"is_active,default=true" example:"true" query:"is_active"`
	Limit    int32   `json:"limit" form:"limit,default=10" example:"10" query:"limit"`
	Offset   int32   `json:"offset" form:"offset,default=0" example:"0" query:"offset"`
	Language string  `json:"language" form:"language,default=en" example:"en" query:"language"`
}

// Response types
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
	courses.CourseWithTranslation
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

type SaveUserExerciseSubmissionRequest struct {
	Type       db.ExerciseType `json:"type" binding:"required" example:"code"`
	Submission json.RawMessage `json:"submission" binding:"required"`
}

type UserExerciseSubmissionCode fs.Entry
type UserExerciseSubmissionQuiz json.RawMessage

type UserExerciseSubmissionResponse struct {
	models.ExecuteResponse
	Passed           bool       `json:"passed" binding:"required" example:"false"`
	NextLessonUuid   *uuid.UUID `json:"next_lesson_uuid,omitempty"`
	NextExerciseUuid *uuid.UUID `json:"next_exercise_uuid,omitempty"`
	Reward           int32      `json:"reward" binding:"required" example:"10"`
}

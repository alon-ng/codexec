package progress

import (
	"codim/pkg/db"
	"encoding/json"
)

type ListUserCoursesWithProgressRequest struct {
	Subject  *string `json:"subject" form:"subject" example:"Programming" query:"subject"`
	IsActive *bool   `json:"is_active" form:"is_active,default=true" example:"true" query:"is_active"`
	Limit    int32   `json:"limit" form:"limit,default=10" example:"10" query:"limit"`
	Offset   int32   `json:"offset" form:"offset,default=0" example:"0" query:"offset"`
	Language string  `json:"language" form:"language,default=en" example:"en" query:"language"`
}

type SaveUserExerciseSubmissionRequest struct {
	Type       db.ExerciseType `json:"type" binding:"required" example:"code"`
	Submission json.RawMessage `json:"submission" binding:"required"`
}

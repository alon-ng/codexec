package lessons

import "github.com/google/uuid"

type CreateLessonRequest struct {
	CourseUuid  uuid.UUID `json:"course_uuid" binding:"required"`
	Name        string    `json:"name" binding:"required" example:"Python Basics"`
	Description string    `json:"description" binding:"required" example:"Learn Python fundamentals"`
	OrderIndex  int16     `json:"order_index" binding:"required" example:"1"`
	IsPublic    bool      `json:"is_public" example:"false"`
}

type UpdateRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
	UpdateLessonRequest
}

type UpdateLessonRequest struct {
	CourseUuid  uuid.UUID `json:"course_uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"Python Basics"`
	Description string    `json:"description" example:"Learn Python fundamentals"`
	OrderIndex  int16     `json:"order_index" example:"1"`
	IsPublic    bool      `json:"is_public" example:"false"`
}

type ListLessonsRequest struct {
	Limit      int32      `json:"limit" form:"limit,default=10" example:"10"`
	Offset     int32      `json:"offset" form:"offset,default=0" example:"0"`
	CourseUuid *uuid.UUID `json:"course_uuid" form:"course_uuid"`
}

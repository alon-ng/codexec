package lessons

import "github.com/google/uuid"

type CreateLessonRequest struct {
	CourseUuid  uuid.UUID `json:"course_uuid" binding:"required"`
	OrderIndex  int16     `json:"order_index" binding:"required" example:"1"`
	IsPublic    bool      `json:"is_public" example:"false"`
	Language    string    `json:"language" binding:"required" example:"en"`
	Name        string    `json:"name" binding:"required" example:"Python Basics"`
	Description string    `json:"description" binding:"required" example:"Learn Python fundamentals"`
}

type UpdateLessonRequest struct {
	Uuid        uuid.UUID `json:"uuid" binding:"required"`
	Language    string    `json:"language" binding:"required" example:"en"`
	OrderIndex  *int16    `json:"order_index" example:"1"`
	IsPublic    *bool     `json:"is_public" example:"false"`
	Name        *string   `json:"name" example:"Python Basics"`
	Description *string   `json:"description" example:"Learn Python fundamentals"`
}

type ListLessonsRequest struct {
	Limit      int32      `json:"limit" form:"limit,default=10" example:"10"`
	Offset     int32      `json:"offset" form:"offset,default=0" example:"0"`
	CourseUuid *uuid.UUID `json:"course_uuid" form:"course_uuid"`
	Language   string     `json:"language" form:"language,default=en" example:"en"`
}

type AddLessonTranslationRequest struct {
	LessonUuid  uuid.UUID `json:"lesson_uuid" binding:"required"`
	Language    string    `json:"language" binding:"required" example:"es"`
	Name        string    `json:"name" binding:"required" example:"Fundamentos de Python"`
	Description string    `json:"description" binding:"required" example:"Aprende los fundamentos de Python"`
}

type IDRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
}

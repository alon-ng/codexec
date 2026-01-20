package courses

import (
	"codim/pkg/api/v1/modules/lessons"
	"time"

	"github.com/google/uuid"
)

type CreateCourseRequest struct {
	Subject     string `json:"subject" binding:"required" example:"Programming"`
	Price       int16  `json:"price" binding:"required" example:"99"`
	Discount    int16  `json:"discount" example:"0"`
	IsActive    bool   `json:"is_active" example:"true"`
	Difficulty  int16  `json:"difficulty" example:"1"`
	Language    string `json:"language" binding:"required" example:"en"`
	Name        string `json:"name" binding:"required" example:"Introduction to Python"`
	Description string `json:"description" binding:"required" example:"Learn Python basics"`
	Bullets     string `json:"bullets" example:"Learn basics\nPractice exercises"`
}

type UpdateCourseRequest struct {
	Uuid        uuid.UUID `json:"uuid" binding:"required"`
	Language    string    `json:"language" binding:"required" example:"en"`
	Subject     *string   `json:"subject" example:"Programming"`
	Price       *int16    `json:"price" example:"99"`
	Discount    *int16    `json:"discount" example:"0"`
	IsActive    *bool     `json:"is_active" example:"true"`
	Difficulty  *int16    `json:"difficulty" example:"1"`
	Name        *string   `json:"name" example:"Introduction to Python"`
	Description *string   `json:"description" example:"Learn Python basics"`
	Bullets     *string   `json:"bullets" example:"Learn basics\nPractice exercises"`
}

type ListCoursesRequest struct {
	Limit    int32   `json:"limit" form:"limit,default=10" example:"10" query:"limit"`
	Offset   int32   `json:"offset" form:"offset,default=0" example:"0" query:"offset"`
	Subject  *string `json:"subject" form:"subject" example:"Programming" query:"subject"`
	IsActive *bool   `json:"is_active" form:"is_active,default=true" example:"true" query:"is_active"`
	Language string  `json:"language" form:"language,default=en" example:"en" query:"language"`
}

type AddCourseTranslationRequest struct {
	CourseUuid  uuid.UUID `json:"course_uuid" binding:"required"`
	Language    string    `json:"language" binding:"required" example:"es"`
	Name        string    `json:"name" binding:"required" example:"Introducción a Python"`
	Description string    `json:"description" binding:"required" example:"Aprende los fundamentos de Python"`
	Bullets     string    `json:"bullets" example:"Aprende lo básico\nPractica ejercicios"`
}

// Response types
type Course struct {
	Uuid       uuid.UUID  `json:"uuid" binding:"required"`
	CreatedAt  time.Time  `json:"created_at" binding:"required"`
	ModifiedAt time.Time  `json:"modified_at" binding:"required"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Subject    string     `json:"subject" binding:"required" example:"Programming"`
	Price      int16      `json:"price" binding:"required" example:"99"`
	Discount   int16      `json:"discount" example:"0"`
	IsActive   bool       `json:"is_active" example:"true"`
	Difficulty int16      `json:"difficulty" example:"1"`
}

type CourseTranslation struct {
	Uuid        uuid.UUID `json:"uuid" binding:"required"`
	CourseUuid  uuid.UUID `json:"course_uuid" binding:"required"`
	Language    string    `json:"language" binding:"required" example:"en"`
	Name        string    `json:"name" binding:"required" example:"Introduction to Python"`
	Description string    `json:"description" binding:"required" example:"Learn Python basics"`
	Bullets     string    `json:"bullets" example:"Learn basics\nPractice exercises"`
}

type CourseWithTranslation struct {
	Course
	Translation CourseTranslation `json:"translation" binding:"required"`
}

type CourseFull struct {
	CourseWithTranslation
	Lessons []lessons.LessonFull `json:"lessons"`
}

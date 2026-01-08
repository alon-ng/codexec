package courses

import "github.com/google/uuid"

type CreateCourseRequest struct {
	Name        string `json:"name" binding:"required" example:"Introduction to Python"`
	Description string `json:"description" binding:"required" example:"Learn Python basics"`
	Subject     string `json:"subject" binding:"required" example:"Programming"`
	Price       int16  `json:"price" binding:"required" example:"99"`
	Discount    int16  `json:"discount" example:"0"`
	IsActive    bool   `json:"is_active" example:"true"`
	Difficulty  int16  `json:"difficulty" example:"1"`
	Bullets     string `json:"bullets" example:"Learn basics\nPractice exercises"`
}

type UpdateRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
	UpdateCourseRequest
}

type UpdateCourseRequest struct {
	Name        string `json:"name" example:"Introduction to Python"`
	Description string `json:"description" example:"Learn Python basics"`
	Subject     string `json:"subject" example:"Programming"`
	Price       int16  `json:"price" example:"99"`
	Discount    int16  `json:"discount" example:"0"`
	IsActive    bool   `json:"is_active" example:"true"`
	Difficulty  int16  `json:"difficulty" example:"1"`
	Bullets     string `json:"bullets" example:"Learn basics\nPractice exercises"`
}

type ListCoursesRequest struct {
	Limit    int32   `json:"limit" form:"limit,default=10" example:"10" query:"limit"`
	Offset   int32   `json:"offset" form:"offset,default=0" example:"0" query:"offset"`
	Subject  *string `json:"subject" form:"subject" example:"Programming" query:"subject"`
	IsActive *bool   `json:"is_active" form:"is_active,default=true" example:"true" query:"is_active"`
}

package models

import (
	"codim/pkg/db"
	"time"

	"github.com/google/uuid"
)

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
	Lessons []LessonFull `json:"lessons"`
}

func ToCourse(d db.Course) (Course, error) {
	return Course{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		Subject:    d.Subject,
		Price:      d.Price,
		Discount:   d.Discount,
		IsActive:   d.IsActive,
		Difficulty: d.Difficulty,
	}, nil
}

func ToCourseTranslation(d db.CourseTranslation) (CourseTranslation, error) {
	return CourseTranslation{
		Uuid:        d.Uuid,
		CourseUuid:  d.CourseUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
		Bullets:     d.Bullets,
	}, nil
}

func ToCourseWithTranslation(d db.CourseWithTranslation) (CourseWithTranslation, error) {
	course, err := ToCourse(d.Course)
	if err != nil {
		return CourseWithTranslation{}, err
	}

	translation, err := ToCourseTranslation(d.Translation)
	if err != nil {
		return CourseWithTranslation{}, err
	}

	return CourseWithTranslation{
		Course:      course,
		Translation: translation,
	}, nil
}

func ToCourseFull(d db.CourseFull) (CourseFull, error) {
	courseWithTranslation, err := ToCourseWithTranslation(d.CourseWithTranslation)
	if err != nil {
		return CourseFull{}, err
	}

	lessonsList := make([]LessonFull, len(d.Lessons))
	for i, lesson := range d.Lessons {
		lessonFull, err := ToLessonFull(lesson)
		if err != nil {
			return CourseFull{}, err
		}
		lessonsList[i] = lessonFull
	}
	return CourseFull{
		CourseWithTranslation: courseWithTranslation,
		Lessons:               lessonsList,
	}, nil
}

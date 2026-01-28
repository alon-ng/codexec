package models

import (
	"codim/pkg/db"
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	Uuid       uuid.UUID  `json:"uuid" binding:"required"`
	CreatedAt  time.Time  `json:"created_at" binding:"required"`
	ModifiedAt time.Time  `json:"modified_at" binding:"required"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	CourseUuid uuid.UUID  `json:"course_uuid" binding:"required"`
	OrderIndex int16      `json:"order_index" binding:"required" example:"1"`
	IsPublic   bool       `json:"is_public" example:"false"`
}

type LessonTranslation struct {
	Uuid        uuid.UUID `json:"uuid" binding:"required"`
	LessonUuid  uuid.UUID `json:"lesson_uuid" binding:"required"`
	Language    string    `json:"language" binding:"required" example:"en"`
	Name        string    `json:"name" binding:"required" example:"Python Basics"`
	Description string    `json:"description" binding:"required" example:"Learn Python fundamentals"`
}

type LessonWithTranslation struct {
	Lesson
	Translation LessonTranslation `json:"translation" binding:"required"`
}

type LessonFull struct {
	LessonWithTranslation
	Exercises []ExerciseWithTranslation `json:"exercises"`
}

func ToLesson(d db.Lesson) (Lesson, error) {
	return Lesson{
		Uuid:       d.Uuid,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		DeletedAt:  d.DeletedAt,
		CourseUuid: d.CourseUuid,
		OrderIndex: d.OrderIndex,
		IsPublic:   d.IsPublic,
	}, nil
}

func ToLessonTranslation(d db.LessonTranslation) (LessonTranslation, error) {
	return LessonTranslation{
		Uuid:        d.Uuid,
		LessonUuid:  d.LessonUuid,
		Language:    d.Language,
		Name:        d.Name,
		Description: d.Description,
	}, nil
}

func ToLessonWithTranslation(d db.LessonWithTranslation) (LessonWithTranslation, error) {
	lesson, err := ToLesson(d.Lesson)
	if err != nil {
		return LessonWithTranslation{}, err
	}

	translation, err := ToLessonTranslation(d.Translation)
	if err != nil {
		return LessonWithTranslation{}, err
	}

	return LessonWithTranslation{
		Lesson:      lesson,
		Translation: translation,
	}, nil
}

func ToLessonFull(d db.LessonFull) (LessonFull, error) {
	lessonWithTranslation, err := ToLessonWithTranslation(d.LessonWithTranslation)
	if err != nil {
		return LessonFull{}, err
	}

	exerciseList := make([]ExerciseWithTranslation, len(d.Exercises))
	for i, exercise := range d.Exercises {
		exerciseWithTranslation, err := ToExerciseWithTranslation(exercise)
		if err != nil {
			return LessonFull{}, err
		}
		exerciseList[i] = exerciseWithTranslation
	}
	return LessonFull{
		LessonWithTranslation: lessonWithTranslation,
		Exercises:             exerciseList,
	}, nil
}

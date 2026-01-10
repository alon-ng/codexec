package db

import (
	"context"

	"github.com/google/uuid"
)

type LessonFull struct {
	Lesson
	Exercises []Exercise `json:"exercises"`
}

type CourseFull struct {
	Course
	Lessons []LessonFull `json:"lessons"`
}

func (q *Queries) GetCourseFull(ctx context.Context, argUuid uuid.UUID) (CourseFull, error) {
	r, err := q.getCourseFull(ctx, argUuid)
	if err != nil {
		return CourseFull{}, err
	}

	lessons := []LessonFull{}
	uqLessonsUUIDs := make(map[uuid.UUID]bool)
	for _, l := range r {
		if l.LessonUuid == nil {
			continue
		}
		if _, ok := uqLessonsUUIDs[*l.LessonUuid]; ok {
			continue
		}

		uqLessonsUUIDs[*l.LessonUuid] = true
		exercises := []Exercise{}
		uqExercisesUUIDs := make(map[uuid.UUID]bool)
		for _, e := range r {
			if e.ExerciseUuid == nil || e.ExerciseLessonUuid == nil || *e.ExerciseLessonUuid != *l.LessonUuid {
				continue
			}
			if _, ok := uqExercisesUUIDs[*e.ExerciseUuid]; ok {
				continue
			}

			uqExercisesUUIDs[*e.ExerciseUuid] = true
			exercises = append(exercises, Exercise{
				Uuid:        *e.ExerciseUuid,
				CreatedAt:   *e.ExerciseCreatedAt,
				ModifiedAt:  *e.ExerciseModifiedAt,
				DeletedAt:   e.ExerciseDeletedAt,
				LessonUuid:  *e.ExerciseLessonUuid,
				Name:        *e.ExerciseName,
				Description: *e.ExerciseDescription,
				OrderIndex:  *e.ExerciseOrderIndex,
				Reward:      *e.ExerciseReward,
				Type:        *e.ExerciseType,
				Data:        *e.ExerciseData,
			})
		}

		lessons = append(lessons, LessonFull{
			Lesson: Lesson{
				Uuid:        *l.LessonUuid,
				CreatedAt:   *l.LessonCreatedAt,
				ModifiedAt:  *l.LessonModifiedAt,
				DeletedAt:   l.LessonDeletedAt,
				CourseUuid:  *l.LessonCourseUuid,
				Name:        *l.LessonName,
				Description: *l.LessonDescription,
				OrderIndex:  *l.LessonOrderIndex,
				IsPublic:    *l.LessonIsPublic,
			},
			Exercises: exercises,
		})
	}

	if len(r) == 0 {
		return CourseFull{}, nil
	}

	return CourseFull{Course: Course{
		Uuid:        r[0].CourseUuid,
		CreatedAt:   r[0].CourseCreatedAt,
		ModifiedAt:  r[0].CourseModifiedAt,
		DeletedAt:   r[0].CourseDeletedAt,
		Name:        r[0].CourseName,
		Description: r[0].CourseDescription,
		Subject:     r[0].CourseSubject,
		Price:       r[0].CoursePrice,
		Discount:    r[0].CourseDiscount,
		IsActive:    r[0].CourseIsActive,
		Difficulty:  r[0].CourseDifficulty,
		Bullets:     r[0].CourseBullets,
	}, Lessons: lessons}, nil
}

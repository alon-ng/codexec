package db

import (
	"context"

	"github.com/google/uuid"
)

type CourseWithTranslation struct {
	Course
	Translation CourseTranslation `json:"translation"`
}

type CourseFull struct {
	CourseWithTranslation
	Lessons []LessonFull `json:"lessons"`
}

func (q *Queries) GetCourseFull(ctx context.Context, u uuid.UUID, language string) (CourseFull, error) {
	r, err := q.getCourseFull(ctx, getCourseFullParams{
		Uuid:     u,
		Language: language,
	})
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
		exercises := []ExerciseWithTranslation{}
		uqExercisesUUIDs := make(map[uuid.UUID]bool)
		for _, e := range r {
			if e.ExerciseUuid == nil || e.ExerciseLessonUuid == nil || *e.ExerciseLessonUuid != *l.LessonUuid {
				continue
			}
			if _, ok := uqExercisesUUIDs[*e.ExerciseUuid]; ok {
				continue
			}

			uqExercisesUUIDs[*e.ExerciseUuid] = true
			exercises = append(exercises, ExerciseWithTranslation{
				Exercise: Exercise{
					Uuid:       *e.ExerciseUuid,
					CreatedAt:  *e.ExerciseCreatedAt,
					ModifiedAt: *e.ExerciseModifiedAt,
					DeletedAt:  e.ExerciseDeletedAt,
					LessonUuid: *e.ExerciseLessonUuid,
					OrderIndex: *e.ExerciseOrderIndex,
					Reward:     *e.ExerciseReward,
					Type:       *e.ExerciseType,
					CodeData:   e.ExerciseCodeData,
					QuizData:   e.ExerciseQuizData,
				},
				Translation: ExerciseTranslation{
					Uuid:         *e.ExerciseTranslationUuid,
					ExerciseUuid: *e.ExerciseUuid,
					Language:     *e.ExerciseTranslationLanguage,
					Name:         *e.ExerciseName,
					Description:  *e.ExerciseDescription,
					CodeData:     e.ExerciseTranslationCodeData,
					QuizData:     e.ExerciseTranslationQuizData,
				},
			})
		}

		lessons = append(lessons, LessonFull{
			LessonWithTranslation: LessonWithTranslation{
				Lesson: Lesson{
					Uuid:       *l.LessonUuid,
					CreatedAt:  *l.LessonCreatedAt,
					ModifiedAt: *l.LessonModifiedAt,
					DeletedAt:  l.LessonDeletedAt,
					CourseUuid: *l.LessonCourseUuid,
					OrderIndex: *l.LessonOrderIndex,
					IsPublic:   *l.LessonIsPublic,
				},
				Translation: LessonTranslation{
					Uuid:        *l.LessonTranslationUuid,
					LessonUuid:  *l.LessonUuid,
					Language:    *l.LessonTranslationLanguage,
					Name:        *l.LessonName,
					Description: *l.LessonDescription,
					Content:     *l.LessonContent,
				},
			},
			Exercises: exercises,
		})
	}

	if len(r) == 0 {
		return CourseFull{}, nil
	}

	return CourseFull{CourseWithTranslation: CourseWithTranslation{
		Course: Course{
			Uuid:       r[0].CourseUuid,
			CreatedAt:  r[0].CourseCreatedAt,
			ModifiedAt: r[0].CourseModifiedAt,
			DeletedAt:  r[0].CourseDeletedAt,
			Subject:    r[0].CourseSubject,
			Price:      r[0].CoursePrice,
			Discount:   r[0].CourseDiscount,
			IsActive:   r[0].CourseIsActive,
			Difficulty: r[0].CourseDifficulty,
		},
		Translation: CourseTranslation{
			Uuid:        r[0].CourseTranslationUuid,
			CourseUuid:  r[0].CourseUuid,
			Language:    r[0].CourseTranslationLanguage,
			Name:        r[0].CourseName,
			Description: r[0].CourseDescription,
			Bullets:     r[0].CourseBullets,
		},
	}, Lessons: lessons}, nil
}

func (l *ListCoursesRow) ToCourseWithTranslation() CourseWithTranslation {
	return CourseWithTranslation{
		Course: Course{
			Uuid:       l.Uuid,
			CreatedAt:  l.CreatedAt,
			ModifiedAt: l.ModifiedAt,
			DeletedAt:  l.DeletedAt,
			Subject:    l.Subject,
			Price:      l.Price,
			Discount:   l.Discount,
			IsActive:   l.IsActive,
			Difficulty: l.Difficulty,
		},
		Translation: CourseTranslation{
			Uuid:        l.Uuid_2,
			CourseUuid:  l.CourseUuid,
			Language:    l.Language,
			Name:        l.Name,
			Description: l.Description,
			Bullets:     l.Bullets,
		},
	}
}

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserExerciseStatus struct {
	ExerciseUuid   uuid.UUID  `json:"exercise_uuid"`
	StartedAt      *time.Time `json:"started_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at"`
	IsCompleted    bool       `json:"is_completed"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type UserLessonStatus struct {
	LessonUuid     uuid.UUID            `json:"lesson_uuid"`
	StartedAt      *time.Time           `json:"started_at"`
	LastAccessedAt *time.Time           `json:"last_accessed_at"`
	IsCompleted    bool                 `json:"is_completed"`
	CompletedAt    *time.Time           `json:"completed_at,omitempty"`
	Exercises      []UserExerciseStatus `json:"exercises"`
}

type UserCourseFull struct {
	CourseUuid     uuid.UUID          `json:"course_uuid"`
	StartedAt      *time.Time         `json:"started_at"`
	LastAccessedAt *time.Time         `json:"last_accessed_at"`
	IsCompleted    bool               `json:"is_completed"`
	CompletedAt    *time.Time         `json:"completed_at,omitempty"`
	Lessons        []UserLessonStatus `json:"lessons"`
}

type UserCourseWithProgress struct {
	CourseWithTranslation
	UserCourseStartedAt      time.Time  `json:"user_course_started_at"`
	UserCourseLastAccessedAt *time.Time `json:"user_course_last_accessed_at"`
	UserCourseCompletedAt    *time.Time `json:"user_course_completed_at"`
	TotalExercises           int32      `json:"total_exercises"`
	CompletedExercises       int32      `json:"completed_exercises"`
	NextLessonUuid           *uuid.UUID `json:"next_lesson_uuid"`
	NextLessonName           *string    `json:"next_lesson_name"`
	NextExerciseUuid         *uuid.UUID `json:"next_exercise_uuid"`
	NextExerciseName         *string    `json:"next_exercise_name"`
}

func (q *Queries) GetUserCourseFull(ctx context.Context, userUuid uuid.UUID, courseUuid uuid.UUID) (UserCourseFull, error) {
	r, err := q.getUserCourseFull(ctx, getUserCourseFullParams{
		UserUuid: userUuid,
		Uuid:     courseUuid,
	})
	if err != nil {
		return UserCourseFull{}, err
	}

	if len(r) == 0 {
		return UserCourseFull{}, nil
	}

	result := UserCourseFull{
		CourseUuid:     r[0].CourseUuid,
		StartedAt:      r[0].UserCourseStartedAt,
		LastAccessedAt: r[0].UserCourseLastAccessedAt,
		IsCompleted:    r[0].CourseIsCompleted,
		CompletedAt:    r[0].CourseCompletedAt,
		Lessons:        []UserLessonStatus{},
	}

	uqLessonsUUIDs := make(map[uuid.UUID]bool)
	for _, row := range r {
		if row.LessonUuid == nil {
			continue
		}
		lessonUuid := *row.LessonUuid
		if _, exists := uqLessonsUUIDs[lessonUuid]; exists {
			continue
		}

		uqLessonsUUIDs[lessonUuid] = true
		exercises := []UserExerciseStatus{}
		uqExercisesUUIDs := make(map[uuid.UUID]bool)

		for _, exerciseRow := range r {
			if exerciseRow.ExerciseUuid == nil ||
				exerciseRow.LessonUuid == nil ||
				*exerciseRow.LessonUuid != lessonUuid {
				continue
			}
			exerciseUuid := *exerciseRow.ExerciseUuid
			if _, exists := uqExercisesUUIDs[exerciseUuid]; exists {
				continue
			}

			uqExercisesUUIDs[exerciseUuid] = true
			exercises = append(exercises, UserExerciseStatus{
				ExerciseUuid:   exerciseUuid,
				StartedAt:      exerciseRow.UserExerciseStartedAt,
				LastAccessedAt: exerciseRow.UserExerciseLastAccessedAt,
				IsCompleted:    exerciseRow.ExerciseIsCompleted,
				CompletedAt:    exerciseRow.ExerciseCompletedAt,
			})
		}

		result.Lessons = append(result.Lessons, UserLessonStatus{
			LessonUuid:     lessonUuid,
			StartedAt:      row.UserLessonStartedAt,
			LastAccessedAt: row.UserLessonLastAccessedAt,
			IsCompleted:    row.LessonIsCompleted,
			CompletedAt:    row.LessonCompletedAt,
			Exercises:      exercises,
		})
	}

	return result, nil
}

func (l *ListUserCoursesWithProgressRow) ToUserCourseWithProgress() UserCourseWithProgress {
	return UserCourseWithProgress{
		CourseWithTranslation: CourseWithTranslation{
			Course: Course{
				Uuid:       l.CourseUuid,
				CreatedAt:  l.CourseCreatedAt,
				ModifiedAt: l.CourseModifiedAt,
				DeletedAt:  l.CourseDeletedAt,
				Subject:    l.CourseSubject,
				Price:      l.CoursePrice,
				Discount:   l.CourseDiscount,
				IsActive:   l.CourseIsActive,
				Difficulty: l.CourseDifficulty,
			},
			Translation: CourseTranslation{
				Uuid:        l.CourseTranslationUuid,
				CourseUuid:  l.CourseUuid,
				Language:    l.CourseTranslationLanguage,
				Name:        l.CourseName,
				Description: l.CourseDescription,
				Bullets:     l.CourseBullets,
			},
		},
		UserCourseStartedAt:      l.UserCourseStartedAt,
		UserCourseLastAccessedAt: l.UserCourseLastAccessedAt,
		UserCourseCompletedAt:    l.UserCourseCompletedAt,
		TotalExercises:           l.TotalExercises,
		CompletedExercises:       l.CompletedExercises,
		NextLessonUuid:           l.NextLessonUuid,
		NextLessonName:           l.NextLessonName,
		NextExerciseUuid:         l.NextExerciseUuid,
		NextExerciseName:         l.NextExerciseName,
	}
}

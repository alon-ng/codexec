package db

type LessonWithTranslation struct {
	Lesson
	Translation LessonTranslation `json:"translation"`
}

type LessonFull struct {
	LessonWithTranslation
	Exercises []ExerciseWithTranslation `json:"exercises"`
}

func (l *GetLessonRow) ToLessonWithTranslation() LessonWithTranslation {
	return LessonWithTranslation{
		Lesson: Lesson{
			Uuid:       l.Uuid,
			CreatedAt:  l.CreatedAt,
			ModifiedAt: l.ModifiedAt,
			DeletedAt:  l.DeletedAt,
			CourseUuid: l.CourseUuid,
			OrderIndex: l.OrderIndex,
			IsPublic:   l.IsPublic,
		},
		Translation: LessonTranslation{
			Uuid:        l.Uuid_2,
			LessonUuid:  l.LessonUuid,
			Language:    l.Language,
			Name:        l.Name,
			Description: l.Description,
			Content:     l.Content,
		},
	}
}

func (l *ListLessonsRow) ToLessonWithTranslation() LessonWithTranslation {
	return LessonWithTranslation{
		Lesson: Lesson{
			Uuid:       l.Uuid,
			CreatedAt:  l.CreatedAt,
			ModifiedAt: l.ModifiedAt,
			DeletedAt:  l.DeletedAt,
			CourseUuid: l.CourseUuid,
			OrderIndex: l.OrderIndex,
			IsPublic:   l.IsPublic,
		},
		Translation: LessonTranslation{
			Uuid:        l.Uuid_2,
			LessonUuid:  l.LessonUuid,
			Language:    l.Language,
			Name:        l.Name,
			Description: l.Description,
			Content:     l.Content,
		},
	}
}

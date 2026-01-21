package db

type ExerciseWithTranslation struct {
	Exercise
	Translation ExerciseTranslation `json:"translation"`
}

func (g *GetExerciseRow) ToExerciseWithTranslation() ExerciseWithTranslation {
	return ExerciseWithTranslation{
		Exercise: Exercise{
			Uuid:       g.Uuid,
			CreatedAt:  g.CreatedAt,
			ModifiedAt: g.ModifiedAt,
			DeletedAt:  g.DeletedAt,
			LessonUuid: g.LessonUuid,
			OrderIndex: g.OrderIndex,
			Reward:     g.Reward,
			Type:       g.Type,
			CodeData:   g.CodeData,
			QuizData:   g.QuizData,
		},
		Translation: ExerciseTranslation{
			Uuid:         g.Uuid_2,
			ExerciseUuid: g.ExerciseUuid,
			Language:     g.Language,
			Name:         g.Name,
			Description:  g.Description,
			CodeData:     g.CodeData_2,
			QuizData:     g.QuizData_2,
		},
	}
}

func (l *ListExercisesRow) ToExerciseWithTranslation() ExerciseWithTranslation {
	return ExerciseWithTranslation{
		Exercise: Exercise{
			Uuid:       l.Uuid,
			CreatedAt:  l.CreatedAt,
			ModifiedAt: l.ModifiedAt,
			DeletedAt:  l.DeletedAt,
			LessonUuid: l.LessonUuid,
			OrderIndex: l.OrderIndex,
			Reward:     l.Reward,
			Type:       l.Type,
			CodeData:   l.CodeData,
			QuizData:   l.QuizData,
		},
		Translation: ExerciseTranslation{
			Uuid:         l.Uuid_2,
			ExerciseUuid: l.ExerciseUuid,
			Language:     l.Language,
			Name:         l.Name,
			Description:  l.Description,
			CodeData:     l.CodeData_2,
			QuizData:     l.QuizData_2,
		},
	}
}

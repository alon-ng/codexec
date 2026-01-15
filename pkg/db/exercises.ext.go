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
			Data:       g.Data,
		},
		Translation: ExerciseTranslation{
			Uuid:         g.Uuid_2,
			ExerciseUuid: g.ExerciseUuid,
			Language:     g.Language,
			Name:         g.Name,
			Description:  g.Description,
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
			Data:       l.Data,
		},
		Translation: ExerciseTranslation{
			Uuid:         l.Uuid_2,
			ExerciseUuid: l.ExerciseUuid,
			Language:     l.Language,
			Name:         l.Name,
			Description:  l.Description,
		},
	}
}

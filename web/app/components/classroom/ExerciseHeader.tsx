import type { ModelsExerciseWithTranslation } from "~/api/generated/model";

export interface ExerciseHeaderProps {
  exercise: ModelsExerciseWithTranslation;
}

export default function ExerciseHeader({ exercise }: ExerciseHeaderProps) {
  return (
    <div className="leading-none">
      <h2 className="text-2xl font-bold">
        {exercise.translation?.name || "Exercise"}
      </h2>
      {exercise.translation?.description && (
        <p className="text-muted-foreground">
          {exercise.translation.description}
        </p>
      )}
    </div>
  );
}

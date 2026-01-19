import type { DbExerciseWithTranslation } from "~/api/generated/model";

export interface ExerciseHeaderProps {
  exercise: DbExerciseWithTranslation;
}

export default function ExerciseHeader({ exercise }: ExerciseHeaderProps) {
  return (
    <div className="border-b pb-4">
      <h2 className="text-2xl font-bold">
        {exercise.translation?.name || "Exercise"}
      </h2>
      {exercise.translation?.description && (
        <p className="text-muted-foreground mt-2">
          {exercise.translation.description}
        </p>
      )}
    </div>
  );
}

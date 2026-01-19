import { Link } from "react-router";
import { CheckCircle2, Circle } from "lucide-react";
import { cn } from "~/lib/utils";
import type { DbUserExerciseStatus } from "~/api/generated/model";

export interface ExerciseListItemProps {
  exercise: DbUserExerciseStatus;
  courseUuid: string;
  lessonUuid: string;
  exerciseIndex: number;
  isSelected: boolean;
}

export default function ExerciseListItem({
  exercise,
  courseUuid,
  lessonUuid,
  exerciseIndex,
  isSelected,
}: ExerciseListItemProps) {
  const isExerciseCompleted = exercise.is_completed;

  return (
    <Link
      to={`/classroom/${courseUuid}/${lessonUuid}/${exercise.exercise_uuid}`}
      className={cn(
        "flex items-center gap-2 py-2 px-3 rounded-md hover:bg-accent transition-colors",
        isSelected && "bg-accent font-medium"
      )}
    >
      {isExerciseCompleted ? (
        <CheckCircle2 className="h-4 w-4 text-green-600 shrink-0" />
      ) : (
        <Circle className="h-4 w-4 text-muted-foreground shrink-0" />
      )}
      <span className="text-sm">{`Exercise ${exerciseIndex + 1}`}</span>
    </Link>
  );
}

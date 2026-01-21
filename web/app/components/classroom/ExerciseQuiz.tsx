import type { ExercisesExerciseWithTranslation } from "~/api/generated/model";

export interface ExerciseQuizProps {
  exercise: ExercisesExerciseWithTranslation;
  language?: string;
  initialCode?: string;
  onChange?: (value: string | undefined) => void;
}

export default function ExerciseQuiz({
  exercise,
  language = "javascript",
  initialCode,
  onChange,
}: ExerciseQuizProps) {
  const getCodeValue = () => {
    if (initialCode !== undefined) {
      return initialCode;
    }
    if (exercise.data && typeof exercise.data === "object" && "code" in exercise.data) {
      return String(exercise.data.code);
    }
    return "";
  };

  return (
    <div className="flex-1 border rounded-lg overflow-hidden">
      Quiz
    </div>
  );
}

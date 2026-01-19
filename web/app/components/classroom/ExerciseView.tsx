import ExerciseHeader from "./ExerciseHeader";
import ExerciseEditor from "./ExerciseEditor";
import ExerciseActions from "./ExerciseActions";
import type { DbExerciseWithTranslation } from "~/api/generated/model";

export interface ExerciseViewProps {
  exercise: DbExerciseWithTranslation;
  language?: string;
  initialCode?: string;
  onCodeChange?: (value: string | undefined) => void;
  onSubmit?: () => void;
}

export default function ExerciseView({
  exercise,
  language,
  initialCode,
  onCodeChange,
  onSubmit,
}: ExerciseViewProps) {
  return (
    <div className="flex flex-col h-full gap-4">
      <ExerciseHeader exercise={exercise} />
      <ExerciseEditor
        exercise={exercise}
        language={language}
        initialCode={initialCode}
        onChange={onCodeChange}
      />
      <ExerciseActions onSubmit={onSubmit} />
    </div>
  );
}

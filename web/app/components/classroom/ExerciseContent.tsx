import { Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useGetMeExercisesExerciseUuid } from "~/api/generated/me/me";
import type { ModelsExerciseWithTranslation } from "~/api/generated/model";
import ExerciseCode from "./ExerciseCode";
import ExerciseQuiz from "./ExerciseQuiz";

export interface ExerciseContentProps {
  exercise?: ModelsExerciseWithTranslation;
  exerciseUuid?: string;
  language: string;
  onExerciseComplete: (exerciseUuid: string, nextLessonUuid?: string, nextExerciseUuid?: string) => void;
}

export default function ExerciseContent({
  exercise,
  exerciseUuid,
  language,
  onExerciseComplete,
}: ExerciseContentProps) {
  const { t } = useTranslation();

  const {
    data: userExerciseData,
    isLoading: isLoadingUserExercise,
  } = useGetMeExercisesExerciseUuid(exerciseUuid || "", {
    query: {
      enabled: !!exerciseUuid,
    },
  });

  if (isLoadingUserExercise) {
    return (
      <div className="flex flex-col h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        <p className="mt-4 text-muted-foreground">{t("common.loading")}</p>
      </div>
    );
  }

  if (!exercise || !userExerciseData) {
    if (exerciseUuid) {
      return (
        <div className="flex flex-col h-full items-center justify-center">
          <p className="text-muted-foreground">{t("common.error")}</p>
        </div>
      );
    }
    return (
      <div className="flex flex-col h-full items-center justify-center">
        <p className="text-muted-foreground">
          {t("common.selectExercise") || "Select an exercise to begin"}
        </p>
      </div>
    );
  }

  if (exercise.type === "code") {
    return (
      <ExerciseCode
        key={exerciseUuid}
        exercise={exercise}
        language={language}
        userExercise={userExerciseData}
        onExerciseComplete={onExerciseComplete}
      />
    );
  }

  if (exercise.type === "quiz") {
    return (
      <ExerciseQuiz
        key={exerciseUuid}
        exercise={exercise}
        language={language}
        userExercise={userExerciseData}
        onExerciseComplete={onExerciseComplete}
      />
    );
  }
}

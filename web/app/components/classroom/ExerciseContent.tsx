import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import { useGetMeExercisesExerciseUuid } from "~/api/generated/me/me";
import ExerciseView from "./ExerciseView";
import type { DbExerciseWithTranslation } from "~/api/generated/model";

export interface ExerciseContentProps {
  exercise?: DbExerciseWithTranslation;
  exerciseUuid?: string;
  language?: string;
  onCodeChange?: (value: string | undefined) => void;
  onSubmit?: () => void;
}

export default function ExerciseContent({
  exercise,
  exerciseUuid,
  language = "en",
  onCodeChange,
  onSubmit,
}: ExerciseContentProps) {
  const { t } = useTranslation();

  // Fetch user exercise data
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

  if (!exercise) {
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

  // Use user's submission code if available, otherwise use exercise default
  let codeValue = "";

  if (userExerciseData?.submission && typeof userExerciseData.submission === "object" && "code" in userExerciseData.submission) {
    codeValue = String(userExerciseData.submission.code);
  } else if (
    exercise.data &&
    typeof exercise.data === "object" &&
    "code" in exercise.data
  ) {
    codeValue = String(exercise.data.code);
  }

  return (
    <ExerciseView
      exercise={exercise}
      language={language}
      initialCode={codeValue}
      onCodeChange={onCodeChange}
      onSubmit={onSubmit}
    />
  );
}

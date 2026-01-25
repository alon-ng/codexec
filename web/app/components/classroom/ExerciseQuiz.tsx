import { CheckCircle, XCircle } from "lucide-react";
import { motion } from "motion/react";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { usePutMeExercisesExerciseUuid } from "~/api/generated/me/me";
import type { ExercisesExerciseTranslationQuizDataQuestion, ExercisesExerciseWithTranslation, MeSaveUserExerciseSubmissionRequestSubmission, MeUserExercise } from "~/api/generated/model";
import type { ExecuteResponse, UserExerciseQuizData } from '~/api/types';
import errorSound from "~/assets/error.mp3";
import { Button } from "~/components/base/Button";
import ExerciseHeader from "~/components/classroom/ExerciseHeader";
import { Label } from "~/components/ui/label";
import { RadioGroup, RadioGroupItem } from "~/components/ui/radio-group";
import { useWebSocket } from "~/hooks/useWebSocket";
import { cn } from "~/lib/utils";
import { blurInVariants } from "~/utils/animations";

export interface ExerciseQuizProps {
  exercise: ExercisesExerciseWithTranslation;
  language: string;
  userExercise: MeUserExercise;
  onExerciseComplete: (exerciseUuid: string, nextLessonUuid?: string, nextExerciseUuid?: string) => void;
}

export default function ExerciseQuiz({
  exercise,
  language,
  userExercise,
  onExerciseComplete,
}: ExerciseQuizProps) {
  const { t } = useTranslation();
  const { mutate: saveMutation } = usePutMeExercisesExerciseUuid();

  const userSubmission = userExercise.submission as unknown as UserExerciseQuizData;
  const initialAnswers = userSubmission?.answers || {};
  const initialResults = userSubmission?.results || {};
  const [answers, setAnswers] = useState<Record<string, string>>(initialAnswers);
  const [questionResults, setQuestionResults] = useState<Record<string, boolean>>(initialResults);
  const previousAnswersRef = useRef<Record<string, string>>(initialAnswers);
  const answersRef = useRef<Record<string, string>>(initialAnswers);
  const questionResultsRef = useRef<Record<string, boolean>>(initialResults);
  const [isRunning, setIsRunning] = useState(false);
  const isCompleted = Boolean(userExercise.completed_at);
  const quizData = exercise.translation.quiz_data || {};

  function onSubmissionResponse(result: ExecuteResponse) {
    setIsRunning(false);

    // Map checker results to question results
    const results: Record<string, boolean> = {};
    result.checker_results?.forEach((checker) => {
      results[checker.type] = checker.success;
    });
    setQuestionResults(results);
    questionResultsRef.current = results;

    // Save the submission with results
    saveQuiz(answersRef.current, results, exercise.uuid);

    if (result.passed) {
      onExerciseComplete(exercise.uuid, result.next_lesson_uuid, result.next_exercise_uuid);
    } else {
      const audio = new Audio(errorSound);
      audio.play();
    }
  }

  const { submit } = useWebSocket(onSubmissionResponse);

  // Helper function to create submission data
  const createSubmissionData = (currentAnswers: Record<string, string>, currentResults: Record<string, boolean>): UserExerciseQuizData => ({
    answers: currentAnswers,
    results: Object.keys(currentResults).length > 0 ? currentResults : undefined,
  });

  // Helper function to save the quiz submission
  const saveQuiz = (currentAnswers: Record<string, string>, currentResults: Record<string, boolean>, currentExerciseUuid: string) => {
    if (Object.keys(currentAnswers).length === 0) {
      return;
    }

    saveMutation({
      exerciseUuid: currentExerciseUuid,
      data: {
        type: "quiz",
        submission: createSubmissionData(currentAnswers, currentResults) as unknown as MeSaveUserExerciseSubmissionRequestSubmission,
      },
    });

    previousAnswersRef.current = { ...currentAnswers };
  };

  // Helper function to check if answers have changed
  const hasAnswersChanged = (current: Record<string, string>, previous: Record<string, string>): boolean => {
    const currentKeys = Object.keys(current).sort().join(',');
    const previousKeys = Object.keys(previous).sort().join(',');
    return currentKeys !== previousKeys ||
      Object.keys(current).some(key => current[key] !== previous[key]);
  };

  // Auto-save functionality: check every 5 seconds if answers changed
  useEffect(() => {
    const checkAndSave = () => {
      const currentAnswers = answersRef.current;
      if (hasAnswersChanged(currentAnswers, previousAnswersRef.current)) {
        saveQuiz(currentAnswers, questionResultsRef.current, exercise.uuid);
      }
    };

    const interval = setInterval(checkAndSave, 5000);

    return () => {
      clearInterval(interval);
      checkAndSave(); // Save on unmount if there are unsaved changes
    };
  }, [exercise.uuid]);

  function updateAnswerKey(questionKey: string, answerKey: string) {
    const newAnswers = { ...answers, [questionKey]: answerKey };
    setAnswers(newAnswers);
    answersRef.current = newAnswers;
  }

  const handleSubmitQuiz = () => {
    setIsRunning(true);
    submit(exercise.uuid, createSubmissionData(answers, questionResults));
  };

  return (
    <div className="flex flex-col gap-8 h-full">
      <motion.div variants={blurInVariants(0.2)} initial="hidden" animate="visible">
        <ExerciseHeader exercise={exercise} />
      </motion.div>
      <div className="flex flex-col gap-12">
        {Object.entries(quizData).map(([key, question], index) => {
          const isCorrect = questionResults[key] === true;
          const hasResult = questionResults[key] !== undefined;

          return (
            <motion.div key={key} variants={blurInVariants(0.2 + index * 0.1)} initial="hidden" animate="visible">
              <ExerciseQuizQuestion
                question={question}
                index={index}
                answerKey={answers[key]}
                hasResult={hasResult}
                isCorrect={isCorrect}
                isDisabled={isCompleted}
                updateAnswerKey={(answerKey) => updateAnswerKey(key, answerKey)}
              />
            </motion.div>
          );
        })}
      </div>
      <motion.div className="self-end me-4" variants={blurInVariants(0.3 + (Object.entries(quizData).length - 1) * 0.1)} initial="hidden" animate="visible">
        <Button className="w-36" onClick={handleSubmitQuiz} isLoading={isRunning} disabled={isCompleted}>
          {t("common.submit")}
        </Button>
      </motion.div>
    </div>
  );
}

export interface ExerciseQuizQuestionProps {
  question: ExercisesExerciseTranslationQuizDataQuestion;
  index: number;
  answerKey?: string;
  updateAnswerKey: (answerKey: string) => void;
  hasResult: boolean;
  isCorrect?: boolean;
  isDisabled?: boolean;
}

export function ExerciseQuizQuestion({ question, index, answerKey, hasResult, isCorrect, isDisabled, updateAnswerKey }: ExerciseQuizQuestionProps) {
  return (
    <motion.div className="flex flex-col gap-2" variants={blurInVariants(0.2 + index * 0.1)} initial="hidden" animate="visible">
      <div className="flex items-center gap-4">
        <div className="flex items-center justify-center text-3xl font-bold leading-none ms-4">
          {index + 1}
        </div>
        <div className="font-medium flex-1">{question.question}</div>
        {hasResult && (
          <div className={cn("flex items-center gap-1", isCorrect ? "text-green-500" : "text-red-500")}>
            {isCorrect ? <CheckCircle className="size-5" /> : <XCircle className="size-5" />}
          </div>
        )}
      </div>
      <RadioGroup className="grid grid-cols-2 gap-2" disabled={isDisabled} value={answerKey} onValueChange={updateAnswerKey}>
        {Object.entries(question.answers).map(([key, answer]) => (
          <div key={key} className={cn(
            "relative border rounded-md p-4 flex items-center gap-2 h-full has-data-[state=checked]:border-primary/50",
            hasResult && isCorrect && "has-data-[state=checked]:border-green-500",
            hasResult && !isCorrect && "has-data-[state=checked]:border-red-500",
          )}>
            <RadioGroupItem
              className="disabled:opacity-100"
              circleClassName={hasResult ? (isCorrect ? "fill-green-500 stroke-green-500" : "fill-red-500 stroke-red-500") : undefined}
              id={`${index}-answer-${key}`}
              value={key}
            />
            <Label htmlFor={`${index}-answer-${key}`} className="font-medium after:absolute after:inset-0 after:cursor-pointer">{answer}</Label>
          </div>
        ))}
      </RadioGroup>
    </motion.div>
  );
}

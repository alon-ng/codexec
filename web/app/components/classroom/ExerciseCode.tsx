import { EditorContent, useEditor } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import CodeMirror from '@uiw/react-codemirror';
import { Play } from "lucide-react";
import { motion } from "motion/react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { usePutMeExercisesExerciseUuid } from "~/api/generated/me/me";
import type { ExercisesExerciseCodeData, ExercisesExerciseWithTranslation, MeSaveUserExerciseSubmissionRequestSubmission, MeUserExercise } from "~/api/generated/model";
import type { ExecuteResponse } from '~/api/types';
import errorSound from "~/assets/error.mp3";
import { Button } from "~/components/base/Button";
import { useWebSocket } from "~/hooks/useWebSocket";
import { useLanguage } from '~/lib/useLanguage';
import { blurInVariants } from "~/utils/animations";
import LANGUAGE_MAP from "~/utils/codeLang";
import { getCodeMirrorExtensions } from "~/utils/codeMirror";
import ExerciseCodeResults from './ExerciseCodeResults';
import ExerciseHeader from "./ExerciseHeader";

export interface ExerciseCodeProps {
  exercise: ExercisesExerciseWithTranslation;
  language: string;
  userExercise: MeUserExercise;
  onExerciseComplete: (exerciseUuid: string, nextLessonUuid?: string, nextExerciseUuid?: string) => void;
}

function getCodeValue(submission: ExercisesExerciseCodeData | undefined): string {
  if (!submission?.content) {
    return "";
  }

  return submission.content;
}

function getSubmissionFromCode(code: string, language: string): ExercisesExerciseCodeData {
  const ext = LANGUAGE_MAP[language] || "txt";
  return {
    name: `main.${ext}`,
    content: code,
  };
}

export default function ExerciseCode({
  exercise,
  language,
  userExercise,
  onExerciseComplete,
}: ExerciseCodeProps) {
  const { t } = useTranslation();
  const { dir } = useLanguage();
  const { mutate: saveMutation } = usePutMeExercisesExerciseUuid();

  const initialCode = useMemo(() => {
    const userSubmission = userExercise.submission as unknown as ExercisesExerciseCodeData;
    const hasUserSubmission = Boolean(userSubmission?.name && userSubmission?.content);
    return getCodeValue(hasUserSubmission ? userSubmission : exercise.code_data);
  }, [userExercise.submission, exercise.code_data]);

  const previousCodeRef = useRef<string>(initialCode);
  const codeValueRef = useRef<string>(initialCode);
  const [codeValue, setCodeValue] = useState(initialCode);

  const [isRunning, setIsRunning] = useState(false);
  const [resultTab, setResultTab] = useState<string>("console");

  useEffect(() => {
    codeValueRef.current = initialCode;
    setCodeValue(initialCode);
  }, [initialCode]);

  function onSubmissionResponse(result: ExecuteResponse) {
    setIsRunning(false);
    if (result.passed) {
      onExerciseComplete(exercise.uuid, result.next_lesson_uuid, result.next_exercise_uuid);
    } else {
      const audio = new Audio(errorSound);
      audio.play();
    }

    if (result.stderr) {
      setResultTab("errors");
    }
  }
  const { submit, lastResult } = useWebSocket(onSubmissionResponse);

  const readOnlyLines: number[] = [];

  const extensions = useMemo(() => {
    return getCodeMirrorExtensions(language, readOnlyLines);
  }, [language, readOnlyLines]);

  const editor = useEditor({
    extensions: [StarterKit],
    content: exercise.translation?.code_data?.instructions || "No instructions available.",
    editable: false,
    immediatelyRender: false,
  })

  // Helper function to save the code
  const saveCode = (currentCode: string, currentExerciseUuid: string) => {
    if (currentCode.trim() === "") {
      return;
    }

    const s = getSubmissionFromCode(currentCode, language);

    saveMutation({
      exerciseUuid: currentExerciseUuid,
      data: {
        type: "code",
        submission: s as unknown as MeSaveUserExerciseSubmissionRequestSubmission,
      },
    });

    // Update the ref to the current value
    previousCodeRef.current = currentCode;
  };

  // Auto-save functionality: check every 5 seconds if code changed
  useEffect(() => {
    const interval = setInterval(() => {
      const currentCode = codeValueRef.current;
      if (currentCode !== previousCodeRef.current) {
        saveCode(currentCode, exercise.uuid);
      }
    }, 5000); // Check every 5 seconds

    return () => {
      clearInterval(interval);
      // Before unmount, check if there are unsaved changes and save them
      const currentCode = codeValueRef.current;
      if (currentCode !== previousCodeRef.current) {
        saveCode(currentCode, exercise.uuid);
      }
    };
  }, [exercise.uuid]);

  const handleCodeChange = (value: string) => {
    setCodeValue(value);
    codeValueRef.current = value;
  };

  const handleRunCode = () => {
    const s = getSubmissionFromCode(codeValue, language);
    setIsRunning(true);
    submit(exercise.uuid, s);
  };

  return (
    <div className="flex justify-start h-full">
      <div className="flex flex-col w-1/2 gap-4">
        <motion.div variants={blurInVariants(0.2)} initial="hidden" animate="visible">
          <ExerciseHeader exercise={exercise} />
        </motion.div>
        <motion.div variants={blurInVariants(0.3)} initial="hidden" animate="visible">
          <EditorContent className="prose prose-sm dark:prose-invert" editor={editor} />
        </motion.div>
      </div>
      <div className="w-1/2 h-full flex flex-col gap-2">
        <motion.div className="flex justify-end" variants={blurInVariants(0.5)} initial="hidden" animate="visible">
          <Button variant="outline" onClick={handleRunCode} isLoading={isRunning} disabled={Boolean(userExercise.completed_at)}>
            {t("common.run")}
            <Play className="size-4" />
          </Button>
        </motion.div>
        <motion.div className="flex flex-col flex-1 border rounded-lg overflow-hidden" variants={blurInVariants(0.4)} initial="hidden" animate="visible">
          <CodeMirror
            dir="ltr"
            className="flex-1"
            height="100%"
            value={codeValue}
            onChange={handleCodeChange}
            extensions={extensions}
            theme="light"
          />
          <ExerciseCodeResults resultTab={resultTab} setResultTab={setResultTab} lastResult={lastResult} />
        </motion.div>
      </div>
    </div>
  );
}

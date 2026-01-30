import { EditorContent, useEditor } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import CodeMirror from '@uiw/react-codemirror';
import { Play } from "lucide-react";
import { motion } from "motion/react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { usePutMeExercisesExerciseUuid } from "~/api/generated/me/me";
import type { MeSaveUserExerciseSubmissionRequestSubmission, ModelsExerciseCodeData, ModelsExerciseWithTranslation, ModelsUserExercise } from "~/api/generated/model";
import type { ExecuteResponse } from '~/api/types';
import codyAvatar from "~/assets/cody-256.png";
import errorSound from "~/assets/error.mp3";
import { Button } from "~/components/base/Button";
import Chat from '~/components/chat/Chat';
import ExerciseCodeResults from '~/components/classroom/ExerciseCodeResults';
import ExerciseHeader from '~/components/classroom/ExerciseHeader';
import { useWebSocket } from "~/hooks/useWebSocket";
import { blurInVariants } from "~/utils/animations";
import LANGUAGE_MAP from "~/utils/codeLang";
import { getCodeMirrorExtensions } from "~/utils/codeMirror";
import { prose } from '~/utils/prose';

export interface ExerciseCodeProps {
  exercise: ModelsExerciseWithTranslation;
  language: string;
  userExercise: ModelsUserExercise;
  onExerciseComplete: (exerciseUuid: string, nextLessonUuid?: string, nextExerciseUuid?: string) => void;
}

function getCodeValue(submission: ModelsExerciseCodeData | undefined): string {
  if (!submission?.content) {
    return "";
  }

  return submission.content;
}

function getSubmissionFromCode(code: string, language: string): ModelsExerciseCodeData {
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
  const { mutate: saveMutation } = usePutMeExercisesExerciseUuid();

  const initialCode = useMemo(() => {
    const userSubmission = userExercise.submission as unknown as ModelsExerciseCodeData;
    const hasUserSubmission = Boolean(userSubmission?.name && userSubmission?.content);
    return getCodeValue(hasUserSubmission ? userSubmission : exercise.code_data);
  }, [userExercise.submission, exercise.code_data]);

  const previousCodeRef = useRef<string>(initialCode);
  const codeValueRef = useRef<string>(initialCode);
  const [codeValue, setCodeValue] = useState(initialCode);

  const [isRunning, setIsRunning] = useState(false);
  const [resultTab, setResultTab] = useState<string>("console");
  const [isChatOpen, setIsChatOpen] = useState(false);

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
    <div className="flex justify-start h-full gap-2">
      <div className="flex-1 flex flex-col gap-4">
        <motion.div variants={blurInVariants(0.2)} initial="hidden" animate="visible">
          <ExerciseHeader exercise={exercise} />
        </motion.div>
        <motion.div variants={blurInVariants(0.3)} initial="hidden" animate="visible">
          <EditorContent className={prose} editor={editor} />
        </motion.div>
      </div>
      <div className="flex-1 h-full flex flex-col gap-2">
        <motion.div className="flex justify-end gap-2" variants={blurInVariants(0.5)} initial="hidden" animate="visible">
          <Button variant="outline" onClick={handleRunCode} isLoading={isRunning} disabled={Boolean(userExercise.completed_at)}>
            {t("common.run")}
            <Play className="size-4" />
          </Button>
        </motion.div>
        <motion.div className="flex flex-col flex-1 border rounded-lg overflow-hidden relative" variants={blurInVariants(0.4)} initial="hidden" animate="visible">
          <CodeMirror
            dir="ltr"
            className="flex-1"
            height="100%"
            value={codeValue}
            onChange={handleCodeChange}
            extensions={extensions}
            theme="light"
            readOnly={Boolean(userExercise.completed_at)}
          />
          <ExerciseCodeResults resultTab={resultTab} setResultTab={setResultTab} lastResult={lastResult} />
          <img src={codyAvatar} className="size-16 absolute bottom-2 right-2 cursor-pointer hover:translate-y-[-0.25rem] transition-all duration-200" onClick={() => setIsChatOpen(!isChatOpen)} />
        </motion.div>
      </div>
      {isChatOpen && <Chat exerciseInstructions={exercise.translation?.code_data?.instructions || ""} exerciseCode={codeValue} exerciseUuid={exercise.uuid} />}
    </div>
  );
}

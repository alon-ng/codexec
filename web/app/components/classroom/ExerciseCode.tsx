import { useState, useMemo, useEffect, useRef } from "react";
import CodeMirror from '@uiw/react-codemirror';
import type { ExercisesExerciseCodeData, ExercisesExerciseWithTranslation, MeUserExercise } from "~/api/generated/model";
import ExerciseHeader from "./ExerciseHeader";
import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import { loadLanguage, type LanguageName } from '@uiw/codemirror-extensions-langs';
import { usePutMeExercisesExerciseUuid } from "~/api/generated/me/me";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";
import { Button } from "../base/Button";
import { Play } from "lucide-react";
import { useTranslation } from "react-i18next";

export interface ExerciseEditorProps {
  exercise: ExercisesExerciseWithTranslation;
  language: string;
  userExercise: MeUserExercise;
  onChange?: (value: string | undefined) => void;
}

const LANGUAGE_MAP: Record<string, LanguageName> = {
  "javascript": "js",
  "python": "py",
};

function getCodeValue(submission: ExercisesExerciseCodeData | undefined): string {
  if (!submission) {
    return "";
  }

  const mainFile = submission.files[0];
  return mainFile.content;
}

export default function ExerciseCode({
  exercise,
  language,
  userExercise,
  onChange,
}: ExerciseEditorProps) {
  const { t } = useTranslation();
  
  // Determine the submission structure and initial code value
  const { submission, initialCode } = useMemo(() => {
    const userSubmission = userExercise.submission as unknown as ExercisesExerciseCodeData;
    const hasUserSubmission = userSubmission && 
      typeof userSubmission === 'object' && 
      'name' in userSubmission && 
      userSubmission.name;
    
    if (hasUserSubmission) {
      return {
        submission: userSubmission,
        initialCode: getCodeValue(userSubmission),
      };
    } else {
      const exerciseCodeData = exercise.code_data;
      return {
        submission: exerciseCodeData || {
          name: "root",
          directories: [],
          files: [{ name: "main.py", ext: "py", content: "" }],
        },
        initialCode: getCodeValue(exerciseCodeData),
      };
    }
  }, [userExercise.submission, exercise.code_data]);

  const [codeValue, setCodeValue] = useState(initialCode);
  const previousCodeRef = useRef<string>(initialCode);
  const codeValueRef = useRef<string>(initialCode);
  const saveMutation = usePutMeExercisesExerciseUuid();

  // Update refs when initial code changes (e.g., when switching exercises)
  useEffect(() => {
    previousCodeRef.current = initialCode;
    codeValueRef.current = initialCode;
    setCodeValue(initialCode);
  }, [initialCode]);

  const extensions = useMemo(() => {
    const langName = LANGUAGE_MAP[language];
    const langExt = loadLanguage(langName);
    return langExt ? [langExt] : [];
  }, [language]);

  const editor = useEditor({
    extensions: [StarterKit],
    content: exercise.translation?.code_data?.instructions || "No instructions available.",
    editable: false,
    immediatelyRender: false, 
  })

  // Helper function to save the code
  const saveCode = (currentCode: string, currentSubmission: ExercisesExerciseCodeData, currentLanguage: string, currentExerciseUuid: string) => {
    if (currentCode.trim() === "") {
      return;
    }

    const submissionData: ExercisesExerciseCodeData = {
      name: currentSubmission.name || "root",
      directories: currentSubmission.directories || [],
      files: currentSubmission.files.length > 0
        ? currentSubmission.files.map((file, index) => 
            index === 0 
              ? { ...file, content: currentCode }
              : file
          )
        : [{ 
            name: currentLanguage === "python" ? "main.py" : "main.js", 
            ext: currentLanguage === "python" ? "py" : "js", 
            content: currentCode 
          }],
    };

    saveMutation.mutate({
      exerciseUuid: currentExerciseUuid,
      data: {
        type: "code",
        submission: submissionData as unknown as Record<string, unknown>,
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
        saveCode(currentCode, submission, language, exercise.uuid);
      }
    }, 5000); // Check every 5 seconds

    return () => {
      clearInterval(interval);
      // Before unmount, check if there are unsaved changes and save them
      const currentCode = codeValueRef.current;
      if (currentCode !== previousCodeRef.current) {
        saveCode(currentCode, submission, language, exercise.uuid);
      }
    };
  }, [exercise.uuid, submission, saveMutation, language]);

  const handleCodeChange = (value: string) => {
    setCodeValue(value);
    codeValueRef.current = value; // Keep ref in sync
    onChange?.(value);
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
          <Button variant="outline">
            {t("common.run")}
            <Play className="size-4" />
          </Button>
        </motion.div>
        <motion.div className="flex-1 border rounded-lg overflow-hidden" variants={blurInVariants(0.4)} initial="hidden" animate="visible">
          <CodeMirror
            dir="ltr"
            className="h-full"
            height="100%"
            value={codeValue} 
            onChange={handleCodeChange} 
            extensions={extensions} 
            theme="light" 
          />
        </motion.div>
      </div>
    </div>
  );
}

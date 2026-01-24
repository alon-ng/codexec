import { EditorContent, useEditor } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import CodeMirror from '@uiw/react-codemirror';
import { Ban, CheckCircle, ChevronRightSquare, FlaskConical, Play, XCircle } from "lucide-react";
import { motion } from "motion/react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { usePutMeExercisesExerciseUuid } from "~/api/generated/me/me";
import type { ExercisesExerciseCodeData, ExercisesExerciseWithTranslation, MeSaveUserExerciseSubmissionRequestSubmission, MeUserExercise } from "~/api/generated/model";
import { Button } from "~/components/base/Button";
import { useWebSocket } from "~/hooks/useWebSocket";
import { useLanguage } from '~/lib/useLanguage';
import { cn } from '~/lib/utils';
import { blurInVariants } from "~/utils/animations";
import LANGUAGE_MAP from "~/utils/codeLang";
import { getCodeMirrorExtensions } from "~/utils/codeMirror";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import ExerciseHeader from "./ExerciseHeader";

export interface ExerciseEditorProps {
  exercise: ExercisesExerciseWithTranslation;
  language: string;
  userExercise: MeUserExercise;
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
}: ExerciseEditorProps) {
  const { t } = useTranslation();
  const { submitCode, lastResult } = useWebSocket();
  const { dir } = useLanguage();
  const saveMutation = usePutMeExercisesExerciseUuid();

  const { submission, initialCode } = useMemo(() => {
    const userSubmission = userExercise.submission as unknown as ExercisesExerciseCodeData;
    const hasUserSubmission = Boolean(userSubmission?.name && userSubmission?.content);

    if (hasUserSubmission) {
      return {
        submission: userSubmission,
        initialCode: getCodeValue(userSubmission),
      };
    }

    const code = getCodeValue(exercise.code_data);
    return {
      submission: getSubmissionFromCode(code, language),
      initialCode: code,
    };
  }, [userExercise.submission, exercise.code_data, language]);

  const [codeValue, setCodeValue] = useState(initialCode);
  const previousCodeRef = useRef<string>(initialCode);
  const codeValueRef = useRef<string>(initialCode);
  const [resultTab, setResultTab] = useState<string>("console");
  const [isRunning, setIsRunning] = useState(false);

  useEffect(() => {
    previousCodeRef.current = initialCode;
    codeValueRef.current = initialCode;
    setCodeValue(initialCode);
  }, [initialCode]);

  useEffect(() => {
    setIsRunning(false);
    if (lastResult && lastResult.stderr) {
      setResultTab("errors");
    }
  }, [lastResult]);

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
  const saveCode = (currentCode: string, currentSubmission: ExercisesExerciseCodeData, currentLanguage: string, currentExerciseUuid: string) => {
    if (currentCode.trim() === "") {
      return;
    }

    const s = getSubmissionFromCode(codeValue, language);

    saveMutation.mutate({
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
  }, [exercise.uuid]);

  const handleCodeChange = (value: string) => {
    setCodeValue(value);
    codeValueRef.current = value;
  };

  const handleRunCode = () => {
    const s = getSubmissionFromCode(codeValue, language);
    setIsRunning(true);
    submitCode(exercise.uuid, s);
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
          <Button variant="outline" onClick={handleRunCode} isLoading={isRunning}>
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
          {lastResult && (
            <motion.div
              variants={blurInVariants()}
              initial="hidden"
              animate="visible"
              className="flex h-48 text-xs overflow-auto"
            >
              <Tabs className="w-full" value={resultTab} onValueChange={setResultTab}>
                <TabsList className="h-8 w-full justify-start rounded-none p-1 border-t shadow-sm" dir={dir}>
                  <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="console">
                    <ChevronRightSquare className="size-3" />
                    {t("common.console")}
                  </TabsTrigger>
                  <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="errors">
                    <Ban className="size-3" />
                    {t("common.errors")}
                  </TabsTrigger>
                  <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="tests">
                    <FlaskConical className="size-3" />
                    {t("common.tests")}
                  </TabsTrigger>
                </TabsList>
                <TabsContent className="text-xs px-3 font-mono" value="console">
                  <div className="whitespace-pre-wrap">
                    {lastResult.stdout || <span className="text-muted-foreground">{t("common.noOutput") || "No output"}</span>}
                  </div>
                </TabsContent>
                <TabsContent className="text-xs px-3 font-mono" value="errors">
                  <div className="text-red-400 whitespace-pre-wrap">
                    {lastResult.stderr || <span className="text-muted-foreground">{t("common.noErrors") || "No errors"}</span>}
                  </div>
                </TabsContent>
                <TabsContent className="text-xs font-mono" value="tests">
                  {lastResult.checker_results?.length === 0 ? <span className="text-muted-foreground">{t("common.noTests") || "No tests"}</span> : (
                    <>
                      {lastResult.checker_results?.map((result) => (
                        <div className={cn("flex items-center gap-1.5 py-1 px-3", result.success ? "text-green-400 bg-green-50" : "text-red-400 bg-red-50")} key={result.type}>
                          {result.success ? <CheckCircle className="size-3" /> : <XCircle className="size-3" />}
                          <span>{result.message}</span>
                        </div>
                      ))}
                    </>
                  )}
                </TabsContent>
              </Tabs>
            </motion.div>
          )}
        </motion.div>
      </div>
    </div>
  );
}

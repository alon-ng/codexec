import Editor from "@monaco-editor/react";
import type { DbExerciseWithTranslation } from "~/api/generated/model";

export interface ExerciseEditorProps {
  exercise: DbExerciseWithTranslation;
  language?: string;
  initialCode?: string;
  onChange?: (value: string | undefined) => void;
}

export default function ExerciseEditor({
  exercise,
  language = "javascript",
  initialCode,
  onChange,
}: ExerciseEditorProps) {
  const getCodeValue = () => {
    // Use initialCode if provided (from user submission), otherwise use exercise default
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
      <Editor
        height="100%"
        defaultLanguage={language}
        theme="vs-dark"
        value={getCodeValue()}
        onChange={onChange}
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: "on",
          scrollBeyondLastLine: false,
          automaticLayout: true,
        }}
      />
    </div>
  );
}

import ListItem from "@tiptap/extension-list-item";
import { EditorContent, useEditor } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import { motion } from "motion/react";
import { useMemo } from "react";
import type { ModelsLessonFull } from "~/api/generated/model";
import { blurInVariants } from "~/utils/animations";
import { prose } from "~/utils/prose";

export interface LessonContentProps {
  lesson?: ModelsLessonFull;
}

export default function LessonContent({ lesson }: LessonContentProps) {
  if (!lesson) {
    return null;
  }

  console.log(lesson);


  const editor = useEditor({
    extensions: [StarterKit.configure({
      listItem: false,
    }), ListItem.extend({ content: "text*" })],
    content: lesson.translation?.content || "No content available.",
    editable: false,
    immediatelyRender: false,
  });

  const lessonName = useMemo(() => {
    return lesson.translation?.name || "Lesson";
  }, [lesson.translation?.name]);

  const lessonDescription = useMemo(() => {
    return lesson.translation?.description || "";
  }, [lesson.translation?.description]);

  return (
    <div className="flex flex-col h-full gap-4">
      <motion.div variants={blurInVariants(0.2)} initial="hidden" animate="visible">
        <div className="leading-none">
          <h2 className="text-2xl font-bold">
            {lessonName}
          </h2>
          {lessonDescription && (
            <p className="text-muted-foreground">
              {lessonDescription}
            </p>
          )}
        </div>
      </motion.div>
      <motion.div variants={blurInVariants(0.3)} initial="hidden" animate="visible" className="flex-1 overflow-y-auto">
        <EditorContent className={prose} editor={editor} />
      </motion.div>
    </div>
  );
}

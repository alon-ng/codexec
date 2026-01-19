import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Accordion } from "~/components/ui/accordion";
import LessonAccordionItem from "./LessonAccordionItem";
import type { DbUserCourseFull } from "~/api/generated/model";
import type { DbCourseFull } from "~/api/generated/model";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";

export interface LessonSidebarProps {
  userCourseData: DbUserCourseFull;
  courseData: DbCourseFull;
  courseUuid: string;
  selectedLessonUuid?: string;
  selectedExerciseUuid?: string;
  onAccordionChange: (value: string | undefined) => void;
}

export default function LessonSidebar({
  userCourseData,
  courseData,
  courseUuid,
  selectedLessonUuid,
  selectedExerciseUuid,
  onAccordionChange,
}: LessonSidebarProps) {
  const { t } = useTranslation();
  const [openLesson, setOpenLesson] = useState<string | undefined>(
    selectedLessonUuid || undefined
  );

  useEffect(() => {
    if (selectedLessonUuid) {
      setOpenLesson(selectedLessonUuid);
    }
  }, [selectedLessonUuid]);

  const handleValueChange = (value: string | undefined) => {
    setOpenLesson(value);
    onAccordionChange(value);
  };

  const lessons = userCourseData.lessons || [];

  return (
    <aside className="w-80 border-r pr-6 overflow-y-auto">
      <motion.h2 variants={blurInVariants()} initial="hidden" animate="visible" className="text-lg font-semibold mb-4">
        {t("course.lessons") || "Lessons"}
      </motion.h2>
      <Accordion
        type="single"
        collapsible
        value={openLesson}
        onValueChange={handleValueChange}
        className="w-full"
      >
        {lessons.map((lesson, index) => {
          // Find corresponding lesson data from courseData
          const lessonData = courseData.lessons?.find(
            (l) => l.uuid === lesson.lesson_uuid
          );

          return (
            <LessonAccordionItem
              key={lesson.lesson_uuid || index}
              lesson={lesson}
              lessonIndex={index}
              courseUuid={courseUuid}
              selectedLessonUuid={selectedLessonUuid}
              selectedExerciseUuid={selectedExerciseUuid}
              lessonData={lessonData}
            />
          );
        })}
      </Accordion>
    </aside>
  );
}

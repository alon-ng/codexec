import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Accordion } from "~/components/ui/accordion";
import { Button } from "~/components/base/Button";
import { PanelLeftOpen, PanelLeftClose } from "lucide-react";
import LessonAccordionItem from "./LessonAccordionItem";
import type { MeUserCourseFull } from "~/api/generated/model";
import type { CoursesCourseFull } from "~/api/generated/model";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";
import { cn } from "~/lib/utils";

export interface LessonSidebarProps {
  userCourseData: MeUserCourseFull;
  courseData: CoursesCourseFull;
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
  const [isCollapsed, setIsCollapsed] = useState(false);
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

  const toggleCollapse = () => {
    setIsCollapsed(!isCollapsed);
  };

  const lessons = userCourseData.lessons || [];

  return (
    <aside
      className={cn("relative border-e overflow-y-auto transition-all duration-300", isCollapsed ? "w-12" : "w-80")}
    >
      <div className={cn("flex items-center mb-4", isCollapsed ? "justify-center" : "justify-between pe-6")}>
        {!isCollapsed && (
          <motion.h2
            variants={blurInVariants()}
            initial="hidden"
            animate="visible"
            className="text-lg font-semibold"
          >
            {t("course.lessons") || "Lessons"}
          </motion.h2>
        )}
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleCollapse}
          className="shrink-0 rtl:rotate-180"
          aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {isCollapsed ? (
            <PanelLeftOpen className="h-4 w-4" />
          ) : (
            <PanelLeftClose className="h-4 w-4" />
          )}
        </Button>
      </div>
      {!isCollapsed && (
        <div className="pe-6">
          <Accordion
            type="single"
            collapsible
            value={openLesson}
            onValueChange={handleValueChange}
            className="w-full"
          >
            {lessons.map((lesson, index) => {
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
        </div>
      )}
    </aside>
  );
}

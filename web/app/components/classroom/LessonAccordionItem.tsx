import {
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "~/components/ui/accordion";
import { CheckCircle2, Circle } from "lucide-react";
import { cn } from "~/lib/utils";
import ExerciseListItem from "./ExerciseListItem";
import type { DbUserLessonStatus } from "~/api/generated/model";
import type { DbLessonFull } from "~/api/generated/model";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";

export interface LessonAccordionItemProps {
  lesson: DbUserLessonStatus;
  lessonIndex: number;
  courseUuid: string;
  selectedLessonUuid?: string;
  selectedExerciseUuid?: string;
  lessonData?: DbLessonFull;
}

export default function LessonAccordionItem({
  lesson,
  lessonIndex,
  courseUuid,
  selectedLessonUuid,
  selectedExerciseUuid,
  lessonData,
}: LessonAccordionItemProps) {
  const exercises = lesson.exercises || [];
  const isLessonCompleted = lesson.is_completed;
  const isSelected = lesson.lesson_uuid === selectedLessonUuid;

  const lessonName = lessonData?.translation?.name || `Lesson ${lessonIndex + 1}`;

  return (
    <motion.div variants={blurInVariants(lessonIndex * 0.1)} initial="hidden" animate="visible">
      <AccordionItem
        key={lesson.lesson_uuid || lessonIndex}
        value={lesson.lesson_uuid || `lesson-${lessonIndex}`}
        className="border-b"
      >
        <AccordionTrigger
          className={cn("hover:no-underline", isSelected && "font-semibold")}
        >
          <div className="flex items-center gap-2 flex-1">
            {isLessonCompleted ? (
              <CheckCircle2 className="h-4 w-4 text-green-600" />
            ) : (
              <Circle className="h-4 w-4 text-muted-foreground" />
            )}
            <span>{lessonName}</span>
          </div>
        </AccordionTrigger>
        <AccordionContent>
          <div className="flex flex-col gap-1 pl-6">
            {exercises.map((exercise, exIndex) => (
              <ExerciseListItem
                key={exercise.exercise_uuid || exIndex}
                exercise={exercise}
                courseUuid={courseUuid}
                lessonUuid={lesson.lesson_uuid!}
                exerciseIndex={exIndex}
                isSelected={exercise.exercise_uuid === selectedExerciseUuid}
              />
            ))}
          </div>
        </AccordionContent>
      </AccordionItem>
    </motion.div>
  );
}

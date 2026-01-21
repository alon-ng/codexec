import {
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "~/components/ui/accordion";
import { CheckCircle2, Circle } from "lucide-react";
import { cn } from "~/lib/utils";
import ExerciseListItem from "./ExerciseListItem";
import type { MeUserLessonStatus } from "~/api/generated/model";
import type { LessonsLessonFull } from "~/api/generated/model";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";

export interface LessonAccordionItemProps {
  lesson: MeUserLessonStatus;
  lessonIndex: number;
  courseUuid: string;
  selectedLessonUuid?: string;
  selectedExerciseUuid?: string;
  lessonData?: LessonsLessonFull;
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

  function getExerciseData(exerciseUuid: string) {
    return lessonData?.exercises?.find((e) => e.uuid === exerciseUuid);
  }

  return (
    <motion.div variants={blurInVariants(lessonIndex * 0.1)} initial="hidden" animate="visible">
      <AccordionItem
        key={lesson.lesson_uuid || lessonIndex}
        value={lesson.lesson_uuid || `lesson-${lessonIndex}`}
        className="border-b"
      >
        <AccordionTrigger
          className={cn("hover:no-underline cursor-pointer", isSelected && "font-semibold")}
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
          <div className="flex flex-col gap-1 ps-6">
            {exercises.map((exercise, exIndex) => (
              <ExerciseListItem
                key={exercise.exercise_uuid || exIndex}
                exercise={exercise}
                exerciseData={getExerciseData(exercise.exercise_uuid)}
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

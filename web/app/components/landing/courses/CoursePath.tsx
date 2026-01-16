import { cn } from "~/lib/utils";
import { Check, Lock, LockOpen } from "lucide-react";
import type { DbCourseFull, DbLessonFull } from "~/api/generated/model";
import { useLanguage } from "~/lib/useLanguage";
import { motion } from "motion/react";
import { blurInVariants } from "~/utils/animations";

interface CoursePathProps {
    course: Required<DbCourseFull>;
}

const ANIMATION_DELAY = 0.1;
const LESSONS_PER_ROW = 3;
const PATH_COLOR = "bg-muted-foreground/20";
const CONNECTOR_BORDER_COLOR = "border-muted-foreground/20";

/**
 * Chunks an array into groups of specified size
 */
function chunkArray<T>(array: T[], chunkSize: number): T[][] {
    const chunks: T[][] = [];
    for (let i = 0; i < array.length; i += chunkSize) {
        chunks.push(array.slice(i, i + chunkSize));
    }
    return chunks;
}

interface RowConnectorProps {
    isEvenRow: boolean;
    isRTL: boolean;
    animationDelay: number;
}

/**
 * Renders the U-shaped connector between rows
 */
function RowConnector({ isEvenRow, isRTL, animationDelay }: RowConnectorProps) {
    const connectorClasses = cn(
        "absolute top-[17.7%] w-8 sm:w-12 h-[calc(100%+0.25rem)] border-y-4 z-20",
        CONNECTOR_BORDER_COLOR,
        isEvenRow
            ? cn(
                "end-0 border-e-4 rounded-e-[2rem] border-s-0",
                isRTL ? "-translate-x-full" : "translate-x-full"
            )
            : cn(
                "start-0 border-s-4 rounded-s-[2rem] border-e-0",
                isRTL ? "translate-x-full" : "-translate-x-full"
            )
    );

    return (
        <motion.div
            className={connectorClasses}
            variants={blurInVariants(animationDelay, 0)}
            initial="hidden"
            animate="visible"
        />
    );
}

interface LessonNodeProps {
    lesson: DbLessonFull;
    globalIndex: number;
    isCompleted: boolean;
    showLeftLine: boolean;
    showRightLine: boolean;
    animationDelay: number;
}

/**
 * Renders a single lesson node with its connection lines
 */
function LessonNode({
    lesson,
    globalIndex,
    isCompleted,
    showLeftLine,
    showRightLine,
    animationDelay,
}: LessonNodeProps) {
    const nodeClasses = cn(
        "relative w-10 h-10 rounded-full border-2 flex items-center justify-center bg-background transition-all duration-200 cursor-pointer",
        isCompleted
            ? "border-primary bg-primary text-primary-foreground"
            : "border-muted-foreground group-hover:border-primary hover:shadow-sm group-hover:scale-115"
    );

    return (
        <motion.div
            className="relative w-1/3 flex flex-col items-center z-10 group"
            variants={blurInVariants(animationDelay, 0)}
            initial="hidden"
            animate="visible"
        >
            {/* Horizontal connection lines */}
            {showLeftLine && (
                <div className={cn("absolute top-4.5 start-0 w-1/2 h-1 -z-10", PATH_COLOR)} />
            )}
            {showRightLine && (
                <div className={cn("absolute top-4.5 end-0 w-1/2 h-1 -z-10", PATH_COLOR)} />
            )}

            {/* Lesson number/check icon */}
            <div className={nodeClasses}>
                {isCompleted ?
                    <Check className="w-5 h-5" /> :
                    <span className="font-semibold text-sm">{globalIndex + 1}</span>
                }
                {lesson.is_public ?
                    <LockOpen className="absolute -bottom-3 -end-3 w-6 h-6 text-muted-foreground p-1 bg-background rounded-md" /> :
                    <Lock className="absolute -bottom-3 -end-3 w-6 h-6 text-muted-foreground p-1 bg-background rounded-md" />}
            </div>

            {/* Lesson details */}
            <div className="mt-3 text-center px-1 w-full">
                <h4 className="font-semibold text-sm leading-tight mb-1 wrap-break-word">
                    {lesson.translation?.name || `Lesson ${globalIndex + 1}`}
                </h4>
                <p className="text-xs text-muted-foreground line-clamp-3 max-w-[140px] mx-auto hidden sm:block">
                    {lesson.translation?.description}
                </p>
            </div>
        </motion.div>
    );
}

interface CourseRowProps {
    rowLessons: DbLessonFull[];
    rowIndex: number;
    totalRows: number;
    isRTL: boolean;
}

/**
 * Determines if connection lines should be shown for a lesson
 */
function shouldShowConnectionLines(
    isEvenRow: boolean,
    isFirstInRow: boolean,
    isLastInRow: boolean,
    isFirstRow: boolean,
    isLastRow: boolean
) {
    // For even rows (left-to-right): show left line if not first item OR not first row
    // For odd rows (right-to-left): show left line if not last item OR not last row
    const showLeftLine = isEvenRow
        ? !isFirstInRow || !isFirstRow
        : !isLastInRow || !isLastRow;

    // For even rows: show right line if not last item OR not last row
    // For odd rows: show right line if not first item OR not first row
    const showRightLine = isEvenRow
        ? !isLastInRow || !isLastRow
        : !isFirstInRow || !isFirstRow;

    return { showLeftLine, showRightLine };
}

/**
 * Renders a single row of lessons
 */
function CourseRow({ rowLessons, rowIndex, totalRows, isRTL }: CourseRowProps) {
    const isEvenRow = rowIndex % 2 === 0;
    const isLastRow = rowIndex === totalRows - 1;
    const isFirstRow = rowIndex === 0;

    const rowDirection = isEvenRow ? "flex-row" : "flex-row-reverse";

    // Calculate connector animation delay - appears after the last lesson in the row
    const connectorAnimationDelay = !isLastRow
        ? (rowIndex * LESSONS_PER_ROW + rowLessons.length) * ANIMATION_DELAY
        : 0;

    return (
        <div key={rowIndex} className="relative h-48">
            {/* U-shaped connector to next row */}
            {!isLastRow && (
                <RowConnector
                    isEvenRow={isEvenRow}
                    isRTL={isRTL}
                    animationDelay={connectorAnimationDelay}
                />
            )}

            {/* Lesson nodes */}
            <div className={cn("flex h-full items-start pt-4", rowDirection)}>
                {rowLessons.map((lesson, index) => {
                    const globalIndex = rowIndex * LESSONS_PER_ROW + index;
                    const isLastInRow = index === rowLessons.length - 1;
                    const isFirstInRow = index === 0;

                    // TODO: Replace with actual completion status from user progress
                    const isCompleted = false;

                    const { showLeftLine, showRightLine } = shouldShowConnectionLines(
                        isEvenRow,
                        isFirstInRow,
                        isLastInRow,
                        isFirstRow,
                        isLastRow
                    );

                    const animationDelay = globalIndex * ANIMATION_DELAY;

                    return (
                        <LessonNode
                            key={lesson.uuid || globalIndex}
                            lesson={lesson}
                            globalIndex={globalIndex}
                            isCompleted={isCompleted}
                            showLeftLine={showLeftLine}
                            showRightLine={showRightLine}
                            animationDelay={animationDelay}
                        />
                    );
                })}
            </div>
        </div>
    );
}

export default function CoursePath({ course }: CoursePathProps) {
    const { isRTL } = useLanguage();
    const lessons = course.lessons;
    const rows = chunkArray(lessons, LESSONS_PER_ROW);

    return (
        <div className="w-full max-w-3xl mx-auto p-4 md:p-8">
            <div className="relative">
                {rows.map((rowLessons, rowIndex) => (
                    <CourseRow
                        key={rowIndex}
                        rowLessons={rowLessons}
                        rowIndex={rowIndex}
                        totalRows={rows.length}
                        isRTL={isRTL}
                    />
                ))}
            </div>
        </div>
    );
}

import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { useGetCoursesUuid } from "~/api/generated/courses/courses";
import { useGetMeCoursesCourseUuid } from "~/api/generated/me/me";
import type { ModelsUserExerciseStatus } from "~/api/generated/model";
import courseCompleteImage from "~/assets/course-complete.png";
import successSound from "~/assets/success.mp3";
import PageHeader, { type BreadcrumbProps } from "~/components/PageHeader";
import { Button } from "~/components/base/Button";
import ClassroomError from "~/components/classroom/ClassroomError";
import ClassroomLoading from "~/components/classroom/ClassroomLoading";
import ExerciseContent from "~/components/classroom/ExerciseContent";
import LessonContent from "~/components/classroom/LessonContent";
import LessonSidebar from "~/components/classroom/LessonSidebar";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "~/components/ui/dialog";

export default function Classroom() {
    const { t, i18n } = useTranslation();
    const navigate = useNavigate();
    const [showCourseCompletedDialog, setShowCourseCompletedDialog] = useState(false);

    const { courseUuid, lessonUuid, exerciseUuid } = useParams<{
        courseUuid: string;
        lessonUuid?: string;
        exerciseUuid?: string;
    }>();

    const {
        data: userCourseData,
        isLoading: isLoadingUserCourse,
        error: userCourseError,
        refetch: refetchUserCourse,
    } = useGetMeCoursesCourseUuid(courseUuid!, {
        query: {
            enabled: !!courseUuid,
        },
    });

    const {
        data: courseData,
        isLoading: isLoadingCourse,
        error: courseError,
    } = useGetCoursesUuid(courseUuid!, {
        language: i18n.language,
    }, {
        query: {
            enabled: !!courseUuid,
        },
    });

    // Look up selected exercise from courseData
    const selectedExerciseData = useMemo(() => {
        if (!courseData || !lessonUuid || !exerciseUuid) return undefined;
        const lesson = courseData.lessons?.find((l) => l.uuid === lessonUuid);
        return lesson?.exercises?.find((exercise) => exercise.uuid === exerciseUuid);
    }, [courseData, lessonUuid, exerciseUuid]);

    // Look up selected lesson from courseData
    const selectedLessonData = useMemo(() => {
        if (!courseData || !lessonUuid) return undefined;
        return courseData.lessons?.find((l) => l.uuid === lessonUuid);
    }, [courseData, lessonUuid]);

    useEffect(() => {
        if (userCourseError || courseError) {
            toast.error(t("common.error"));
        }
    }, [userCourseError, courseError, t]);

    // Auto-navigation logic based on completion status
    useEffect(() => {
        if (!userCourseData || isLoadingUserCourse || isLoadingCourse) return;

        const lessons = userCourseData.lessons || [];

        if (lessons.length === 0) return;

        if (userCourseData.is_completed) {
            const firstLesson = lessons[0];
            if (firstLesson.lesson_uuid && firstLesson.exercises && firstLesson.exercises.length > 0) {
                const firstExercise = firstLesson.exercises[0];
                if (firstExercise.exercise_uuid) {
                    navigate(`/classroom/${courseUuid}/${firstLesson.lesson_uuid}/${firstExercise.exercise_uuid}`);
                    return;
                }
            }
        }

        const findFirstIncompleteExercise = (
            lessonExercises: ModelsUserExerciseStatus[] | undefined
        ) => {
            if (!lessonExercises || lessonExercises.length === 0) return null;

            const incomplete = lessonExercises.find((ex: ModelsUserExerciseStatus) => !ex.is_completed);
            if (incomplete) return incomplete;

            return lessonExercises[lessonExercises.length - 1];
        };

        // Only auto-navigate if no lesson or exercise is selected
        // Don't auto-navigate if user explicitly selected a lesson (lessonUuid but no exerciseUuid)
        if (!lessonUuid && !exerciseUuid) {
            for (const lesson of lessons) {
                if (!lesson.lesson_uuid) continue;

                const targetExercise = findFirstIncompleteExercise(lesson.exercises);
                if (targetExercise?.exercise_uuid) {
                    navigate(`/classroom/${courseUuid}/${lesson.lesson_uuid}/${targetExercise.exercise_uuid}`);
                    return;
                }
            }

            if (lessons.length > 0) {
                const lastLesson = lessons[lessons.length - 1];
                if (lastLesson.lesson_uuid && lastLesson.exercises && lastLesson.exercises.length > 0) {
                    const lastExercise = lastLesson.exercises[lastLesson.exercises.length - 1];
                    if (lastExercise.exercise_uuid) {
                        navigate(`/classroom/${courseUuid}/${lastLesson.lesson_uuid}/${lastExercise.exercise_uuid}`);
                        return;
                    }
                }
            }
        }

        // Removed auto-navigation when lessonUuid is present but exerciseUuid is not
        // This allows users to view lesson content without being redirected to an exercise

    }, [userCourseData, courseUuid, lessonUuid, exerciseUuid, navigate, isLoadingUserCourse, isLoadingCourse]);

    const handleAccordionChange = (value: string | undefined) => {
        if (!value || !userCourseData) return;

        const lesson = userCourseData.lessons?.find((l) => l.lesson_uuid === value);
        if (lesson?.exercises && lesson.exercises.length > 0) {
            // Find first incomplete exercise, or last if all complete
            const incomplete = lesson.exercises.find((ex) => !ex.is_completed);
            const targetExercise = incomplete || lesson.exercises[lesson.exercises.length - 1];

            if (targetExercise?.exercise_uuid) {
                navigate(
                    `/classroom/${courseUuid}/${value}/${targetExercise.exercise_uuid}`
                );
            }
        } else {
            navigate(`/classroom/${courseUuid}/${value}`);
        }
    };

    if (isLoadingUserCourse || isLoadingCourse) {
        return <ClassroomLoading />;
    }

    if (!userCourseData || !courseData) {
        return <ClassroomError />;
    }

    const safeUserCourseData: typeof userCourseData = userCourseData;
    const safeCourseData: typeof courseData = courseData;

    const breadcrumbs: BreadcrumbProps[] = [
        { label: "navigation.classroom", to: "/classroom" },
        { label: "navigation.myCourses", to: "/classroom/courses" },
        {
            label: safeCourseData.translation?.name ?? t("common.loading"),
            to: `/classroom/${courseUuid}`,
        },
    ];

    const onExerciseComplete = (exerciseUuid: string, nextLessonUuid?: string, nextExerciseUuid?: string) => {
        if (nextLessonUuid && nextExerciseUuid) {
            navigate(`/classroom/${courseUuid}/${nextLessonUuid}/${nextExerciseUuid}`);
        } else if (nextLessonUuid) {
            navigate(`/classroom/${courseUuid}/${nextLessonUuid}`);
        }

        const audio = new Audio(successSound);
        audio.play();

        refetchUserCourse();

        if (!nextExerciseUuid && !nextLessonUuid) {
            // No next exercise or lesson, meaning the course is complete
            setShowCourseCompletedDialog(true);
        }
    };

    const showCertificate = () => {
        console.log("TODO: Show certificate");
    };

    return (
        <div className="flex flex-col h-full gap-4">
            <PageHeader
                title={safeCourseData.translation?.name ?? t("common.loading")}
                breadcrumbs={breadcrumbs}
            />

            <div className="flex flex-1 gap-6 overflow-hidden">
                <LessonSidebar
                    userCourseData={safeUserCourseData}
                    courseData={safeCourseData}
                    courseUuid={courseUuid!}
                    selectedLessonUuid={lessonUuid}
                    selectedExerciseUuid={exerciseUuid}
                    onAccordionChange={handleAccordionChange}
                />

                <main className="flex-1 flex flex-col overflow-hidden">
                    {exerciseUuid ? (
                        <ExerciseContent
                            exercise={selectedExerciseData}
                            exerciseUuid={exerciseUuid}
                            language={courseData.subject}
                            onExerciseComplete={onExerciseComplete}
                        />
                    ) : selectedLessonData ? (
                        <LessonContent lesson={selectedLessonData} />
                    ) : (
                        <div className="flex flex-col h-full items-center justify-center">
                            <p className="text-muted-foreground">
                                {t("common.selectExercise") || "Select an exercise to begin"}
                            </p>
                        </div>
                    )}
                </main>
            </div>
            <Dialog open={showCourseCompletedDialog} onOpenChange={setShowCourseCompletedDialog}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader className="flex flex-col items-center">
                        <img src={courseCompleteImage} alt={t("common.courseComplete")} className="h-128 rounded-lg mb-4" />
                        <DialogTitle>
                            {t("course.courseComplete")}
                        </DialogTitle>
                        <DialogDescription>
                            {t("course.courseCompleteDescription")}
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter className="justify-center! mt-4">
                        <Button onClick={() => showCertificate()}>
                            {t("course.showCertificate")}
                        </Button>
                        <Button onClick={() => navigate("/classroom/courses")}>
                            {t("navigation.myCourses")}
                        </Button>
                        <Button variant="outline" onClick={() => setShowCourseCompletedDialog(false)}>
                            {t("common.close")}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}

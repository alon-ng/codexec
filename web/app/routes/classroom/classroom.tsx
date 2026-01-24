import { useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { useGetCoursesUuid } from "~/api/generated/courses/courses";
import { useGetMeCoursesCourseUuid } from "~/api/generated/me/me";
import type { MeUserExerciseStatus } from "~/api/generated/model";
import successSound from "~/assets/success.mp3";
import PageHeader, { type BreadcrumbProps } from "~/components/PageHeader";
import ClassroomError from "~/components/classroom/ClassroomError";
import ClassroomLoading from "~/components/classroom/ClassroomLoading";
import ExerciseContent from "~/components/classroom/ExerciseContent";
import LessonSidebar from "~/components/classroom/LessonSidebar";

export default function Classroom() {
    const { t, i18n } = useTranslation();
    const navigate = useNavigate();

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

        // Helper function to find first incomplete exercise in a lesson
        const findFirstIncompleteExercise = (
            lessonExercises: MeUserExerciseStatus[] | undefined
        ) => {
            if (!lessonExercises || lessonExercises.length === 0) return null;

            const incomplete = lessonExercises.find((ex: MeUserExerciseStatus) => !ex.is_completed);
            if (incomplete) return incomplete;

            // If all completed, return the last one
            return lessonExercises[lessonExercises.length - 1];
        };

        // Case 1: Only courseUuid provided - find first incomplete lesson and exercise
        if (!lessonUuid && !exerciseUuid) {
            for (const lesson of lessons) {
                if (!lesson.lesson_uuid) continue;

                const targetExercise = findFirstIncompleteExercise(lesson.exercises);
                if (targetExercise?.exercise_uuid) {
                    navigate(`/classroom/${courseUuid}/${lesson.lesson_uuid}/${targetExercise.exercise_uuid}`);
                    return;
                }
            }

            // If all lessons/exercises are completed, go to the last exercise of the last lesson
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

        // Case 2: courseUuid + lessonUuid provided - find first incomplete exercise in that lesson
        if (lessonUuid && !exerciseUuid) {
            const lesson = lessons.find((l) => l.lesson_uuid === lessonUuid);
            if (lesson) {
                const targetExercise = findFirstIncompleteExercise(lesson.exercises);
                if (targetExercise?.exercise_uuid) {
                    navigate(`/classroom/${courseUuid}/${lessonUuid}/${targetExercise.exercise_uuid}`);
                    return;
                }
            }
        }

        // Case 3: All params provided - do nothing, just load everything
    }, [userCourseData, courseUuid, lessonUuid, exerciseUuid, navigate, isLoadingUserCourse, isLoadingCourse]);

    // Handle accordion change - navigate to first exercise or lesson
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
                    <ExerciseContent
                        exercise={selectedExerciseData}
                        exerciseUuid={exerciseUuid}
                        language={courseData.subject}
                        onExerciseComplete={onExerciseComplete}
                    />
                </main>
            </div>
        </div>
    );
}

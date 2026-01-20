import type { MeUserCourseWithProgress } from "~/api/generated/model";

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { useTranslation } from "react-i18next";
import { courseIcons } from "~/utils/course";
import { Progress } from "../base/Progress";
import { Link } from "react-router";
import Blob from "~/assets/blob.svg?react";
import { Skeleton } from "../ui/skeleton";

export interface UserCourseCardProps {
    course?: MeUserCourseWithProgress;
}

export default function UserCourseCard({ course }: UserCourseCardProps) {
    const { t } = useTranslation();

    if (!course) {
        return (
            <Card className="relative h-72 bg-transparent overflow-hidden hover:shadow-md hover:cursor-pointer hover:translate-y-[-4px] transition-all duration-200">
                <CardHeader className="flex items-center gap-4 z-10">
                    <Skeleton className="h-10 w-10 rounded-full" />
                    <div className="space-y-2">
                        <Skeleton className="h-5 w-32" />
                        <Skeleton className="h-4 w-24" />
                    </div>
                </CardHeader>
                <CardContent className="flex-1 space-y-2 z-10">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-[90%]" />
                    <Skeleton className="h-4 w-[80%]" />
                </CardContent>
                <CardFooter className="flex items-center justify-between z-10">
                    <Skeleton className="h-4 w-16" />
                    <div className="flex items-center gap-2">
                        <Skeleton className="h-9 w-24" />
                        <Skeleton className="h-9 w-24" />
                    </div>
                </CardFooter>
                <Blob className="absolute -top-32 -start-24 size-128 z-0 blur-3xl text-codim-pink/5" />
                <Blob className="absolute -bottom-32 -end-24 size-128 z-0 blur-3xl text-codim-purple/5" />
            </Card>
        );
    }

    const nextExercisePath = `/classroom/${course.uuid}/${course.next_lesson_uuid}/${course.next_exercise_uuid}`;

    return (
        <Card className="relative h-72 bg-transparent overflow-hidden hover:shadow-md hover:cursor-pointer hover:translate-y-[-4px] transition-all duration-200">
            <CardHeader className="flex items-center gap-4 z-10">
                <img src={courseIcons[course.subject as keyof typeof courseIcons]} alt={course.subject} className="h-10 w-10" />
                <div>
                    <CardTitle className="text-lg leading-none">{course.translation.name}</CardTitle>
                    <div className="text-sm text-muted-foreground">
                        <span>{t("difficulty.title")}: </span>
                        {t(`difficulty.${course.difficulty}`)}
                    </div>
                </div>
            </CardHeader>
            <CardContent className="text-sm flex-1 z-10">
                <div>{course.translation.description}</div>
            </CardContent>
            <CardFooter className="flex flex-col items-start gap-1 z-10">
                <Progress value={course.completed_exercises / course.total_exercises * 100} showPercentage={true} />
                <div className="flex items-center gap-1 text-sm text-muted-foreground w-full">
                    <span className="shrink-0">{t("classroom.continueWhereYouLeftOff")}:</span>
                    <Link className="truncate underline hover:text-primary" to={nextExercisePath}>
                        {course.next_lesson_name} ({course.next_exercise_name})
                    </Link>
                </div>
            </CardFooter>
            <Blob className="absolute -top-32 -start-24 size-128 z-0 blur-3xl text-codim-pink/5" />
            <Blob className="absolute -bottom-32 -end-24 size-128 z-0 blur-3xl text-codim-purple/5" />
        </Card>
    );
}
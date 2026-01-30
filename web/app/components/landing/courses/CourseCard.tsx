import { ChevronsRightIcon, ShoppingCartIcon } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import type { ModelsCourseWithTranslation } from "~/api/generated/model";
import Blob from "~/assets/blob.svg?react";
import { Button } from "~/components/base/Button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { Skeleton } from "~/components/ui/skeleton";
import { courseIcons } from "~/utils/course";

export interface CourseCardProps {
    course?: ModelsCourseWithTranslation;
}

export default function CourseCard({ course }: CourseCardProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    if (!course) {
        return (
            <Card className="relative h-72 overflow-hidden">
                <CardHeader className="flex flex-row items-center gap-4 space-y-0 z-10">
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

    return (
        <Card className="relative h-72 overflow-hidden hover:shadow-md hover:cursor-pointer hover:translate-y-[-4px] transition-all duration-200" onClick={() => navigate(`/courses/${course.uuid}`)}>
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
            <CardFooter className="flex items-center justify-between z-10">
                {
                    course.discount! > 0 ? (
                        <div className="flex items-center gap-2">
                            <div className="text-sm text-gray-400 line-through leading-none">{course.price!}₪</div>
                            <div className="text-sm leading-none">{course.price! - course.discount!}₪</div>
                        </div>
                    ) : (
                        <div className="text-sm leading-none">{course.price!}₪</div>
                    )
                }
                <div className="flex items-center gap-2">
                    <Button variant="outline">
                        <ShoppingCartIcon className="h-4 w-4" />
                        {t("common.buy")}
                    </Button>
                    <Button variant="outline" onClick={() => navigate(`/courses/${course.uuid}`)}>
                        {t("common.readMore")}
                        <ChevronsRightIcon className="rtl:rotate-180 h-4 w-4" />
                    </Button>
                </div>
            </CardFooter>
            <Blob className="absolute -top-32 -start-24 size-128 z-0 blur-3xl text-codim-pink/5" />
            <Blob className="absolute -bottom-32 -end-24 size-128 z-0 blur-3xl text-codim-purple/5" />
        </Card>
    );
}
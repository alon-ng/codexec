
import type { DbCourse } from "~/api/generated/model";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { useTranslation } from "react-i18next";
import { courseIcons } from "~/utils/course";
import { Button } from "~/components/base/Button";
import { Skeleton } from "~/components/ui/skeleton";
import { ChevronsRightIcon, ShoppingCartIcon } from "lucide-react";
import { useNavigate } from "react-router";

export interface CourseCardProps {
    course?: DbCourse
}

export default function CourseCard({ course }: CourseCardProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    if (!course) {
        return (
            <Card className="h-72">
                <CardHeader className="flex flex-row items-center gap-4 space-y-0">
                    <Skeleton className="h-10 w-10 rounded-full" />
                    <div className="space-y-2">
                        <Skeleton className="h-5 w-32" />
                        <Skeleton className="h-4 w-24" />
                    </div>
                </CardHeader>
                <CardContent className="flex-1 space-y-2">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-[90%]" />
                    <Skeleton className="h-4 w-[80%]" />
                </CardContent>
                <CardFooter className="flex items-center justify-between">
                    <Skeleton className="h-4 w-16" />
                    <div className="flex items-center gap-2">
                        <Skeleton className="h-9 w-24" />
                        <Skeleton className="h-9 w-24" />
                    </div>
                </CardFooter>
            </Card>
        );
    }

    return (
        <Card className="h-72">
            <CardHeader className="flex items-center gap-4">
                <img src={courseIcons[course.subject as keyof typeof courseIcons]} alt={course.subject} className="h-10 w-10" />
                <div>
                    <CardTitle className="text-lg leading-none">{t(course.name!)}</CardTitle>
                    <div className="text-sm text-muted-foreground">
                        <span>{t("difficulty.title")}: </span>
                        {t(`difficulty.${course.difficulty}`)}
                    </div>
                </div>
            </CardHeader>
            <CardContent className="text-sm flex-1">
                <div>{t(course.description!)}</div>
            </CardContent>
            <CardFooter className="flex items-center justify-between">
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
        </Card>
    );
}
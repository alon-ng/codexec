import type { Route } from "./+types/course";
import PageHeader from "~/components/PageHeader";
import { useTranslation } from "react-i18next";
import { type BreadcrumbProps } from "~/components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { getCoursesUuid } from "~/api/generated/courses/courses";
import { useEffect } from "react";
import { toast } from "sonner";
import { Card, CardContent } from "~/components/ui/card";
import { Award, BadgeCheck, BookMarked, CircleCheckBig, Clock, ShoppingCartIcon } from "lucide-react";
import { Button } from "~/components/base/Button";
import { motion } from "motion/react";
import { blurInVariants } from "~/utils/animations";
import CoursePath from "~/components/landing/courses/CoursePath";

export default function Course({ params }: Route.ComponentProps) {
    const { t, i18n } = useTranslation();

    const { data, isLoading, error } = useQuery({
        queryKey: ['course', params.uuid],
        queryFn: () => getCoursesUuid(params.uuid, { language: i18n.language }),
    });

    const breadcrumbs: BreadcrumbProps[] = [
        { label: "navigation.home", to: "/" },
        { label: "navigation.courses", to: "/courses" },
        { label: data?.translation?.name || t("common.loading"), to: isLoading ? undefined : `/courses/${params.uuid}` },
    ]

    useEffect(() => {
        if (error) {
            toast.error(t("common.error"));
        }
    }, [error, t]);

    if (!data) {
        return <div className="flex flex-col h-full">
            <PageHeader title={t("common.loading")} breadcrumbs={breadcrumbs} />
        </div>
    }

    return (
        <div className="flex flex-col h-full gap-12">
            <PageHeader title={data.translation?.name!} breadcrumbs={breadcrumbs} />
            <div className="flex justify-between gap-4">
                <motion.div className="flex flex-col gap-4" variants={blurInVariants(0.1)} initial="hidden" animate="visible">
                    <div>{data.translation?.description}</div>
                    <div>
                        {data.translation?.bullets?.split("\n").map((bullet, index) => (
                            <div key={index} className="flex items-center gap-2">
                                <CircleCheckBig className="w-4 h-4" /> {bullet}</div>
                        ))}
                    </div>
                </motion.div>
                <motion.div variants={blurInVariants(0.2)} initial="hidden" animate="visible">
                    <Card className="w-64">
                        <CardContent className="flex flex-col gap-2">
                            <div className="flex gap-2">
                                <div className="text-2xl font-bold">₪{data.price! - data.discount!}</div>
                                <div className="text-xl text-gray-500 line-through">₪{data.price!}</div>
                            </div>

                            <div>
                                <div className="flex items-center gap-2 text-gray-500 text-sm">
                                    <Award className="w-4 h-4" /> {t("course.certificateNote")}
                                </div>
                                <div className="flex items-center gap-2 text-gray-500 text-sm">
                                    <BookMarked className="w-4 h-4" /> {data.lessons?.length ?? 0} {t("course.lessonsNote")}
                                </div>
                                <div className="flex items-center gap-2 text-gray-500 text-sm">
                                    <BadgeCheck className="w-4 h-4" /> {data.lessons?.map((lesson) => lesson.exercises?.length ?? 0).flat().reduce((a, b) => a + b, 0)} {t("course.exerciesNote")}
                                </div>
                                <div className="flex items-center gap-2 text-gray-500 text-sm">
                                    <Clock className="w-4 h-4" /> 15 {t("course.durationNote")}
                                </div>
                            </div>

                            <Button>
                                <ShoppingCartIcon className="w-4 h-4" />
                                {t("common.buy")}
                            </Button>
                        </CardContent>
                    </Card>
                </motion.div>
            </div>
            <CoursePath course={data} />
        </div >
    );
}
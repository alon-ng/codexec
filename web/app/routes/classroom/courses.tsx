import { useInfiniteQuery } from "@tanstack/react-query";
import { getMeCourses } from "~/api/generated/me/me";
import { useTranslation } from "react-i18next";
import PageHeader, { type BreadcrumbProps } from "~/components/PageHeader";
import { toast } from "sonner";
import UserCourseCard from "~/components/classroom/UserCourseCard";
import { useCallback, useEffect, useRef } from "react";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";

const breadcrumbs: BreadcrumbProps[] = [
  { label: "navigation.classroom", to: "/classroom" },
  { label: "navigation.myCourses", to: "/classroom/courses" },
];

const LIMIT = 9;

export default function UserCourses() {
  const { t, i18n } = useTranslation();
  const observerTarget = useRef<HTMLDivElement>(null);

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    status,
    error
  } = useInfiniteQuery({
    queryKey: ['courses', 'infinite'],
    queryFn: ({ pageParam = 0 }) => getMeCourses({ limit: LIMIT, offset: pageParam, language: i18n.language }),
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage && lastPage.length < LIMIT) return undefined;
      return allPages.length * LIMIT;
    },
    initialPageParam: 0,
  });

  const handleObserver = useCallback((entries: IntersectionObserverEntry[]) => {
    const [target] = entries;
    if (target.isIntersecting && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

  useEffect(() => {
    const element = observerTarget.current;
    const option = { threshold: 0.1 };
    const observer = new IntersectionObserver(handleObserver, option);

    if (element) observer.observe(element);

    return () => {
      if (element) observer.unobserve(element);
    };
  }, [handleObserver, observerTarget]);

  useEffect(() => {
    if (error) {
      toast.error(t("common.error"));
    }
  }, [error]);

  const courses = data?.pages.flatMap((page) => page) || [];
  const showSkeleton = status === 'pending';

  return (
    <>
      <PageHeader title={t("navigation.myCourses")} breadcrumbs={breadcrumbs} />
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4 pb-8">
        {courses.map((course, index) => (
          <motion.div key={`${course.uuid}-${index}`} variants={blurInVariants(index % LIMIT * 0.1)} initial="hidden" animate="visible">
            <UserCourseCard course={course} />
          </motion.div>
        ))}

        {(showSkeleton || isFetchingNextPage) && (
          Array.from({ length: 3 }).map((_, index) => (
            <motion.div key={`skeleton-${index}`} variants={blurInVariants(index * 0.1)} initial="hidden" animate="visible">
              <UserCourseCard />
            </motion.div>
          ))
        )}
      </div>

      {/* Sentinel element for infinite scroll */}
      <div ref={observerTarget} className="h-4 w-full" />
    </>
  );
}

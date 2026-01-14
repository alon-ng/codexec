import { Link } from "react-router";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "~/components/ui/breadcrumb";
import { Fragment } from "react";
import { useTranslation } from "react-i18next";
import { motion } from "motion/react";
import { blurInVariants } from "~/utils/animations";

export interface BreadcrumbProps {
    label: string;
    to?: string;
}

export interface PageHeaderProps {
    title: string;
    breadcrumbs?: BreadcrumbProps[];
}

export default function PageHeader({ title, breadcrumbs }: PageHeaderProps) {
    const { t } = useTranslation();

    return (
        <motion.div variants={blurInVariants()} initial="hidden" animate="visible" className="flex flex-col mb-4">
            {breadcrumbs && breadcrumbs.length > 0 && (
                <Breadcrumb>
                    <BreadcrumbList>
                        {breadcrumbs.map((crumb, index) => {
                            const isLast = index === breadcrumbs.length - 1;

                            return (
                                <Fragment key={crumb.label + index}>
                                    <BreadcrumbItem>
                                        {isLast || !crumb.to ? (
                                            <BreadcrumbPage>{t(crumb.label)}</BreadcrumbPage>
                                        ) : (
                                            <BreadcrumbLink asChild>
                                                <Link to={crumb.to}>{t(crumb.label)}</Link>
                                            </BreadcrumbLink>
                                        )}
                                    </BreadcrumbItem>
                                    {!isLast && <BreadcrumbSeparator />}
                                </Fragment>
                            );
                        })}
                    </BreadcrumbList>
                </Breadcrumb>
            )}
            <h1 className="text-3xl font-bold tracking-tight">{t(title)}</h1>
        </motion.div>
    );
}

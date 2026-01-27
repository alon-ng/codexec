import { GraduationCap } from "lucide-react";
import { motion } from "motion/react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import logoGradient from "~/assets/logo-gradient.svg";
import { Button } from "~/components/base/Button";
import { LanguageSelector } from "~/components/navbar/LanguageSelector";
import { UserMenu } from "~/components/navbar/UserMenu";
import {
    NavigationMenu,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
} from "~/components/ui/navigation-menu";
import { useLanguage } from "~/lib/useLanguage";
import { cn } from "~/lib/utils";
import { blurInVariants } from "~/utils/animations";

export interface PlatformNavbarNavigationItem {
    label: string;
    href: string;
}

export interface PlatformNavbarProps {
    showUserMenu?: boolean;
    showLoginButton?: boolean;
    showClassroomButton?: boolean;
    navigationItems: PlatformNavbarNavigationItem[];
}

export function PlatformNavbar({ navigationItems, showUserMenu = false, showLoginButton = false, showClassroomButton = false }: PlatformNavbarProps) {
    const { t } = useTranslation();
    const { dir } = useLanguage();

    return (
        <motion.nav variants={blurInVariants()} initial="hidden" animate="visible" dir={dir} className={cn(
            "absolute top-4 left-1/2 -translate-x-1/2 flex items-center justify-between z-50",
            "w-2/3 px-6 h-16 border rounded-lg shadow-md",
            "bg-background/50 backdrop-blur-sm"
        )}>

            <div className="flex gap-4">
                <Link to="/" className="flex items-center gap-2">
                    <img src={logoGradient} alt="Codim" className="h-10 w-10" />
                </Link>

                <NavigationMenu dir={dir}>
                    <NavigationMenuList className="gap-1">
                        {navigationItems.map((item) => (
                            <NavigationMenuItem key={item.href}>
                                <NavigationMenuLink
                                    asChild
                                    className="group inline-flex rounded-md text-sm h-9 px-4 py-2 transition-colors hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
                                >
                                    <Link to={item.href}>{t(item.label)}</Link>
                                </NavigationMenuLink>
                            </NavigationMenuItem>
                        ))}
                    </NavigationMenuList>
                </NavigationMenu>
            </div>
            <div className="flex items-center gap-2">
                <LanguageSelector />
                {showUserMenu && <UserMenu />}
                {showLoginButton &&
                    <>
                        <Button className="font-normal" variant="ghost">
                            <Link to="/login">{t("navigation.login")}</Link>
                        </Button>
                        <Button className="font-normal" variant="ghost" asChild>
                            <Link to="/login?signup=true">{t("navigation.signup")}</Link>
                        </Button>
                    </>}
                {showClassroomButton && <Button className="font-normal" variant="ghost">
                    <GraduationCap className="w-4 h-4" />
                    <Link to="/classroom">{t("navigation.classroom")}</Link>
                </Button>}
            </div>
        </motion.nav >
    );
}

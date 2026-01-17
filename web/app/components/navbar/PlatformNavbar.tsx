import { Link } from "react-router";
import logoGradient from "~/assets/logo-gradient.svg";
import {
    NavigationMenu,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
} from "~/components/ui/navigation-menu";
import { cn } from "~/lib/utils";
import { LanguageSelector } from "./LanguageSelector";
import { UserMenu } from "./UserMenu";
import { useTranslation } from "react-i18next";
import { motion } from "motion/react";
import { blurInVariants } from "~/utils/animations";
import { useLanguage } from "~/lib/useLanguage";

export interface PlatformNavbarNavigationItem {
    label: string;
    href: string;
}

export interface PlatformNavbarProps {
    showUserMenu?: boolean;
    navigationItems: PlatformNavbarNavigationItem[];
}

export function PlatformNavbar({ navigationItems, showUserMenu = false }: PlatformNavbarProps) {
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
                                    className="group inline-flex opacity-80 text-sm font-medium transition-colors hover:bg-transparent hover:opacity-100 focus:bg-transparent focus:text-accent-foreground focus:outline-none"
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
            </div>
        </motion.nav >
    );
}

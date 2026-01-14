import { PlatformNavbar } from "~/components/navbar/PlatformNavbar";
import { Outlet } from "react-router";


const navigationItems = [
    { label: "navigation.home", href: "/" },
    { label: "navigation.courses", href: "/courses" },
    { label: "navigation.about", href: "/about" },
    { label: "navigation.contact", href: "/contact" },
];

export default function LandingLayout() {
    return (
        <div className="flex h-screen flex-col bg-background overflow-hidden">
            <PlatformNavbar navigationItems={navigationItems} showUserMenu={false} />
            <main className="flex-1 pt-28 pb-8 px-24 overflow-y-auto">
                <Outlet />
            </main>
        </div>
    );
}
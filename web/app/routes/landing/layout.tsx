import { useEffect } from "react";
import { Outlet } from "react-router";
import { PlatformNavbar } from "~/components/navbar/PlatformNavbar";
import { useMeStore } from "~/stores/meStore";


const navigationItems = [
    { label: "navigation.home", href: "/" },
    { label: "navigation.courses", href: "/courses" },
    { label: "navigation.about", href: "/about" },
    { label: "navigation.contact", href: "/contact" },
];

export default function LandingLayout() {
    const { loadMe, userUUID, isLoggedIn } = useMeStore();

    useEffect(() => {
        if (!userUUID) {
            console.log("loading me");
            loadMe();
        }
    }, []);

    return (
        <div className="flex h-screen flex-col bg-background overflow-hidden">
            <PlatformNavbar navigationItems={navigationItems} showUserMenu={false} showLoginButton={!isLoggedIn} showClassroomButton={isLoggedIn} />
            <main className="flex-1 pt-28 pb-8 px-24 overflow-y-auto">
                <Outlet />
            </main>
        </div>
    );
}
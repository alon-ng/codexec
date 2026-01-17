import { useEffect } from "react";
import { Outlet } from "react-router";
import { useMeStore } from "~/stores/meStore";
import { PlatformNavbar } from "~/components/navbar/PlatformNavbar";

const navigationItems = [
  {
    label: "navigation.overview",
    href: "/classroom",
  },
  {
    label: "navigation.myCourses",
    href: "/classroom/courses",
  },
];

export default function PlatformLayout() {
  const { loadMe, isLoading, isLoggedIn } = useMeStore();

  useEffect(() => {
    if (!isLoggedIn) {
      loadMe();
    }
  }, []);

  // Show loading screen while authenticating
  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-100">
        <div className="text-center">
          <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-gray-300 border-r-gray-600"></div>
          <p className="text-gray-600">Authenticating...</p>
        </div>
      </div>
    );
  }

  if (!isLoggedIn) {
    return null;
  }

  return (
    <div className="flex h-screen flex-col bg-background overflow-hidden">
      <PlatformNavbar navigationItems={navigationItems} showUserMenu={true} />
      <main className="flex-1 pt-28 pb-8 px-24 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  );
}

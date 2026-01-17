/// <reference types="vite-plugin-svgr/client" />

import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useLoaderData,
} from "react-router";

import type { Route } from "./+types/root";
import "./app.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState, useEffect } from "react";
import { Toaster } from "sonner";
import { useTranslation } from "react-i18next";
import { initClientI18n } from "~/lib/i18n";
import "~/lib/i18n"; // Initialize client i18n
import { cn } from "./lib/utils";
import { detectLanguageFromRequest } from "~/lib/i18n.server";
import { useLanguage } from "~/lib/useLanguage";

// Loader runs on both server and client
export async function loader({ request }: Route.LoaderArgs) {
  const language = detectLanguageFromRequest(request);
  return { language };
}

export const links: Route.LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Open+Sans:ital,wght@0,300..800;1,300..800&display=swap",
  },
];

export function Layout({ children }: { children: React.ReactNode }) {
  const loaderData = useLoaderData<typeof loader>();
  const language = loaderData?.language ?? "en";

  // Initialize i18n with server-detected language synchronously
  // This ensures SSR and client hydration match
  initClientI18n(language);

  const { i18n } = useTranslation();
  const [lang, setLang] = useState(language);
  const isRTL = lang === "he";

  useEffect(() => {
    // Sync language on client-side changes
    const handleLanguageChanged = (lng: string) => {
      setLang(lng);
      if (typeof document !== "undefined") {
        document.documentElement.lang = lng;
        document.documentElement.dir = lng === "he" ? "rtl" : "ltr";
      }
    };

    // Set initial language and direction
    handleLanguageChanged(i18n.language || language);

    // Listen for language changes
    i18n.on("languageChanged", handleLanguageChanged);

    return () => {
      i18n.off("languageChanged", handleLanguageChanged);
    };
  }, [i18n, language]);

  return (
    <html lang={lang} dir={isRTL ? "rtl" : "ltr"}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body className={cn("overflow-hidden", isRTL ? "rtl" : "ltr")}>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  }));
  const { isRTL } = useLanguage();

  return (
    <QueryClientProvider client={queryClient}>
      <Outlet />
      <Toaster position={isRTL ? "top-left" : "top-right"} closeButton={true} />
    </QueryClientProvider>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  let message = "Oops!";
  let details = "An unexpected error occurred.";
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : "Error";
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}

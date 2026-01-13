import { useTranslation } from "react-i18next";
import { useRouteLoaderData, useMatches } from "react-router";
import { useState, useEffect } from "react";

/**
 * Hook to get the current language that matches SSR rendering
 * This prevents hydration mismatches by using the server-detected language
 */
export function useLanguage() {
  const { i18n } = useTranslation();
  
  // Try to get language from root loader data
  // In React Router v7, we can use useMatches to find the root route
  const matches = useMatches();
  const rootMatch = matches.find((match) => {
    // Root route is typically the first match or has id "root"
    return match.id === "root" || match.pathname === "/" || matches.indexOf(match) === 0;
  });
  const rootLoaderData = rootMatch?.data as { language?: string } | undefined;
  const serverLanguage = rootLoaderData?.language;
  
  // Use server language if available (for SSR consistency), otherwise use i18n.language
  // This ensures the initial render matches what the server rendered
  const [currentLang, setCurrentLang] = useState(() => {
    // Prioritize server language to match SSR, fallback to i18n.language
    return serverLanguage || i18n.language || "en";
  });

  useEffect(() => {
    // Sync with i18n changes after hydration
    const handleLanguageChanged = (lng: string) => {
      setCurrentLang(lng);
    };
    
    // Update if i18n language changes
    if (i18n.language && i18n.language !== currentLang) {
      setCurrentLang(i18n.language);
    }
    
    i18n.on("languageChanged", handleLanguageChanged);
    return () => {
      i18n.off("languageChanged", handleLanguageChanged);
    };
  }, [i18n]);

  return currentLang;
}

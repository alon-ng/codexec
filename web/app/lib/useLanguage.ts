import { useTranslation } from "react-i18next";
import { useMatches } from "react-router";
import { useState, useEffect } from "react";

type Direction = "ltr" | "rtl";

/**
 * Hook to get the current language that matches SSR rendering
 * This prevents hydration mismatches by using the server-detected language
 */
export function useLanguage() {
  const { i18n } = useTranslation();
  
  // Try to get language from root loader data
  const matches = useMatches();
  const rootMatch = matches.find((match) => {
    return match.id === "root" || match.pathname === "/" || matches.indexOf(match) === 0;
  });
  const rootLoaderData = rootMatch?.data as { language?: string } | undefined;
  const serverLanguage = rootLoaderData?.language;
  
  const [currentLang, setCurrentLang] = useState(() => {
    return serverLanguage || i18n.language || "en";
  });

  useEffect(() => {
    const handleLanguageChanged = (lng: string) => {
      setCurrentLang(lng);
    };
    
    if (i18n.language && i18n.language !== currentLang) {
      setCurrentLang(i18n.language);
    }
    
    i18n.on("languageChanged", handleLanguageChanged);
    return () => {
      i18n.off("languageChanged", handleLanguageChanged);
    };
  }, [i18n, currentLang]);

  const isRTL = currentLang === "he";
  const dir: Direction = isRTL ? "rtl" : "ltr";

  return {
    language: currentLang,
    isRTL,
    dir,
  };
}

/**
 * Convenience hook to get just the direction and RTL status
 */
export function useDirection() {
  const { isRTL, dir } = useLanguage();
  return { isRTL, dir };
}

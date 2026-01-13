import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import enTranslations from "~/locales/en.json";
import heTranslations from "~/locales/he.json";

const resources = {
  en: {
    translation: enTranslations,
  },
  he: {
    translation: heTranslations,
  },
};

// Client-side i18n instance - use default instance for react-i18next compatibility
// Initialize with a default language first
i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    lng: "en", // Default language for initial render
    fallbackLng: "en",
    supportedLngs: ["en", "he"],
    interpolation: {
      escapeValue: false,
    },
    react: {
      useSuspense: false, // Important for SSR compatibility
    },
    detection: {
      order: ["localStorage", "navigator"],
      caches: ["localStorage"],
      // Don't detect on init - we'll set it manually from server
      checkWhitelist: true,
    },
  });

/**
 * Initialize client i18n with a specific language
 * This is called to sync with server-rendered language
 * Can be called on both server and client
 * Sets the language synchronously to prevent hydration mismatches
 */
export function initClientI18n(language: string) {
  if ((language === "en" || language === "he") && i18n.language !== language) {
    // Set language synchronously to ensure SSR/client match
    i18n.language = language;
    // Also trigger changeLanguage to ensure all i18n internals are updated
    i18n.changeLanguage(language);
  }
}

export default i18n;

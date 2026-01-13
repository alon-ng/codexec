import i18n from "i18next";
import { initReactI18next } from "react-i18next";
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

// Server-side i18n instance (no language detector)
const serverI18n = i18n.createInstance();

serverI18n.use(initReactI18next).init({
  resources,
  fallbackLng: "en",
  supportedLngs: ["en", "he"],
  interpolation: {
    escapeValue: false,
  },
  react: {
    useSuspense: false, // Important for SSR
  },
});

export default serverI18n;

/**
 * Detect language from request headers/cookies
 * This runs on the server during SSR
 */
export function detectLanguageFromRequest(request: Request): string {
  // Check cookie first
  const cookieHeader = request.headers.get("Cookie");
  if (cookieHeader) {
    const cookies = cookieHeader.split(";").reduce((acc, cookie) => {
      const [key, value] = cookie.trim().split("=");
      acc[key] = value;
      return acc;
    }, {} as Record<string, string>);

    const langCookie = cookies["i18next"];
    if (langCookie === "en" || langCookie === "he") {
      return langCookie;
    }
  }

  // Check Accept-Language header
  const acceptLanguage = request.headers.get("Accept-Language");
  if (acceptLanguage) {
    // Parse Accept-Language header (e.g., "en-US,en;q=0.9,he;q=0.8")
    const languages = acceptLanguage
      .split(",")
      .map((lang) => {
        const [code] = lang.trim().split(";");
        return code.split("-")[0].toLowerCase();
      });

    // Check if Hebrew is preferred
    if (languages.includes("he")) {
      return "he";
    }
    // Check if English is preferred
    if (languages.includes("en")) {
      return "en";
    }
  }

  // Default to English
  return "en";
}

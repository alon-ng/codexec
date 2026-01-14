import { useTranslation } from "react-i18next";
import { useLanguage } from "~/lib/useLanguage";
import { Select } from "~/components/base/Select";

const languages = [
    { value: "en", label: "English", flag: "🇺🇸" },
    { value: "he", label: "עברית", flag: "🇮🇱" },
];

export function LanguageSelector() {
    const { i18n } = useTranslation();
    const { language: currentLang } = useLanguage();

    const switchLanguage = (langCode: string) => {
        i18n.changeLanguage(langCode);
        // Set cookie for SSR detection on next request
        if (typeof document !== "undefined") {
            document.cookie = `i18next=${langCode}; path=/; max-age=${60 * 60 * 24 * 365}; SameSite=Lax`;
        }
        // Redirect to the current page to ensure the new language is detected
        window.location.href = window.location.pathname;
    };

    return (
        <Select
            value={currentLang}
            onValueChange={switchLanguage}
            options={languages}
            triggerClassName="w-fit gap-1 px-2 border-none shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 cursor-pointer"
            contentClassName="min-w-[150px]"
            renderTrigger={(option) => (
                <div className="flex items-center gap-1 opacity-90 hover:opacity-100 transition-opacity duration-200">
                    <span className="text-md leading-none">{option?.flag}</span>
                    <span className="text-xs font-medium">{option?.value.toUpperCase()}</span>
                </div>
            )}
            renderOption={(option) => (
                <div className="flex items-center justify-between w-full gap-2">
                    <span className="text-lg leading-none">{option.flag}</span>
                    <span>{option.label}</span>
                </div>
            )}
        />
    );
}

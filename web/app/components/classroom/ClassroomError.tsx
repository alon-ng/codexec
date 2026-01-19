import { useTranslation } from "react-i18next";

export default function ClassroomError() {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col h-full items-center justify-center">
      <p className="text-muted-foreground">{t("common.error")}</p>
    </div>
  );
}

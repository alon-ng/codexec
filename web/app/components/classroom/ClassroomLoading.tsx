import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";

export default function ClassroomLoading() {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col h-full items-center justify-center">
      <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      <p className="mt-4 text-muted-foreground">{t("common.loading")}</p>
    </div>
  );
}

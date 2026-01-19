import { useTranslation } from "react-i18next";
import { toast } from "sonner";

export interface ExerciseActionsProps {
  onSubmit?: () => void;
}

export default function ExerciseActions({ onSubmit }: ExerciseActionsProps) {
  const { t } = useTranslation();

  const handleSubmit = () => {
    if (onSubmit) {
      onSubmit();
    } else {
      // Default behavior
      toast.info(t("common.comingSoon") || "Coming soon!");
    }
  };

  return (
    <div className="flex justify-end gap-2">
      <button
        className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
        onClick={handleSubmit}
      >
        {t("common.submit")}
      </button>
    </div>
  );
}

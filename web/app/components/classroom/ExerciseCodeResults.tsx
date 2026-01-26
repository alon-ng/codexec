import { Ban, CheckCircle, ChevronRightSquare, FlaskConical, XCircle } from "lucide-react";
import { motion } from "motion/react";
import { useTranslation } from "react-i18next";
import type { ExecuteResponse } from "~/api/types";
import { useLanguage } from '~/lib/useLanguage';
import { cn } from '~/lib/utils';
import { blurInVariants } from "~/utils/animations";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";

export interface ExerciseCodeResultsProps {
  resultTab: string;
  setResultTab: (tab: string) => void;
  lastResult?: ExecuteResponse | null;
}

export default function ExerciseCodeResults({
  resultTab,
  setResultTab,
  lastResult,
}: ExerciseCodeResultsProps) {
  const { t } = useTranslation();
  const { dir } = useLanguage();

  if (!lastResult) {
    return null;
  }

  return (
    <motion.div
      variants={blurInVariants()}
      initial="hidden"
      animate="visible"
      className="flex h-48 text-xs overflow-auto"
    >
      <Tabs className="w-full" value={resultTab} onValueChange={setResultTab}>
        <TabsList className="h-8 w-full justify-start rounded-none p-1 border-t shadow-sm" dir={dir}>
          <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="console">
            <ChevronRightSquare className="size-3" />
            {t("common.console")}
          </TabsTrigger>
          <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="errors">
            <Ban className="size-3" />
            {t("common.errors")}
          </TabsTrigger>
          <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="tests">
            <FlaskConical className="size-3" />
            {t("common.tests")}
          </TabsTrigger>
        </TabsList>
        <TabsContent className="text-xs px-3 font-mono" value="console">
          <div className="whitespace-pre-wrap">
            {lastResult.stdout || <span className="text-muted-foreground">{t("common.noOutput") || "No output"}</span>}
          </div>
        </TabsContent>
        <TabsContent className="text-xs px-3 font-mono" value="errors">
          <div className="text-red-400 whitespace-pre-wrap">
            {lastResult.stderr || <span className="text-muted-foreground">{t("common.noErrors") || "No errors"}</span>}
          </div>
        </TabsContent>
        <TabsContent className="text-xs font-mono" value="tests">
          {lastResult.checker_results?.length === 0 ? <span className="text-muted-foreground">{t("common.noTests") || "No tests"}</span> : (
            <>
              {lastResult.checker_results?.map((result) => (
                <div className={cn("flex items-center gap-1.5 py-1 px-3", result.success ? "text-green-400 bg-green-50" : "text-red-400 bg-red-50")} key={result.type}>
                  {result.success ? <CheckCircle className="size-3" /> : <XCircle className="size-3" />}
                  <span>{result.message}</span>
                </div>
              ))}
            </>
          )}
        </TabsContent>
      </Tabs>
    </motion.div>
  );
}

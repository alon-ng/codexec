import { Ban, CheckCircle, ChevronRightSquare, Database, FlaskConical, UserRound, XCircle } from "lucide-react";
import { motion } from "motion/react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import Blob from "~/assets/blob.svg?react";
import codyAvatar from "~/assets/cody-256.png";
import javaIcon from "~/assets/java.png";
import nodejsIcon from "~/assets/nodejs.png";
import PythonIcon from "~/assets/python.svg?react";
import { Button } from "~/components/base/Button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card";
import { Label } from "~/components/ui/label";
import { RadioGroup, RadioGroupItem } from "~/components/ui/radio-group";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { cn } from "~/lib/utils";
import { blurInVariants } from "~/utils/animations";
import { prose } from '~/utils/prose';

export default function LandingPage() {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col gap-24 py-16">
      <section className="flex flex-col items-center gap-8 text-center">
        <motion.div
          variants={blurInVariants(0.1)}
          initial="hidden"
          animate="visible"
          className="flex flex-col gap-6 max-w-4xl"
        >
          <h1 className="text-5xl md:text-6xl lg:text-7xl font-bold leading-tight">
            <span className="bg-gradient-to-r from-[var(--color-codim-pink)] to-[var(--color-codim-purple)] bg-clip-text text-transparent">
              {t("landing.hero.title")}
            </span>
          </h1>
          <p className="text-xl md:text-2xl text-muted-foreground">
            {t("landing.hero.subtitle")}
          </p>
          <div className="flex justify-center gap-4 pt-4">
            <Button asChild size="lg" className="text-lg px-8">
              <Link to="/auth/login">{t("landing.hero.cta")}</Link>
            </Button>
          </div>
        </motion.div>

        <motion.div
          variants={blurInVariants(0.3)}
          initial="hidden"
          animate="visible"
          className="w-full max-w-5xl aspect-video bg-muted border-2 border-dashed border-muted-foreground/30 rounded-xl flex items-center justify-center"
        >
          {/* TODO: Replace with actual product demo video/animation */}
          <p className="text-muted-foreground text-lg">{t("landing.hero.demoPlaceholder")}</p>
        </motion.div>
      </section>

      <section className="flex flex-col gap-8">
        <motion.div
          variants={blurInVariants(0.2)}
          initial="hidden"
          animate="visible"
          className="text-center max-w-3xl mx-auto"
        >
          <h2 className="text-4xl font-bold mb-4">{t("landing.testFeedback.title")}</h2>
          <p className="text-lg text-muted-foreground">
            {t("landing.testFeedback.description")}
          </p>
        </motion.div>

        <motion.div
          variants={blurInVariants(0.3)}
          initial="hidden"
          animate="visible"
          className="max-w-4xl mx-auto w-full"
        >
          <Card dir="ltr" className="pt-0 relative overflow-hidden">
            <CardContent className="p-0 z-5">
              {/* Mock IDE Code Editor */}
              <div className="bg-[#1e1e1e] p-4 font-mono text-sm text-gray-300">
                <div className="flex items-center gap-2 mb-4">
                  <div className="w-3 h-3 rounded-full bg-red-500"></div>
                  <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                  <div className="w-3 h-3 rounded-full bg-green-500"></div>
                  <span className="ml-4 text-xs text-gray-500">main.py</span>
                </div>
                <pre className="text-xs">
                  {`def calculate_sum(numbers):
    sum = 0
    for num in numbers:
        sum += num
    return sum`}
                </pre>
              </div>

              {/* Mock Terminal with Tabs */}
              <div className="flex flex-col h-64 border-t">
                <Tabs defaultValue="tests" className="w-full h-full">
                  <TabsList className="h-10 w-full justify-start rounded-none p-1 border-b shadow-sm">
                    <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="console">
                      <ChevronRightSquare className="size-3" />
                      {t("landing.testFeedback.output")}
                    </TabsTrigger>
                    <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="errors">
                      <Ban className="size-3" />
                      {t("landing.testFeedback.errors")}
                    </TabsTrigger>
                    <TabsTrigger className="flex items-center gap-1.5 transition-colors cursor-pointer text-xs" value="tests">
                      <FlaskConical className="size-3" />
                      {t("landing.testFeedback.tests")}
                    </TabsTrigger>
                  </TabsList>
                  <TabsContent className="text-xs px-3 font-mono h-full overflow-auto m-0" value="console">
                    <div className="whitespace-pre-wrap py-2">
                      <span className="text-muted-foreground">{t("landing.testFeedback.noOutput")}</span>
                    </div>
                  </TabsContent>
                  <TabsContent className="text-xs px-3 font-mono h-full overflow-auto m-0" value="errors">
                    <div className="text-red-400 whitespace-pre-wrap py-2">
                      <span className="text-muted-foreground">{t("landing.testFeedback.noErrors")}</span>
                    </div>
                  </TabsContent>
                  <TabsContent className="text-xs font-mono h-full overflow-auto m-0" value="tests">
                    <div className="flex flex-col">
                      <div className={cn("flex items-center gap-1.5 py-2 px-3", "text-green-600 bg-green-50 dark:text-green-400 dark:bg-green-950/30")}>
                        <CheckCircle className="size-4" />
                        <span>test_calculate_sum_empty_list: PASSED</span>
                      </div>
                      <div className={cn("flex items-center gap-1.5 py-2 px-3", "text-green-600 bg-green-50 dark:text-green-400 dark:bg-green-950/30")}>
                        <CheckCircle className="size-4" />
                        <span>test_calculate_sum_single_element: PASSED</span>
                      </div>
                      <div className={cn("flex items-center gap-1.5 py-2 px-3", "text-green-600 bg-green-50 dark:text-green-400 dark:bg-green-950/30")}>
                        <CheckCircle className="size-4" />
                        <span>test_calculate_sum_multiple_elements: PASSED</span>
                      </div>
                      <div className={cn("flex items-center gap-1.5 py-2 px-3", "text-red-600 bg-red-50 dark:text-red-400 dark:bg-red-950/30")}>
                        <XCircle className="size-4" />
                        <span>test_calculate_sum_negative_numbers: FAILED - AssertionError: Expected -10, got 0</span>
                      </div>
                    </div>
                  </TabsContent>
                </Tabs>
              </div>
            </CardContent>
            <Blob className="absolute -top-32 -start-24 size-128 z-0 blur-3xl text-codim-pink/5" />
            <Blob className="absolute -bottom-32 -end-24 size-128 z-0 blur-3xl text-codim-purple/5" />
          </Card>
          <p className="text-center text-sm text-muted-foreground mt-4">
            {t("landing.testFeedback.caption")}
          </p>
        </motion.div>
      </section>

      {/* AI Mentor Feature Section */}
      <section className="flex flex-col md:flex-row gap-12 items-center">
        <motion.div
          variants={blurInVariants(0.2)}
          initial="hidden"
          animate="visible"
          className="flex-1 flex flex-col gap-4"
        >
          <h2 className="text-4xl font-bold mb-4">{t("landing.aiMentor.title")}</h2>
          <p className="text-lg text-muted-foreground">
            {t("landing.aiMentor.description")}
          </p>
        </motion.div>

        <motion.div
          variants={blurInVariants(0.3)}
          initial="hidden"
          animate="visible"
          className="flex-1 max-w-md"
        >
          <Card className="shadow-lg relative overflow-hidden">
            <CardContent className="p-4 flex flex-col gap-3 h-80 overflow-auto">
              {/* Mock Chat Messages */}
              <div className="flex gap-2 w-full flex-row-reverse">
                <UserRound className="size-9 bg-white rounded-full p-1 shadow-sm shrink-0 dark:bg-card" />
                <div className="text-xs font-medium shadow-sm rounded-lg p-3 bg-background max-w-[80%]">
                  {t("landing.aiMentor.userMessage1")}
                </div>
              </div>
              <div className="flex gap-2 w-full flex-row">
                <img src={codyAvatar} className="size-9 bg-white rounded-full p-1 shadow-sm shrink-0 dark:bg-card" alt={t("landing.aiMentor.altText")} />
                <div className={cn("text-xs font-medium shadow-sm rounded-lg p-3 bg-background max-w-[80%]", prose)}>
                  <p>{t("landing.aiMentor.botMessage1")}</p>
                  <p className="mt-2">{t("landing.aiMentor.botMessage2")} <code className="bg-muted px-1 rounded">{t("landing.aiMentor.botCode1")}</code></p>
                  <p className="mt-1">{t("landing.aiMentor.botMessage3")} <code className="bg-muted px-1 rounded">{t("landing.aiMentor.botCode2")}</code></p>
                </div>
              </div>
              <div className="flex gap-2 w-full flex-row-reverse">
                <UserRound className="size-9 bg-white rounded-full p-1 shadow-sm shrink-0 dark:bg-card" />
                <div className="text-xs font-medium shadow-sm rounded-lg p-3 bg-background max-w-[80%]">
                  {t("landing.aiMentor.userMessage2")}
                </div>
              </div>
            </CardContent>
            <Blob className="absolute -top-16 -end-16 size-64 z-0 blur-3xl text-codim-purple/10" />
          </Card>
        </motion.div>
      </section>

      {/* Quiz Feature Section */}
      <section className="flex flex-col gap-8">
        <motion.div
          variants={blurInVariants(0.2)}
          initial="hidden"
          animate="visible"
          className="text-center max-w-3xl mx-auto"
        >
          <h2 className="text-4xl font-bold mb-4">{t("landing.quiz.title")}</h2>
          <p className="text-lg text-muted-foreground">
            {t("landing.quiz.description")}
          </p>
        </motion.div>

        <motion.div
          variants={blurInVariants(0.3)}
          initial="hidden"
          animate="visible"
          className="max-w-2xl mx-auto w-full"
        >
          <Card className="relative overflow-hidden">
            <CardHeader>
              <CardTitle>{t("landing.quiz.question")}</CardTitle>
            </CardHeader>
            <CardContent>
              <RadioGroup defaultValue="b" disabled className="grid gap-3">
                <div className="relative border rounded-md p-4 flex items-center gap-3 h-full border-input">
                  <RadioGroupItem
                    id="quiz-answer-a"
                    value="a"
                  />
                  <Label htmlFor="quiz-answer-a" className="font-medium after:absolute after:inset-0 after:cursor-default">
                    {t("landing.quiz.answerA")}
                  </Label>
                </div>
                <div className={cn(
                  "relative border rounded-md p-4 flex items-center gap-3 h-full",
                  "border-green-500 bg-green-50 dark:bg-green-950/30"
                )}>
                  <RadioGroupItem
                    circleClassName="fill-green-500 stroke-green-500"
                    id="quiz-answer-b"
                    value="b"
                  />
                  <Label htmlFor="quiz-answer-b" className="font-medium after:absolute after:inset-0 after:cursor-default">
                    {t("landing.quiz.answerB")}
                  </Label>
                  <CheckCircle className="size-5 text-green-500 ml-auto" />
                </div>
                <div className="relative border rounded-md p-4 flex items-center gap-3 h-full border-input">
                  <RadioGroupItem
                    id="quiz-answer-c"
                    value="c"
                  />
                  <Label htmlFor="quiz-answer-c" className="font-medium after:absolute after:inset-0 after:cursor-default">
                    {t("landing.quiz.answerC")}
                  </Label>
                </div>
                <div className="relative border rounded-md p-4 flex items-center gap-3 h-full border-input">
                  <RadioGroupItem
                    id="quiz-answer-d"
                    value="d"
                  />
                  <Label htmlFor="quiz-answer-d" className="font-medium after:absolute after:inset-0 after:cursor-default">
                    {t("landing.quiz.answerD")}
                  </Label>
                </div>
              </RadioGroup>
            </CardContent>
            <Blob className="absolute -top-24 -start-24 size-96 z-0 blur-3xl text-codim-pink/5" />
          </Card>
        </motion.div>
      </section>

      {/* Curriculum Overview Section */}
      <section className="flex flex-col gap-8">
        <motion.div
          variants={blurInVariants(0.2)}
          initial="hidden"
          animate="visible"
          className="text-center max-w-3xl mx-auto"
        >
          <h2 className="text-4xl font-bold mb-4">{t("landing.curriculum.title")}</h2>
          <p className="text-lg text-muted-foreground">
            {t("landing.curriculum.description")}
          </p>
        </motion.div>

        <motion.div
          variants={blurInVariants(0.3)}
          initial="hidden"
          animate="visible"
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-7xl mx-auto w-full"
        >
          <Card className="relative overflow-hidden hover:shadow-md transition-all duration-200">
            <CardHeader>
              <div className="flex items-center gap-4 mb-2">
                <div className="size-12 rounded-lg bg-gradient-to-br from-blue-500/10 to-yellow-400/10 flex items-center justify-center border border-yellow-500">
                  <PythonIcon className="size-6" />
                </div>
                <CardTitle className="text-xl">{t("landing.curriculum.python.title")}</CardTitle>
              </div>
              <CardDescription>
                {t("landing.curriculum.python.description")}
              </CardDescription>
            </CardHeader>
            <Blob className="absolute -top-16 -start-16 size-64 z-0 blur-3xl text-codim-pink/5" />
          </Card>

          <Card className="relative overflow-hidden hover:shadow-md transition-all duration-200">
            <CardHeader>
              <div className="flex items-center gap-4 mb-2">
                <div className="size-12 rounded-lg bg-gradient-to-r from-green-500/5 to-green-600/5 flex items-center justify-center border border-green-500">
                  <img src={nodejsIcon} alt="Node.js" className="size-6" />
                </div>
                <CardTitle className="text-xl">{t("landing.curriculum.nodejs.title")}</CardTitle>
              </div>
              <CardDescription>
                {t("landing.curriculum.nodejs.description")}
              </CardDescription>
            </CardHeader>
            <Blob className="absolute -top-16 -end-16 size-64 z-0 blur-3xl text-codim-purple/5" />
          </Card>

          <Card className="relative overflow-hidden hover:shadow-md transition-all duration-200">
            <CardHeader>
              <div className="flex items-center gap-4 mb-2">
                <div className="size-12 rounded-lg bg-gradient-to-r from-red-500/5 to-red-600/5 flex items-center justify-center border border-red-500">
                  <img src={javaIcon} alt="Java" className="size-6" />
                </div>
                <CardTitle className="text-xl">{t("landing.curriculum.java.title")}</CardTitle>
              </div>
              <CardDescription>
                {t("landing.curriculum.java.description")}
              </CardDescription>
            </CardHeader>
            <Blob className="absolute -bottom-16 -start-16 size-64 z-0 blur-3xl text-codim-pink/5" />
          </Card>

          <Card className="relative overflow-hidden hover:shadow-md transition-all duration-200">
            <CardHeader>
              <div className="flex items-center gap-4 mb-2">
                <div className="size-12 rounded-lg bg-gradient-to-r from-blue-500/5 to-blue-600/5 flex items-center justify-center border border-blue-500">
                  <Database className="size-6 text-blue-500" />
                </div>
                <CardTitle className="text-xl">{t("landing.curriculum.sql.title")}</CardTitle>
              </div>
              <CardDescription>
                {t("landing.curriculum.sql.description")}
              </CardDescription>
            </CardHeader>
            <Blob className="absolute -bottom-16 -end-16 size-64 z-0 blur-3xl text-codim-purple/5" />
          </Card>
        </motion.div>
      </section>
    </div>
  );
}

import { Link } from "react-router";
import { useState } from "react";
import { Button } from "~/components/base/Button";
import { Input } from "~/components/base/Input";
import { usePostAuthLogin } from "~/api/generated/auth/auth";
import { toast } from "sonner";
import { blurInVariants } from "~/utils/animations";
import { motion } from "motion/react";
import { useTranslation } from "react-i18next";

interface LoginFormProps {
    onSuccess?: () => void;
}

export function LoginForm({ onSuccess }: LoginFormProps) {
    const { t } = useTranslation();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const loginMutation = usePostAuthLogin({
        mutation: {
            onSuccess: () => {
                toast.success(t("auth.loginSuccessful"));
                onSuccess?.();
            },
            onError: (error) => {
                toast.error(error.error || t("auth.loginError"), {
                    className: "error-toast",
                });
            },
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        loginMutation.mutate({
            data: {
                email,
                password,
            },
        });
    };

    return (
        <motion.div className="flex flex-col gap-12" variants={blurInVariants()} initial="hidden" animate="visible">
            <div>
                <div className="text-4xl font-bold">
                    {t("auth.welcomeBack")}
                </div>
                <div className="text-sm text-gray-500">
                    {t("auth.signInToAccount")}
                </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4 ">
                <Input
                    label={t("common.email")}
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    placeholder={t("auth.emailPlaceholder")}
                />
                <Input
                    label={t("common.password")}
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    placeholder={t("auth.passwordPlaceholder")}
                />
                <div className="flex justify-end rtl:justify-start items-center">
                    <Link to="/forgot-password" className="text-xs text-gray-500 hover:underline">{t("auth.forgotPassword")}</Link>
                </div>
                <Button type="submit" className="w-full" isLoading={loginMutation.isPending} disabled={loginMutation.isPending}>
                    {t("common.signIn")}
                </Button>

                <div className="flex items-center justify-center">
                    <div className="flex-1 h-px bg-gray-200" />
                    <div className="text-xs px-2 text-gray-400">{t("auth.signInWith")}</div>
                    <div className="flex-1 h-px bg-gray-200" />
                </div>

                <Button type="button" className="w-full" variant="outline">
                    {t("common.google")}
                </Button>
            </form>

            <div className="flex items-center justify-center mt-12">
                <div className="text-xs text-gray-500 me-1">{t("auth.dontHaveAccount")}</div>
                <Link to="/login?signup=true" className="text-xs text-gray-500 hover:underline">{t("common.signUp")}</Link>
            </div>
        </motion.div>
    );
}

import { motion } from "motion/react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import { toast } from "sonner";
import { usePostAuthSignup } from "~/api/generated/auth/auth";
import { Button } from "~/components/base/Button";
import { Input } from "~/components/base/Input";
import { blurInVariants } from "~/utils/animations";

interface SignupFormProps {
    onSuccess?: () => void;
}

export function SignupForm({ onSuccess }: SignupFormProps) {
    const { t } = useTranslation();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");

    const signupMutation = usePostAuthSignup({
        mutation: {
            onSuccess: () => {
                toast.success(t("auth.signupSuccessful"));
                onSuccess?.();
            },
            onError: (error) => {
                toast.error(error.error || t("auth.signupError"), {
                    className: "error-toast",
                });
            },
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        signupMutation.mutate({
            data: {
                email,
                password,
                first_name: firstName,
                last_name: lastName,
            },
        });
    };

    return (
        <motion.div className="flex flex-col gap-12" variants={blurInVariants()} initial="hidden" animate="visible">
            <div>
                <div className="text-4xl font-bold">
                    {t("auth.createAccount")}
                </div>
                <div className="text-sm text-gray-500">
                    {t("auth.signUpToGetStarted")}
                </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                    <Input
                        label={t("common.firstName")}
                        type="text"
                        value={firstName}
                        onChange={(e) => setFirstName(e.target.value)}
                        required
                        placeholder={t("auth.firstNamePlaceholder")}
                    />
                    <Input
                        label={t("common.lastName")}
                        type="text"
                        value={lastName}
                        onChange={(e) => setLastName(e.target.value)}
                        required
                        placeholder={t("auth.lastNamePlaceholder")}
                    />
                </div>
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
                <Button type="submit" className="w-full" isLoading={signupMutation.isPending} disabled={signupMutation.isPending}>
                    {t("common.signUp")}
                </Button>

                <div className="flex items-center justify-center">
                    <div className="flex-1 h-px bg-gray-200" />
                    <div className="text-xs px-2 text-gray-400">{t("auth.signUpWith")}</div>
                    <div className="flex-1 h-px bg-gray-200" />
                </div>

                <Button type="button" className="w-full" variant="outline">
                    {t("common.google")}
                </Button>
            </form>

            <div className="flex items-center justify-center">
                <div className="text-xs text-gray-500 me-1">{t("auth.alreadyHaveAccount")}</div>
                <Link to="/login" className="text-xs text-gray-500 hover:underline">{t("common.signIn")}</Link>
            </div>
        </motion.div>
    );
}

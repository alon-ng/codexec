import { motion } from "motion/react";
import { useEffect } from "react";
import { Navigate, useNavigate, useSearchParams } from "react-router";
import loginImage from "~/assets/login.png";
import logoGradient from "~/assets/logo-gradient.svg";
import { LoginForm } from "~/components/login/LoginForm";
import { SignupForm } from "~/components/login/SignupForm";
import { useMeStore } from "~/stores/meStore";
import { blurInVariants } from "~/utils/animations";

export default function Login() {
  const navigate = useNavigate();

  const [searchParams] = useSearchParams();
  const isSignup = searchParams.get("signup") === "true";

  const { isLoggedIn, loadMe } = useMeStore();

  useEffect(() => {
    loadMe();
  }, []);

  if (isLoggedIn) {
    return <Navigate to="/classroom" />;
  }

  const handleSuccess = () => {
    navigate("/classroom");
  };

  return (
    <div className="flex p-3 h-screen w-screen">
      <div className="flex-1 flex flex-col justify-center xl:px-36 md:px-24 gap-12">
        <motion.div
          variants={blurInVariants()}
          initial="hidden"
          animate="visible"
          className="flex gap-2 items-center cursor-pointer w-fit"
          onClick={() => navigate("/")}
        >
          <img src={logoGradient} alt="Logo" className="w-8 h-8" />
          <h1 className="text-2xl font-bold">Codim</h1>
        </motion.div>

        {isSignup ? (
          <SignupForm onSuccess={handleSuccess} />
        ) : (
          <LoginForm onSuccess={handleSuccess} />
        )}
      </div>

      <img src={loginImage} alt="Login" className="h-full w-1/2 rounded-xl object-cover" />
    </div>
  );
}

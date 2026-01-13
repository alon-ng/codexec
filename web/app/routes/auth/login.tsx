import { useEffect } from "react";
import { Navigate, useNavigate, useSearchParams } from "react-router";
import logoGradient from "~/assets/logo-gradient.svg";
import loginImage from "~/assets/login.png";
import { LoginForm } from "~/components/login/LoginForm";
import { SignupForm } from "~/components/login/SignupForm";
import { useMeStore } from "~/stores/meStore";

export default function Login() {
  const navigate = useNavigate();

  const [searchParams] = useSearchParams();
  const isSignup = searchParams.get("signup") === "true";

  const { loadMe, isLoading, isLoggedIn } = useMeStore();

  useEffect(() => {
    if (!isLoggedIn) {
      loadMe();
    }
  }, []);

  // Show loading screen while authenticating
  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-100">
        <div className="text-center">
          <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-gray-300 border-r-gray-600"></div>
          <p className="text-gray-600">Authenticating...</p>
        </div>
      </div>
    );
  }

  if (isLoggedIn) {
    return <Navigate to="/platform" />;
  }

  const handleSuccess = () => {
    navigate("/platform");
  };

  return (
    <div className="flex p-3 h-screen w-screen">
      <div className="flex-1 flex flex-col justify-center xl:px-36 md:px-24 gap-12">
        <div className="flex gap-2 items-center">
          <img src={logoGradient} alt="Logo" className="w-8 h-8" />
          <h1 className="text-2xl font-bold">Codim</h1>
        </div>

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

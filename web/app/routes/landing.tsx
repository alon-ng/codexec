import { Link } from "react-router";
import { Button } from "~/components/base/Button";

export default function Landing() {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center items-center">
      <h1 className="text-5xl font-bold mb-8 text-gray-900">Codim</h1>
      <p className="text-xl text-gray-600 mb-12">Learn to code interactively.</p>
      <div className="flex gap-4">
        <Button asChild size="lg">
          <Link to="/login">Login</Link>
        </Button>
        <Button asChild variant="outline" size="lg">
          <Link to="/signup">Sign Up</Link>
        </Button>
      </div>
    </div>
  );
}

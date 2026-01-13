import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/landing.tsx"),
  route("login", "routes/auth/login.tsx"),
  route("platform", "routes/platform/layout.tsx", [
      index("routes/platform/dashboard.tsx"),
      // Add more platform routes here
  ]),
] satisfies RouteConfig;

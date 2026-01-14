import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
  layout("routes/landing/layout.tsx", [
    index("routes/landing/landing.tsx"),
    route("courses", "routes/landing/courses.tsx"),
  ]),
  route("login", "routes/auth/login.tsx"),
  route("platform", "routes/platform/layout.tsx", [
      index("routes/platform/dashboard.tsx"),
      // Add more platform routes here
  ]),
] satisfies RouteConfig;

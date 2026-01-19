import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
  layout("routes/landing/layout.tsx", [
    index("routes/landing/landing.tsx"),
    route("courses", "routes/landing/courses.tsx"),
    route("courses/:uuid", "routes/landing/course.tsx"),
  ]),
  route("login", "routes/auth/login.tsx"),
  route("classroom", "routes/classroom/layout.tsx", [
    index("routes/classroom/dashboard.tsx"),
    route("courses", "routes/classroom/courses.tsx"),
    route(":courseUuid/:lessonUuid?/:exerciseUuid?", "routes/classroom/classroom.tsx"),
  ]),
] satisfies RouteConfig;

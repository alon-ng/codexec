-- name: CreateUserCourse :one
INSERT INTO "user_courses" (
  "user_uuid", 
  "course_uuid", 
  "completed_at"
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateUserCourse :one
UPDATE "user_courses"
SET "user_uuid" = COALESCE(sqlc.narg('user_uuid'), "user_uuid"), 
    "course_uuid" = COALESCE(sqlc.narg('course_uuid'), "course_uuid"), 
    "completed_at" = COALESCE(sqlc.narg('completed_at'), "completed_at")
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteUserCourse :exec
DELETE FROM "user_courses"
WHERE "uuid" = $1;

-- name: getUserCourseFull :many
SELECT  courses.uuid            AS "course_uuid", 
        user_courses.uuid       AS "user_course_uuid",
        user_courses.started_at AS "user_course_started_at",
        user_courses.last_accessed_at AS "user_course_last_accessed_at",
        (user_courses.completed_at IS NOT NULL)::boolean AS "course_is_completed",
        user_courses.completed_at AS "course_completed_at",
        lessons.uuid            AS "lesson_uuid", 
        user_lessons.uuid       AS "user_lesson_uuid",
        user_lessons.started_at AS "user_lesson_started_at",
        user_lessons.last_accessed_at AS "user_lesson_last_accessed_at",
        (user_lessons.completed_at IS NOT NULL)::boolean AS "lesson_is_completed",
        user_lessons.completed_at AS "lesson_completed_at",
        exercises.uuid          AS "exercise_uuid", 
        user_exercises.uuid     AS "user_exercise_uuid",
        user_exercises.started_at AS "user_exercise_started_at",
        user_exercises.last_accessed_at AS "user_exercise_last_accessed_at",
        (user_exercises.completed_at IS NOT NULL)::boolean AS "exercise_is_completed",
        user_exercises.completed_at AS "exercise_completed_at"
FROM "courses"
LEFT JOIN "lessons"   ON "courses"."uuid" = "lessons"."course_uuid"   AND "lessons"."deleted_at" IS NULL
LEFT JOIN "exercises" ON "lessons"."uuid" = "exercises"."lesson_uuid" AND "exercises"."deleted_at" IS NULL
LEFT JOIN "user_courses" ON "courses"."uuid" = "user_courses"."course_uuid" AND "user_courses"."user_uuid" = $1
LEFT JOIN "user_lessons" ON "lessons"."uuid" = "user_lessons"."lesson_uuid" AND "user_lessons"."user_uuid" = $1
LEFT JOIN "user_exercises" ON "exercises"."uuid" = "user_exercises"."exercise_uuid" AND "user_exercises"."user_uuid" = $1
WHERE "courses"."uuid" = $2 AND "courses"."deleted_at" IS NULL
ORDER BY "lessons"."order_index" ASC, "exercises"."order_index" ASC;

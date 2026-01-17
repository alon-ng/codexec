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

-- name: ListUserCoursesWithProgress :many
WITH course_exercise_counts AS (
    SELECT 
        lessons.course_uuid,
        COUNT(exercises.uuid)::BIGINT AS total_exercises
    FROM lessons
    JOIN exercises ON lessons.uuid = exercises.lesson_uuid
    WHERE lessons.deleted_at IS NULL
      AND exercises.deleted_at IS NULL
    GROUP BY lessons.course_uuid
),
user_exercise_counts AS (
    SELECT 
        lessons.course_uuid,
        user_exercises.user_uuid,
        COUNT(user_exercises.uuid)::BIGINT AS completed_exercises
    FROM lessons
    JOIN exercises ON lessons.uuid = exercises.lesson_uuid
    JOIN user_exercises ON exercises.uuid = user_exercises.exercise_uuid
    WHERE lessons.deleted_at IS NULL
      AND exercises.deleted_at IS NULL
      AND user_exercises.completed_at IS NOT NULL
    GROUP BY lessons.course_uuid, user_exercises.user_uuid
),
next_lessons AS (
    SELECT DISTINCT ON (lessons.course_uuid)
        lessons.course_uuid,
        lessons.uuid AS lesson_uuid,
        lesson_translations.name AS lesson_name
    FROM lessons
    LEFT JOIN user_lessons ON lessons.uuid = user_lessons.lesson_uuid AND user_lessons.user_uuid = $1
    LEFT JOIN lesson_translations ON lessons.uuid = lesson_translations.lesson_uuid AND lesson_translations.language = $2
    WHERE lessons.deleted_at IS NULL
      AND (user_lessons.completed_at IS NULL OR user_lessons.uuid IS NULL)
    ORDER BY lessons.course_uuid, lessons.order_index ASC
),
next_exercises AS (
    SELECT DISTINCT ON (lessons.course_uuid)
        lessons.course_uuid,
        exercises.uuid AS exercise_uuid,
        exercise_translations.name AS exercise_name,
        lessons.uuid AS lesson_uuid,
        lesson_translations.name AS lesson_name
    FROM lessons
    JOIN exercises ON lessons.uuid = exercises.lesson_uuid
    LEFT JOIN user_exercises ON exercises.uuid = user_exercises.exercise_uuid AND user_exercises.user_uuid = $1
    LEFT JOIN exercise_translations ON exercises.uuid = exercise_translations.exercise_uuid AND exercise_translations.language = $2
    LEFT JOIN lesson_translations ON lessons.uuid = lesson_translations.lesson_uuid AND lesson_translations.language = $2
    WHERE lessons.deleted_at IS NULL
      AND exercises.deleted_at IS NULL
      AND (user_exercises.completed_at IS NULL OR user_exercises.uuid IS NULL)
    ORDER BY lessons.course_uuid, lessons.order_index ASC, exercises.order_index ASC
)
SELECT  courses.uuid                        AS "course_uuid", 
        courses.created_at                  AS "course_created_at", 
        courses.modified_at                 AS "course_modified_at", 
        courses.deleted_at                  AS "course_deleted_at", 
        courses.subject                     AS "course_subject", 
        courses.price                       AS "course_price", 
        courses.discount                    AS "course_discount", 
        courses.is_active                   AS "course_is_active", 
        courses.difficulty                  AS "course_difficulty",
        course_translations.uuid            AS "course_translation_uuid",
        course_translations.language        AS "course_translation_language",
        course_translations.name            AS "course_name", 
        course_translations.description     AS "course_description", 
        course_translations.bullets         AS "course_bullets", 
        user_courses.started_at             AS "user_course_started_at",
        user_courses.last_accessed_at       AS "user_course_last_accessed_at",
        user_courses.completed_at           AS "user_course_completed_at",
        COALESCE(course_exercise_counts.total_exercises, 0)::INTEGER   AS total_exercises,
        COALESCE(user_exercise_counts.completed_exercises, 0)::INTEGER AS completed_exercises,
        next_lessons.lesson_uuid            AS "next_lesson_uuid",
        next_lessons.lesson_name            AS "next_lesson_name",
        next_exercises.exercise_uuid        AS "next_exercise_uuid",
        next_exercises.exercise_name        AS "next_exercise_name"
FROM "user_courses"
JOIN "courses"                    ON "user_courses"."course_uuid" = "courses"."uuid" AND "courses"."deleted_at" IS NULL
JOIN "course_translations"        ON "courses"."uuid" = "course_translations"."course_uuid" AND "course_translations"."language" = $2
LEFT JOIN course_exercise_counts  ON "courses"."uuid" = course_exercise_counts.course_uuid
LEFT JOIN user_exercise_counts    ON "courses"."uuid" = user_exercise_counts.course_uuid AND "user_courses"."user_uuid" = user_exercise_counts.user_uuid
LEFT JOIN next_lessons            ON "courses"."uuid" = next_lessons.course_uuid
LEFT JOIN next_exercises          ON "courses"."uuid" = next_exercises.course_uuid
WHERE "user_courses"."user_uuid" = $1
AND   (sqlc.narg('subject')::text IS NULL OR "courses"."subject" = sqlc.narg('subject'))
AND   (sqlc.narg('is_active')::boolean IS NULL OR "courses"."is_active" = sqlc.narg('is_active'))
ORDER BY "user_courses"."last_accessed_at" DESC
LIMIT $3 OFFSET $4;
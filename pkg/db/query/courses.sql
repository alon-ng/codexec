-- name: GetCourse :one
SELECT * FROM "courses"
JOIN "course_translations" ON "courses"."uuid" = "course_translations"."course_uuid" AND "course_translations"."language" = $2
WHERE "courses"."uuid" = $1 AND "courses"."deleted_at" IS NULL 
LIMIT 1;

-- name: getCourseFull :many
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
        lessons.uuid                        AS "lesson_uuid", 
        lessons.created_at                  AS "lesson_created_at", 
        lessons.modified_at                 AS "lesson_modified_at", 
        lessons.deleted_at                  AS "lesson_deleted_at", 
        lessons.course_uuid                 AS "lesson_course_uuid", 
        lessons.order_index                 AS "lesson_order_index", 
        lessons.is_public                   AS "lesson_is_public", 
        lesson_translations.uuid            AS "lesson_translation_uuid",
        lesson_translations.language        AS "lesson_translation_language",
        lesson_translations.name            AS "lesson_name", 
        lesson_translations.description     AS "lesson_description", 
        exercises.uuid                      AS "exercise_uuid", 
        exercises.created_at                AS "exercise_created_at", 
        exercises.modified_at               AS "exercise_modified_at", 
        exercises.deleted_at                AS "exercise_deleted_at", 
        exercises.lesson_uuid               AS "exercise_lesson_uuid", 
        exercises.order_index               AS "exercise_order_index", 
        exercises.reward                    AS "exercise_reward", 
        exercises.type                      AS "exercise_type",
        exercises.code_data                 AS "exercise_code_data",
        exercises.quiz_data                 AS "exercise_quiz_data",
        exercise_translations.uuid          AS "exercise_translation_uuid",
        exercise_translations.language      AS "exercise_translation_language",
        exercise_translations.name          AS "exercise_name", 
        exercise_translations.description   AS "exercise_description",
        exercise_translations.code_data     AS "exercise_translation_code_data",
        exercise_translations.quiz_data     AS "exercise_translation_quiz_data"
FROM "courses"
JOIN "course_translations"        ON "courses"."uuid" = "course_translations"."course_uuid" AND "course_translations"."language" = $2
LEFT JOIN "lessons"               ON "courses"."uuid" = "lessons"."course_uuid"   AND "lessons"."deleted_at" IS NULL
LEFT JOIN "lesson_translations"   ON "lessons"."uuid" = "lesson_translations"."lesson_uuid" AND "lesson_translations"."language" = $2
LEFT JOIN "exercises"             ON "lessons"."uuid" = "exercises"."lesson_uuid" AND "exercises"."deleted_at" IS NULL
LEFT JOIN "exercise_translations" ON "exercises"."uuid" = "exercise_translations"."exercise_uuid" AND "exercise_translations"."language" = $2
WHERE "courses"."uuid" = $1 AND "courses"."deleted_at" IS NULL
ORDER BY "lessons"."order_index" ASC, "exercises"."order_index" ASC;

-- name: ListCourses :many
SELECT * FROM "courses"
JOIN "course_translations" ON "courses"."uuid" = "course_translations"."course_uuid" AND "course_translations"."language" = $3
WHERE "deleted_at" IS NULL
AND   (sqlc.narg('subject')::text IS NULL OR "subject" = sqlc.narg('subject'))
AND   (sqlc.narg('is_active')::boolean IS NULL OR "is_active" = sqlc.narg('is_active'))
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateCourse :one
INSERT INTO "courses" (
  "subject",
  "price",
  "discount",
  "is_active",
  "difficulty"
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateCourse :one
UPDATE "courses"
SET "subject" = COALESCE(sqlc.narg('subject'), "subject"), 
    "price" = COALESCE(sqlc.narg('price'), "price"),
    "discount" = COALESCE(sqlc.narg('discount'), "discount"),
    "is_active" = COALESCE(sqlc.narg('is_active'), "is_active"),
    "difficulty" = COALESCE(sqlc.narg('difficulty'), "difficulty"),
    "modified_at" = NOW()
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteCourse :exec
UPDATE "courses"
SET "deleted_at" = NOW()
WHERE "uuid" = $1;

-- name: HardDeleteCourse :exec
DELETE FROM "courses"
WHERE "uuid" = $1;

-- name: UndeleteCourse :exec
UPDATE "courses"
SET "deleted_at" = NULL
WHERE "uuid" = $1;

-- name: CountCourses :one
SELECT COUNT(*) FROM "courses"
WHERE "deleted_at" IS NULL;

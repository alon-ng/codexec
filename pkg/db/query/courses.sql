-- name: GetCourse :one
SELECT * FROM "courses"
WHERE "uuid" = $1 AND "deleted_at" IS NULL 
LIMIT 1;

-- name: GetCourseByName :one
SELECT * FROM "courses"
WHERE "name" = $1 AND "deleted_at" IS NULL 
LIMIT 1;

-- name: getCourseFull :many
SELECT  courses.uuid            AS "course_uuid", 
        courses.created_at      AS "course_created_at", 
        courses.modified_at     AS "course_modified_at", 
        courses.deleted_at      AS "course_deleted_at", 
        courses.name            AS "course_name", 
        courses.description     AS "course_description", 
        courses.subject         AS "course_subject", 
        courses.price           AS "course_price", 
        courses.discount        AS "course_discount", 
        courses.is_active       AS "course_is_active", 
        courses.difficulty      AS "course_difficulty", 
        courses.bullets         AS "course_bullets", 
        lessons.uuid            AS "lesson_uuid", 
        lessons.created_at      AS "lesson_created_at", 
        lessons.modified_at     AS "lesson_modified_at", 
        lessons.deleted_at      AS "lesson_deleted_at", 
        lessons.course_uuid     AS "lesson_course_uuid", 
        lessons.name            AS "lesson_name", 
        lessons.description     AS "lesson_description", 
        lessons.order_index     AS "lesson_order_index", 
        lessons.is_public       AS "lesson_is_public", 
        exercises.uuid          AS "exercise_uuid", 
        exercises.created_at    AS "exercise_created_at", 
        exercises.modified_at   AS "exercise_modified_at", 
        exercises.deleted_at    AS "exercise_deleted_at", 
        exercises.lesson_uuid   AS "exercise_lesson_uuid", 
        exercises.name          AS "exercise_name", 
        exercises.description   AS "exercise_description", 
        exercises.order_index   AS "exercise_order_index", 
        exercises.reward        AS "exercise_reward", 
        exercises.data          AS "exercise_data" 
FROM "courses"
LEFT JOIN "lessons"   ON "courses"."uuid" = "lessons"."course_uuid"   AND "lessons"."deleted_at" IS NULL
LEFT JOIN "exercises" ON "lessons"."uuid" = "exercises"."lesson_uuid" AND "exercises"."deleted_at" IS NULL
WHERE "courses"."uuid" = $1 AND "courses"."deleted_at" IS NULL
ORDER BY "lessons"."order_index" ASC, "exercises"."order_index" ASC;

-- name: ListCourses :many
SELECT * FROM "courses"
WHERE "deleted_at" IS NULL
AND   (sqlc.narg('subject')::text IS NULL OR "subject" = sqlc.narg('subject'))
AND   (sqlc.narg('is_active')::boolean IS NULL OR "is_active" = sqlc.narg('is_active'))
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateCourse :one
INSERT INTO "courses" (
  "name", 
  "description", 
  "subject",
  "price",
  "discount",
  "is_active",
  "difficulty",
  "bullets"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateCourse :one
UPDATE "courses"
SET "name" = COALESCE($2, "name"), 
    "description" = COALESCE($3, "description"), 
    "subject" = COALESCE($4, "subject"), 
    "price" = COALESCE($5, "price"),
    "discount" = COALESCE($6, "discount"),
    "is_active" = COALESCE($7, "is_active"),
    "difficulty" = COALESCE($8, "difficulty"),
    "bullets" = COALESCE($9, "bullets"),
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

-- name: GetLesson :one
SELECT * FROM "lessons"
JOIN "lesson_translations" ON "lessons"."uuid" = "lesson_translations"."lesson_uuid" AND "lesson_translations"."language" = $2
WHERE "lessons"."uuid" = $1 AND "lessons"."deleted_at" IS NULL 
LIMIT 1;

-- name: ListLessons :many
SELECT * FROM "lessons"
JOIN "lesson_translations" ON "lessons"."uuid" = "lesson_translations"."lesson_uuid" AND "lesson_translations"."language" = $3
WHERE "lessons"."deleted_at" IS NULL
AND   (sqlc.narg('course_uuid')::uuid IS NULL OR "course_uuid" = sqlc.narg('course_uuid'))
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateLesson :one
INSERT INTO "lessons" (
  "course_uuid", 
  "order_index",
  "is_public"
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateLesson :one
UPDATE "lessons"
SET "order_index" = COALESCE(sqlc.narg('order_index'), "order_index"),
    "is_public" = COALESCE(sqlc.narg('is_public'), "is_public"),
    "modified_at" = NOW()
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteLesson :exec
UPDATE "lessons"
SET "deleted_at" = NOW()
WHERE "uuid" = $1;

-- name: HardDeleteLesson :exec
DELETE FROM "lessons"
WHERE "uuid" = $1;

-- name: UndeleteLesson :exec
UPDATE "lessons"
SET "deleted_at" = NULL
WHERE "uuid" = $1;

-- name: CountLessons :one
SELECT COUNT(*) FROM "lessons"
WHERE "deleted_at" IS NULL;

-- name: CountLessonsByCourse :one
SELECT COUNT(*) FROM "lessons"
WHERE "course_uuid" = $1 AND "deleted_at" IS NULL;


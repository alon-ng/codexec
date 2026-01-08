-- name: GetLesson :one
SELECT * FROM "lessons"
WHERE "uuid" = $1 AND "deleted_at" IS NULL 
LIMIT 1;

-- name: ListLessons :many
SELECT * FROM "lessons"
WHERE "deleted_at" IS NULL
AND   (sqlc.narg('course_uuid')::uuid IS NULL OR "course_uuid" = sqlc.narg('course_uuid'))
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateLesson :one
INSERT INTO "lessons" (
  "course_uuid", 
  "name", 
  "description",
  "order_index",
  "is_public"
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateLesson :one
UPDATE "lessons"
SET "course_uuid" = COALESCE($2, "course_uuid"), 
    "name" = COALESCE($3, "name"), 
    "description" = COALESCE($4, "description"), 
    "order_index" = COALESCE($5, "order_index"),
    "is_public" = COALESCE($6, "is_public"),
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


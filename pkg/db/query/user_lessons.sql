-- name: CreateUserLesson :one
INSERT INTO "user_lessons" (
  "user_uuid", 
  "lesson_uuid", 
  "completed_at"
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserLesson :one
SELECT * FROM "user_lessons"
WHERE "uuid" = $1
LIMIT 1;

-- name: GetUserLessonByUserAndLesson :one
SELECT * FROM "user_lessons"
WHERE "user_uuid" = $1 AND "lesson_uuid" = $2
LIMIT 1;

-- name: UpdateUserLesson :one
UPDATE "user_lessons"
SET "user_uuid" = COALESCE($2, "user_uuid"), 
    "lesson_uuid" = COALESCE($3, "lesson_uuid"), 
    "completed_at" = COALESCE($4, "completed_at")
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteUserLesson :exec
DELETE FROM "user_lessons"
WHERE "uuid" = $1;

-- name: GetLessonTranslation :one
SELECT * FROM "lesson_translations"
JOIN "lessons" ON "lesson_translations"."lesson_uuid" = "lessons"."uuid"
WHERE "lesson_translations"."uuid" = $1
AND "lessons"."deleted_at" IS NULL
LIMIT 1;

-- name: CreateLessonTranslation :one
INSERT INTO "lesson_translations" (
  "lesson_uuid",
  "language",
  "name",
  "description"
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateLessonTranslation :one
UPDATE "lesson_translations"
SET "name" = COALESCE(sqlc.narg('name'), "name"),
    "description" = COALESCE(sqlc.narg('description'), "description")
FROM "lessons"
WHERE "lesson_translations"."lesson_uuid" = "lessons"."uuid" AND "lesson_translations"."language" = $2
AND "lessons"."uuid" = $1
RETURNING "lesson_translations".*;

-- name: DeleteLessonTranslation :exec
DELETE FROM "lesson_translations"
USING "lessons"
WHERE "lesson_translations"."lesson_uuid" = "lessons"."uuid"
AND "lessons"."uuid" = $1;

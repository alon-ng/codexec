-- name: GetCourseTranslation :one
SELECT * FROM "course_translations"
JOIN "courses" ON "course_translations"."course_uuid" = "courses"."uuid"
WHERE "course_translations"."uuid" = $1
AND "courses"."deleted_at" IS NULL
LIMIT 1;

-- name: CreateCourseTranslation :one
INSERT INTO "course_translations" (
  "course_uuid",
  "language",
  "name",
  "description",
  "bullets"
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateCourseTranslation :one
UPDATE "course_translations"
SET "name" = COALESCE(sqlc.narg('name'), "name"),
    "description" = COALESCE(sqlc.narg('description'), "description"),
    "bullets" = COALESCE(sqlc.narg('bullets'), "bullets")
FROM "courses"
WHERE "course_translations"."course_uuid" = "courses"."uuid"
AND "courses"."uuid" = $1 AND "course_translations"."language" = $2
RETURNING "course_translations".*;

-- name: DeleteCourseTranslation :exec
DELETE FROM "course_translations"
USING "courses"
WHERE "course_translations"."course_uuid" = "courses"."uuid"
AND "courses"."uuid" = $1;
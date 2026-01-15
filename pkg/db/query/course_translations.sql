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
SET "language" = COALESCE($2, "language"),
    "name" = COALESCE($3, "name"),
    "description" = COALESCE($4, "description"),
    "bullets" = COALESCE($5, "bullets")
FROM "courses"
WHERE "course_translations"."course_uuid" = "courses"."uuid"
AND "courses"."uuid" = $1
RETURNING "course_translations".*;

-- name: DeleteCourseTranslation :exec
DELETE FROM "course_translations"
USING "courses"
WHERE "course_translations"."course_uuid" = "courses"."uuid"
AND "courses"."uuid" = $1;
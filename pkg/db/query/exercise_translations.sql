-- name: GetExerciseTranslation :one
SELECT * FROM "exercise_translations"
JOIN "exercises" ON "exercise_translations"."exercise_uuid" = "exercises"."uuid"
WHERE "exercise_translations"."uuid" = $1
AND "exercises"."deleted_at" IS NULL
LIMIT 1;

-- name: CreateExerciseTranslation :one
INSERT INTO "exercise_translations" (
  "exercise_uuid",
  "language",
  "name",
  "description"
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateExerciseTranslation :one
UPDATE "exercise_translations"
SET "name" = COALESCE(sqlc.narg('name'), "name"),
    "description" = COALESCE(sqlc.narg('description'), "description")
FROM "exercises"
WHERE "exercise_translations"."exercise_uuid" = "exercises"."uuid" AND "exercise_translations"."language" = $2
AND "exercises"."uuid" = $1
RETURNING "exercise_translations".*;

-- name: DeleteExerciseTranslation :exec
DELETE FROM "exercise_translations"
USING "exercises"
WHERE "exercise_translations"."exercise_uuid" = "exercises"."uuid"
AND "exercises"."uuid" = $1;

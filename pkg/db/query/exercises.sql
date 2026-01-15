-- name: GetExercise :one
SELECT * FROM "exercises"
JOIN "exercise_translations" ON "exercises"."uuid" = "exercise_translations"."exercise_uuid" AND "exercise_translations"."language" = $2
WHERE "exercises"."uuid" = $1 AND "exercises"."deleted_at" IS NULL 
LIMIT 1;

-- name: ListExercises :many
SELECT * FROM "exercises"
JOIN "exercise_translations" ON "exercises"."uuid" = "exercise_translations"."exercise_uuid" AND "exercise_translations"."language" = $3
WHERE "exercises"."deleted_at" IS NULL
AND   (sqlc.narg('lesson_uuid')::uuid IS NULL OR "lesson_uuid" = sqlc.narg('lesson_uuid'))
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateExercise :one
INSERT INTO "exercises" (
  "lesson_uuid", 
  "order_index",
  "reward",
  "type",
  "data"
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateExercise :one
UPDATE "exercises"
SET "lesson_uuid" = COALESCE($2, "lesson_uuid"), 
    "order_index" = COALESCE($3, "order_index"),
    "reward" = COALESCE($4, "reward"),
    "type" = COALESCE($5, "type"),
    "data" = COALESCE($6, "data"),
    "modified_at" = NOW()
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteExercise :exec
UPDATE "exercises"
SET "deleted_at" = NOW()
WHERE "uuid" = $1;

-- name: HardDeleteExercise :exec
DELETE FROM "exercises"
WHERE "uuid" = $1;

-- name: UndeleteExercise :exec
UPDATE "exercises"
SET "deleted_at" = NULL
WHERE "uuid" = $1;

-- name: CountExercises :one
SELECT COUNT(*) FROM "exercises"
WHERE "deleted_at" IS NULL;


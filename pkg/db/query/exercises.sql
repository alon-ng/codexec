-- name: GetExercise :one
SELECT * FROM "exercises"
WHERE "uuid" = $1 AND "deleted_at" IS NULL 
LIMIT 1;

-- name: ListExercises :many
SELECT * FROM "exercises"
WHERE "deleted_at" IS NULL
AND   (sqlc.narg('lesson_uuid')::uuid IS NULL OR "lesson_uuid" = sqlc.narg('lesson_uuid'))
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateExercise :one
INSERT INTO "exercises" (
  "lesson_uuid", 
  "name", 
  "description",
  "order_index",
  "reward",
  "type",
  "data"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateExercise :one
UPDATE "exercises"
SET "lesson_uuid" = COALESCE($2, "lesson_uuid"), 
    "name" = COALESCE($3, "name"), 
    "description" = COALESCE($4, "description"), 
    "order_index" = COALESCE($5, "order_index"),
    "reward" = COALESCE($6, "reward"),
    "type" = COALESCE($7, "type"),
    "data" = COALESCE($8, "data"),
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


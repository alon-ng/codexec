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
  "code_data",
  "quiz_data",
  "io_checker",
  "code_checker"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateExercise :one
UPDATE "exercises"
SET "order_index" = COALESCE(sqlc.narg('order_index'), "order_index"),
    "reward" = COALESCE(sqlc.narg('reward'), "reward"),
    "type" = COALESCE(sqlc.narg('type'), "type"),
    "code_data" = COALESCE(sqlc.narg('code_data'), "code_data"),
    "quiz_data" = COALESCE(sqlc.narg('quiz_data'), "quiz_data"),
    "io_checker" = COALESCE(sqlc.narg('io_checker'), "io_checker"),
    "code_checker" = COALESCE(sqlc.narg('code_checker'), "code_checker"),
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

-- name: GetExerciseForSubmission :one
SELECT "courses"."subject", "exercises"."type", "exercises"."code_checker", "exercises"."io_checker" FROM "courses"
JOIN "lessons" ON "courses"."uuid" = "lessons"."course_uuid"
JOIN "exercises" ON "lessons"."uuid" = "exercises"."lesson_uuid"
WHERE "exercises"."uuid" = $1
LIMIT 1;


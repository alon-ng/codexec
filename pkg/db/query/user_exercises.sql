-- name: CreateUserExercise :one
INSERT INTO "user_exercises" (
  "user_uuid", 
  "exercise_uuid", 
  "submission",
  "attempts",
  "completed_at"
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserExercise :one
SELECT * FROM "user_exercises"
WHERE "uuid" = $1
LIMIT 1;

-- name: GetUserExerciseByUserAndExercise :one
SELECT * FROM "user_exercises"
WHERE "user_uuid" = $1 AND "exercise_uuid" = $2
LIMIT 1;

-- name: UpdateUserExercise :one
UPDATE "user_exercises"
SET "user_uuid" = COALESCE($2, "user_uuid"), 
    "exercise_uuid" = COALESCE($3, "exercise_uuid"), 
    "submission" = COALESCE($4, "submission"),
    "attempts" = COALESCE($5, "attempts"),
    "completed_at" = COALESCE($6, "completed_at")
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteUserExercise :exec
DELETE FROM "user_exercises"
WHERE "uuid" = $1;

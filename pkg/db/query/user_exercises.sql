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
WHERE "user_uuid" = $1 AND "exercise_uuid" = $2
LIMIT 1;

-- name: UpdateUserExerciseSubmission :one
UPDATE "user_exercises"
SET "submission" = COALESCE(sqlc.narg('submission'), "submission"),
    "last_accessed_at" = NOW()
FROM "exercises"
WHERE "user_uuid" = $1 
AND "exercise_uuid" = $2 
AND "exercises"."type" = $3 
AND "exercises"."uuid" = "user_exercises"."exercise_uuid"
RETURNING "user_exercises".*;

-- name: UpdateUserExerciseSubmissionWithAttempts :one
UPDATE "user_exercises"
SET "submission" = COALESCE(sqlc.narg('submission'), "submission"),
    "attempts" = "attempts" + 1,
    "last_accessed_at" = NOW()
WHERE "user_uuid" = $1 AND "exercise_uuid" = $2
RETURNING *;

-- name: CompleteUserExercise :one
UPDATE "user_exercises"
SET "completed_at" = NOW()
WHERE "user_uuid" = $1 AND "exercise_uuid" = $2
RETURNING *;

-- name: ResetUserExercise :one
UPDATE "user_exercises"
SET "submission" = '{}'::jsonb,
    "attempts" = 0,
    "completed_at" = NULL,
    "last_accessed_at" = NOW()
WHERE "user_uuid" = $1 AND "exercise_uuid" = $2
RETURNING *;


-- name: CreateChatMessage :one
INSERT INTO "chat_messages" (
  "exercise_uuid", 
  "user_uuid", 
  "role", 
  "content",
  "prompt_tokens",
  "completion_tokens"
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListChatMessages :many
SELECT * FROM "chat_messages"
WHERE "exercise_uuid" = $1 AND "user_uuid" = $2
ORDER BY "ts" ASC
LIMIT $3 OFFSET $4;
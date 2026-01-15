-- name: GetUser :one
SELECT * FROM "users"
WHERE "uuid" = $1 AND "deleted_at" IS NULL 
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "users"
WHERE "email" = $1 AND "deleted_at" IS NULL 
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "users"
WHERE "deleted_at" IS NULL
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO "users" (
  "first_name", 
  "last_name", 
  "email",
  "password_hash",
  "is_verified",
  "is_admin"
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE "users"
SET "first_name" = COALESCE(sqlc.narg('first_name'), "first_name"), 
    "last_name" = COALESCE(sqlc.narg('last_name'), "last_name"), 
    "email" = COALESCE(sqlc.narg('email'), "email"), 
    "is_verified" = COALESCE(sqlc.narg('is_verified'), "is_verified"), 
    "is_admin" = COALESCE(sqlc.narg('is_admin'), "is_admin"),
    "modified_at" = NOW()
WHERE "uuid" = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE "users"
SET "password_hash" = $2,
    "modified_at" = NOW()
WHERE "uuid" = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE "users"
SET "deleted_at" = NOW()
WHERE "uuid" = $1;

-- name: HardDeleteUser :exec
DELETE FROM "users"
WHERE "uuid" = $1;

-- name: UndeleteUser :exec
UPDATE "users"
SET "deleted_at" = NULL
WHERE "uuid" = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM "users"
WHERE "deleted_at" IS NULL;
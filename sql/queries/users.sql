-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
    uuid_generate_v4(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateEmailAndPassword :one
UPDATE users
SET email = $1, password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpgradeUserChirpyRedStatus :one
UPDATE users
SET is_chirpy_red = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;
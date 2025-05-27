-- name: CreateAuthToken :one
INSERT INTO auth_tokens (
    user_id,
    auth_token,
    refresh_token,
    user_agent,
    ip_address,
    auth_token_expires_at,
    refresh_token_expires_at,
    auth_token_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAuthTokenByID :one
SELECT * FROM auth_tokens
WHERE id = $1;

-- name: GetAuthTokenByHash :one
SELECT * FROM auth_tokens
WHERE auth_token_hash = $1;

-- name: GetAuthTokenByRefreshToken :one
SELECT * FROM auth_tokens
WHERE refresh_token = $1;

-- name: ListAuthTokensByUserID :many
SELECT * FROM auth_tokens
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAuthTokens :one
UPDATE auth_tokens
SET
  auth_token = $2,
  refresh_token = $3,
  auth_token_expires_at = $4,
  refresh_token_expires_at = $5,
  auth_token_hash = $6,
  revoked = $7
WHERE id = $1
RETURNING *;

-- name: RevokeAuthToken :exec
UPDATE auth_tokens
SET revoked = TRUE
WHERE id = $1;

-- name: DeleteAuthTokenByID :exec
DELETE FROM auth_tokens
WHERE id = $1;

-- name: DeleteAuthTokensByUserID :exec
DELETE FROM auth_tokens
WHERE user_id = $1;

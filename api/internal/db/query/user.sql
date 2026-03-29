-- name: UpsertUserBinding :exec
INSERT INTO users (chain_id, address, bound_to) VALUES ($1, $2, $3)
ON CONFLICT (chain_id, address) DO UPDATE SET bound_to = EXCLUDED.bound_to;

-- name: ClearUserBinding :exec
UPDATE users SET bound_to = '' WHERE address = $1 AND chain_id = $2;

-- name: UpsertUserRecipient :exec
INSERT INTO users (chain_id, address, recipient) VALUES ($1, $2, $3)
ON CONFLICT (chain_id, address) DO UPDATE SET recipient = EXCLUDED.recipient;

-- name: GetUser :one
SELECT address, bound_to, recipient, registered_at FROM users WHERE address = $1 AND chain_id = $2;

-- name: ListUsers :many
SELECT address, bound_to, recipient, registered_at FROM users WHERE chain_id = $1 ORDER BY registered_at DESC LIMIT $2 OFFSET $3;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users WHERE chain_id = $1;

-- name: UpsertUserBalance :exec
INSERT INTO user_balances (chain_id, user_address, total_allocated)
VALUES ($1, $2, $3)
ON CONFLICT (chain_id, user_address) DO UPDATE SET
  total_allocated = user_balances.total_allocated + EXCLUDED.total_allocated;

-- name: GetUserBalance :one
SELECT user_address, total_allocated FROM user_balances WHERE user_address = $1 AND chain_id = $2;

-- name: AddUserAllocated :exec
INSERT INTO user_balances (chain_id, user_address, total_allocated, updated_block) VALUES ($1, $2, $3, $4)
ON CONFLICT (chain_id, user_address) DO UPDATE SET total_allocated = user_balances.total_allocated + EXCLUDED.total_allocated, updated_block = EXCLUDED.updated_block;

-- name: SubtractUserAllocated :exec
UPDATE user_balances SET total_allocated = GREATEST(total_allocated - $3, 0), updated_block = $4 WHERE user_address = $2 AND chain_id = $1;

-- name: InitUserBalance :exec
INSERT INTO user_balances (chain_id, user_address, total_allocated, updated_block) VALUES ($1, $2, 0, 0)
ON CONFLICT (chain_id, user_address) DO NOTHING;

-- name: SetUserRegisteredAt :exec
INSERT INTO users (chain_id, address, registered_at) VALUES ($1, $2, $3)
ON CONFLICT (chain_id, address) DO UPDATE SET registered_at = EXCLUDED.registered_at
WHERE users.registered_at = 0;

-- name: GetUsersByBoundTo :many
SELECT address, bound_to, recipient, registered_at FROM users
WHERE bound_to = $1 AND chain_id = $2 ORDER BY address LIMIT 500;

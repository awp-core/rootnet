-- name: UpsertUserBinding :exec
INSERT INTO users (address, bound_to) VALUES ($1, $2)
ON CONFLICT (address) DO UPDATE SET bound_to = EXCLUDED.bound_to;

-- name: ClearUserBinding :exec
UPDATE users SET bound_to = '' WHERE address = $1;

-- name: UpsertUserRecipient :exec
INSERT INTO users (address, recipient) VALUES ($1, $2)
ON CONFLICT (address) DO UPDATE SET recipient = EXCLUDED.recipient;

-- name: GetUser :one
SELECT address, bound_to, recipient, registered_at FROM users WHERE LOWER(address) = LOWER($1);

-- name: ListUsers :many
SELECT address, bound_to, recipient, registered_at FROM users ORDER BY registered_at DESC LIMIT $1 OFFSET $2;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users;

-- name: UpsertUserBalance :exec
INSERT INTO user_balances (user_address, total_allocated)
VALUES ($1, $2)
ON CONFLICT (user_address) DO UPDATE SET
  total_allocated = user_balances.total_allocated + EXCLUDED.total_allocated;

-- name: GetUserBalance :one
SELECT user_address, total_allocated FROM user_balances WHERE user_address = $1;

-- name: AddUserAllocated :exec
UPDATE user_balances SET total_allocated = total_allocated + $2 WHERE user_address = $1;

-- name: SubtractUserAllocated :exec
UPDATE user_balances SET total_allocated = GREATEST(total_allocated - $2, 0) WHERE user_address = $1;

-- name: InitUserBalance :exec
INSERT INTO user_balances (user_address, total_allocated) VALUES ($1, 0)
ON CONFLICT (user_address) DO NOTHING;

-- name: SetUserRegisteredAt :exec
INSERT INTO users (address, registered_at) VALUES ($1, $2)
ON CONFLICT (address) DO UPDATE SET registered_at = EXCLUDED.registered_at
WHERE users.registered_at = 0;

-- name: GetUsersByBoundTo :many
SELECT address, bound_to, recipient, registered_at FROM users
WHERE bound_to = $1 ORDER BY address;

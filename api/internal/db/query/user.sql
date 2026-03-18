-- name: InsertUser :exec
INSERT INTO users (address, registered_at) VALUES ($1, $2)
ON CONFLICT (address) DO NOTHING;

-- name: GetUser :one
SELECT address, registered_at FROM users WHERE address = $1;

-- name: ListUsers :many
SELECT address, registered_at FROM users ORDER BY registered_at DESC LIMIT $1 OFFSET $2;

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

-- name: UpsertRewardRecipient :exec
INSERT INTO user_reward_recipients (user_address, recipient_address) VALUES ($1, $2)
ON CONFLICT (user_address) DO UPDATE SET recipient_address = EXCLUDED.recipient_address;

-- name: GetRewardRecipient :one
SELECT user_address, recipient_address FROM user_reward_recipients WHERE user_address = $1;

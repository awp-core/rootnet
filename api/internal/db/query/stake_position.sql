-- name: InsertStakePosition :exec
INSERT INTO stake_positions (chain_id, token_id, owner, amount, lock_end_time, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (chain_id, token_id) DO NOTHING;

-- name: UpdateStakePositionOwner :exec
UPDATE stake_positions SET owner = $3 WHERE chain_id = $1 AND token_id = $2;

-- name: UpdateStakePosition :exec
UPDATE stake_positions SET amount = $3, lock_end_time = $4 WHERE chain_id = $1 AND token_id = $2;

-- name: BurnStakePosition :exec
UPDATE stake_positions SET burned = TRUE WHERE chain_id = $1 AND token_id = $2;

-- name: GetStakePosition :one
SELECT * FROM stake_positions WHERE chain_id = $1 AND token_id = $2;

-- name: GetUserStakePositions :many
SELECT * FROM stake_positions WHERE chain_id = $1 AND owner = $2 AND burned = FALSE ORDER BY token_id LIMIT 500;

-- name: GetUserTotalStaked :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_positions WHERE chain_id = $1 AND owner = $2 AND burned = FALSE;

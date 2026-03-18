-- name: InsertStakePosition :exec
INSERT INTO stake_positions (token_id, owner, amount, lock_end_time, created_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (token_id) DO NOTHING;

-- name: UpdateStakePositionOwner :exec
UPDATE stake_positions SET owner = $2 WHERE token_id = $1;

-- name: UpdateStakePosition :exec
UPDATE stake_positions SET amount = $2, lock_end_time = $3 WHERE token_id = $1;

-- name: BurnStakePosition :exec
UPDATE stake_positions SET burned = TRUE WHERE token_id = $1;

-- name: GetStakePosition :one
SELECT * FROM stake_positions WHERE token_id = $1;

-- name: GetUserStakePositions :many
SELECT * FROM stake_positions WHERE owner = $1 AND burned = FALSE ORDER BY token_id;

-- name: GetUserTotalStaked :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_positions WHERE owner = $1 AND burned = FALSE;

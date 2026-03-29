-- name: InsertVanitySalt :exec
INSERT INTO vanity_salts (chain_id, salt, address)
VALUES ($1, $2, $3)
ON CONFLICT (chain_id, salt) DO NOTHING;

-- name: GetRandomAvailableSalt :one
SELECT salt, address FROM vanity_salts WHERE chain_id = $1 AND used = FALSE ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED;

-- name: ClaimRandomSalt :one
WITH locked AS (
  SELECT id FROM vanity_salts WHERE chain_id = sqlc.arg(chain_id)::BIGINT AND used = FALSE ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED
)
UPDATE vanity_salts SET used = TRUE
WHERE id = (SELECT id FROM locked)
RETURNING salt, address;

-- name: MarkSaltUsedByAddress :exec
UPDATE vanity_salts SET used = TRUE, subnet_id = $2
WHERE chain_id = $1 AND address = LOWER($3) AND used = FALSE;

-- name: ListAvailableSalts :many
SELECT salt, address FROM vanity_salts WHERE chain_id = $1 AND used = FALSE ORDER BY id LIMIT $2;

-- name: CountAvailableSalts :one
SELECT COUNT(*) FROM vanity_salts WHERE chain_id = $1 AND used = FALSE;

-- name: ListAllSalts :many
SELECT salt, address, used, subnet_id, created_at FROM vanity_salts WHERE chain_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: InsertVanitySalt :exec
INSERT INTO vanity_salts (salt, address)
VALUES ($1, $2)
ON CONFLICT (salt) DO NOTHING;

-- name: GetRandomAvailableSalt :one
SELECT salt, address FROM vanity_salts WHERE used = FALSE ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED;

-- name: ClaimRandomSalt :one
WITH locked AS (
  SELECT id FROM vanity_salts WHERE used = FALSE ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED
)
UPDATE vanity_salts SET used = TRUE
WHERE id = (SELECT id FROM locked)
RETURNING salt, address;

-- name: MarkSaltUsedByAddress :exec
UPDATE vanity_salts SET used = TRUE, subnet_id = $1
WHERE address = LOWER($2) AND used = FALSE;

-- name: ListAvailableSalts :many
SELECT salt, address FROM vanity_salts WHERE used = FALSE ORDER BY id LIMIT $1;

-- name: CountAvailableSalts :one
SELECT COUNT(*) FROM vanity_salts WHERE used = FALSE;

-- name: ListAllSalts :many
SELECT salt, address, used, subnet_id, created_at FROM vanity_salts ORDER BY id LIMIT $1 OFFSET $2;

-- name: InsertSubnet :exec
INSERT INTO subnets (subnet_id, chain_id, owner, name, symbol, subnet_contract, skills_uri, min_stake, alpha_token, lp_pool, status, created_at, immunity_ends_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'Pending', $11, $12)
ON CONFLICT (subnet_id) DO NOTHING;

-- name: GetSubnet :one
SELECT * FROM subnets WHERE subnet_id = $1;

-- name: ListSubnets :many
SELECT * FROM subnets WHERE chain_id = $1 AND burned = FALSE ORDER BY subnet_id DESC LIMIT $2 OFFSET $3;

-- name: ListSubnetsByStatus :many
SELECT * FROM subnets WHERE chain_id = $1 AND status = $2 AND burned = FALSE ORDER BY subnet_id DESC LIMIT $3 OFFSET $4;

-- name: ListAllSubnets :many
-- NOTE: No chain_id filter — returns subnets from ALL chains.
SELECT * FROM subnets WHERE burned = FALSE ORDER BY subnet_id DESC LIMIT $1 OFFSET $2;

-- name: ListAllSubnetsByStatus :many
-- NOTE: No chain_id filter — returns subnets from ALL chains with status filter.
SELECT * FROM subnets WHERE status = $1 AND burned = FALSE ORDER BY subnet_id DESC LIMIT $2 OFFSET $3;

-- name: UpdateSubnetLP :exec
UPDATE subnets SET lp_pool = $2 WHERE subnet_id = $1;

-- name: UpdateSubnetStatus :exec
UPDATE subnets SET status = $2 WHERE subnet_id = $1;

-- name: UpdateSubnetActivated :exec
UPDATE subnets SET status = 'Active', activated_at = $2 WHERE subnet_id = $1;

-- name: UpdateSubnetSkillsURI :exec
UPDATE subnets SET skills_uri = $2 WHERE subnet_id = $1;

-- name: UpdateSubnetMinStake :exec
UPDATE subnets SET min_stake = $2 WHERE subnet_id = $1;

-- name: UpdateSubnetMetadataURI :exec
UPDATE subnets SET metadata_uri = $2 WHERE subnet_id = $1;

-- name: GetSubnetSkills :one
SELECT skills_uri FROM subnets WHERE subnet_id = $1;

-- name: UpdateSubnetOwner :exec
UPDATE subnets SET owner = $2 WHERE subnet_id = $1;

-- name: UpdateSubnetBurned :exec
UPDATE subnets SET burned = TRUE WHERE subnet_id = $1;

-- name: GetActiveSubnets :many
SELECT * FROM subnets WHERE chain_id = $1 AND status = 'Active' AND burned = FALSE ORDER BY subnet_id;

-- name: ListActiveAlphaTokens :many
SELECT alpha_token FROM subnets WHERE chain_id = $1 AND status = 'Active' AND alpha_token != '';

-- name: CountAllSubnets :one
SELECT COUNT(*) FROM subnets WHERE burned = FALSE;

-- name: ListActiveAlphaTokensWithSubnetID :many
SELECT subnet_id, alpha_token FROM subnets WHERE chain_id = $1 AND status = 'Active' AND alpha_token != '';

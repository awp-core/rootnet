-- name: ListChains :many
SELECT * FROM chains WHERE status = 'active' ORDER BY chain_id;

-- name: GetChain :one
SELECT * FROM chains WHERE chain_id = $1;

-- name: InsertChain :exec
INSERT INTO chains (chain_id, name, rpc_url, dex, explorer, awp_registry, awp_token, awp_emission, staking_vault, stake_nft, subnet_nft, dao_address, lp_manager, pool_manager, deploy_block)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

-- name: DeleteChain :exec
UPDATE chains SET status = 'inactive' WHERE chain_id = $1;

-- name: ListAllChains :many
SELECT * FROM chains ORDER BY chain_id;

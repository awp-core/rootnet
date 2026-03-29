-- name: DeleteStakeAllocationsAfterBlock :exec
DELETE FROM stake_allocations WHERE updated_block > $1 AND chain_id = $2;

-- name: DeleteUserBalancesAfterBlock :exec
DELETE FROM user_balances WHERE updated_block > $1 AND chain_id = $2;

-- name: UpsertStakeAllocation :exec
INSERT INTO stake_allocations (chain_id, user_address, agent_address, subnet_id, amount, frozen, updated_block)
VALUES ($1, $2, $3, $4, $5, FALSE, $6)
ON CONFLICT (chain_id, user_address, agent_address, subnet_id) DO UPDATE SET
  amount = stake_allocations.amount + EXCLUDED.amount,
  updated_block = EXCLUDED.updated_block;

-- name: SubtractStakeAllocation :exec
UPDATE stake_allocations SET amount = GREATEST(amount - $5, 0), updated_block = $6
WHERE chain_id = $1 AND user_address = $2 AND agent_address = $3 AND subnet_id = $4;

-- name: SetStakeAllocationFrozen :exec
UPDATE stake_allocations SET frozen = $5
WHERE chain_id = $1 AND user_address = $2 AND agent_address = $3 AND subnet_id = $4;

-- name: FreezeAgentAllocations :exec
UPDATE stake_allocations SET frozen = TRUE
WHERE chain_id = $1 AND user_address = $2 AND agent_address = $3;

-- name: GetStakeAllocation :one
SELECT user_address, agent_address, subnet_id, amount, frozen FROM stake_allocations
WHERE chain_id = $1 AND user_address = $2 AND agent_address = $3 AND subnet_id = $4;

-- name: GetAllocationsByUser :many
SELECT user_address, agent_address, subnet_id, amount, frozen FROM stake_allocations
WHERE chain_id = $1 AND user_address = $2 AND amount > 0 ORDER BY subnet_id LIMIT $3 OFFSET $4;

-- name: GetAllocationsByAgent :many
SELECT user_address, agent_address, subnet_id, amount, frozen FROM stake_allocations
WHERE chain_id = $1 AND agent_address = $2 AND amount > 0 ORDER BY subnet_id;

-- name: GetAgentSubnetStake :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_allocations
WHERE chain_id = $1 AND agent_address = $2 AND subnet_id = $3 AND frozen = FALSE;

-- name: GetAgentSubnets :many
SELECT subnet_id, amount FROM stake_allocations
WHERE chain_id = $1 AND agent_address = $2 AND amount > 0 AND frozen = FALSE ORDER BY subnet_id LIMIT 500;

-- name: GetSubnetTotalStake :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_allocations
WHERE subnet_id = $1 AND frozen = FALSE;

-- name: GetFrozenByUser :many
SELECT user_address, agent_address, subnet_id, amount FROM stake_allocations
WHERE chain_id = $1 AND user_address = $2 AND frozen = TRUE AND amount > 0 LIMIT 500;

-- name: DeleteFrozenAllocations :exec
DELETE FROM stake_allocations WHERE chain_id = $1 AND user_address = $2 AND frozen = TRUE;

-- name: GetUsersWithFrozenAllocations :many
SELECT DISTINCT user_address FROM stake_allocations WHERE chain_id = $1 AND frozen = TRUE;

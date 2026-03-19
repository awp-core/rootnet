-- name: TruncateStakeAllocations :exec
TRUNCATE stake_allocations;

-- name: TruncateUserBalances :exec
TRUNCATE user_balances;

-- name: UpsertStakeAllocation :exec
INSERT INTO stake_allocations (user_address, agent_address, subnet_id, amount, frozen)
VALUES ($1, $2, $3, $4, FALSE)
ON CONFLICT (user_address, agent_address, subnet_id) DO UPDATE SET
  amount = stake_allocations.amount + EXCLUDED.amount;

-- name: SubtractStakeAllocation :exec
UPDATE stake_allocations SET amount = GREATEST(amount - $4, 0)
WHERE user_address = $1 AND agent_address = $2 AND subnet_id = $3;

-- name: SetStakeAllocationFrozen :exec
UPDATE stake_allocations SET frozen = $4
WHERE user_address = $1 AND agent_address = $2 AND subnet_id = $3;

-- name: FreezeAgentAllocations :exec
UPDATE stake_allocations SET frozen = TRUE
WHERE user_address = $1 AND agent_address = $2;

-- name: GetStakeAllocation :one
SELECT user_address, agent_address, subnet_id, amount, frozen FROM stake_allocations
WHERE user_address = $1 AND agent_address = $2 AND subnet_id = $3;

-- name: GetAllocationsByUser :many
SELECT user_address, agent_address, subnet_id, amount, frozen FROM stake_allocations
WHERE user_address = $1 AND amount > 0 ORDER BY subnet_id LIMIT $2 OFFSET $3;

-- name: GetAllocationsByAgent :many
SELECT user_address, agent_address, subnet_id, amount, frozen FROM stake_allocations
WHERE agent_address = $1 AND amount > 0 ORDER BY subnet_id;

-- name: GetAgentSubnetStake :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_allocations
WHERE agent_address = $1 AND subnet_id = $2 AND frozen = FALSE;

-- name: GetAgentSubnets :many
SELECT subnet_id, amount FROM stake_allocations
WHERE agent_address = $1 AND amount > 0 AND frozen = FALSE ORDER BY subnet_id;

-- name: GetSubnetTotalStake :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(78,0) AS total FROM stake_allocations
WHERE subnet_id = $1 AND frozen = FALSE;

-- name: GetFrozenByUser :many
SELECT user_address, agent_address, subnet_id, amount FROM stake_allocations
WHERE user_address = $1 AND frozen = TRUE AND amount > 0;

-- name: DeleteFrozenAllocations :exec
DELETE FROM stake_allocations WHERE user_address = $1 AND frozen = TRUE;

-- name: GetUsersWithFrozenAllocations :many
SELECT DISTINCT user_address FROM stake_allocations WHERE frozen = TRUE;

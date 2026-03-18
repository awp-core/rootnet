-- name: InsertAgent :exec
INSERT INTO agents (agent_address, owner_address, is_manager, removed) VALUES ($1, $2, FALSE, FALSE)
ON CONFLICT (agent_address) DO NOTHING;

-- name: UpsertAgent :exec
INSERT INTO agents (agent_address, owner_address, is_manager, removed) VALUES ($1, $2, FALSE, FALSE)
ON CONFLICT (agent_address) DO UPDATE SET owner_address = EXCLUDED.owner_address, removed = FALSE, removed_at = NULL;

-- name: UpdateAgentRemoved :exec
UPDATE agents SET removed = TRUE, removed_at = $2 WHERE agent_address = $1;

-- name: UpdateAgentManager :exec
UPDATE agents SET is_manager = $2 WHERE agent_address = $1;

-- name: GetAgent :one
SELECT agent_address, owner_address, is_manager, removed, removed_at FROM agents WHERE agent_address = $1;

-- name: GetAgentsByOwner :many
SELECT agent_address, owner_address, is_manager, removed, removed_at FROM agents
WHERE owner_address = $1 ORDER BY agent_address;

-- name: GetActiveAgentsByOwner :many
SELECT agent_address, owner_address, is_manager, removed, removed_at FROM agents
WHERE owner_address = $1 AND removed = FALSE ORDER BY agent_address;

-- name: LookupAgentOwner :one
SELECT owner_address FROM agents WHERE agent_address = $1 AND removed = FALSE;

-- name: InsertProposal :exec
INSERT INTO proposals (chain_id, proposal_id, proposer, description, status, votes_for, votes_against)
VALUES ($1, $2, $3, $4, $5, 0, 0)
ON CONFLICT (chain_id, proposal_id) DO NOTHING;

-- name: GetProposal :one
SELECT proposal_id, proposer, description, status, votes_for, votes_against FROM proposals WHERE chain_id = $1 AND proposal_id = $2;

-- name: ListProposals :many
SELECT proposal_id, proposer, description, status, votes_for, votes_against FROM proposals
WHERE chain_id = $1 ORDER BY proposal_id DESC LIMIT $2 OFFSET $3;

-- name: ListProposalsByStatus :many
SELECT proposal_id, proposer, description, status, votes_for, votes_against FROM proposals
WHERE chain_id = $1 AND status = $2 ORDER BY proposal_id DESC LIMIT $3 OFFSET $4;

-- name: UpdateProposalStatus :exec
UPDATE proposals SET status = $3 WHERE chain_id = $1 AND proposal_id = $2;

-- name: AddProposalVotesFor :exec
UPDATE proposals SET votes_for = votes_for + $3 WHERE chain_id = $1 AND proposal_id = $2;

-- name: AddProposalVotesAgainst :exec
UPDATE proposals SET votes_against = votes_against + $3 WHERE chain_id = $1 AND proposal_id = $2;

-- name: GetSyncState :one
SELECT contract_name, last_block FROM sync_states WHERE chain_id = $1 AND contract_name = $2;

-- name: UpsertSyncState :exec
INSERT INTO sync_states (chain_id, contract_name, last_block) VALUES ($1, $2, $3)
ON CONFLICT (chain_id, contract_name) DO UPDATE SET last_block = EXCLUDED.last_block;

-- name: UpsertIndexedBlock :exec
INSERT INTO indexed_blocks (chain_id, block_number, block_hash) VALUES ($1, $2, $3)
ON CONFLICT (chain_id, block_number) DO UPDATE SET block_hash = EXCLUDED.block_hash;

-- name: GetIndexedBlockHash :one
SELECT block_hash FROM indexed_blocks WHERE chain_id = $1 AND block_number = $2;

-- name: DeleteIndexedBlocksAfter :exec
DELETE FROM indexed_blocks WHERE chain_id = $1 AND block_number > $2;

-- name: PruneIndexedBlocks :exec
DELETE FROM indexed_blocks WHERE chain_id = $1 AND block_number < $2;

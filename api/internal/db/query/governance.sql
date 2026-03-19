-- name: InsertProposal :exec
INSERT INTO proposals (proposal_id, proposer, description, status, votes_for, votes_against)
VALUES ($1, $2, $3, $4, 0, 0)
ON CONFLICT (proposal_id) DO NOTHING;

-- name: GetProposal :one
SELECT proposal_id, proposer, description, status, votes_for, votes_against FROM proposals WHERE proposal_id = $1;

-- name: ListProposals :many
SELECT proposal_id, proposer, description, status, votes_for, votes_against FROM proposals
ORDER BY proposal_id DESC LIMIT $1 OFFSET $2;

-- name: ListProposalsByStatus :many
SELECT proposal_id, proposer, description, status, votes_for, votes_against FROM proposals
WHERE status = $1 ORDER BY proposal_id DESC LIMIT $2 OFFSET $3;

-- name: UpdateProposalStatus :exec
UPDATE proposals SET status = $2 WHERE proposal_id = $1;

-- name: AddProposalVotesFor :exec
UPDATE proposals SET votes_for = votes_for + $2 WHERE proposal_id = $1;

-- name: AddProposalVotesAgainst :exec
UPDATE proposals SET votes_against = votes_against + $2 WHERE proposal_id = $1;

-- name: GetSyncState :one
SELECT contract_name, last_block FROM sync_states WHERE contract_name = $1;

-- name: UpsertSyncState :exec
INSERT INTO sync_states (contract_name, last_block) VALUES ($1, $2)
ON CONFLICT (contract_name) DO UPDATE SET last_block = EXCLUDED.last_block;

-- name: UpsertIndexedBlock :exec
INSERT INTO indexed_blocks (block_number, block_hash) VALUES ($1, $2)
ON CONFLICT (block_number) DO UPDATE SET block_hash = EXCLUDED.block_hash;

-- name: GetIndexedBlockHash :one
SELECT block_hash FROM indexed_blocks WHERE block_number = $1;

-- name: DeleteIndexedBlocksAfter :exec
DELETE FROM indexed_blocks WHERE block_number > $1;

-- name: PruneIndexedBlocks :exec
DELETE FROM indexed_blocks WHERE block_number < $1;

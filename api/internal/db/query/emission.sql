-- name: UpsertEpoch :exec
INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission, dao_emission)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (chain_id, epoch_id) DO UPDATE SET
  start_time = EXCLUDED.start_time,
  daily_emission = EXCLUDED.daily_emission,
  dao_emission = COALESCE(EXCLUDED.dao_emission, epochs.dao_emission);

-- name: GetEpoch :one
SELECT epoch_id, start_time, daily_emission, dao_emission FROM epochs WHERE chain_id = $1 AND epoch_id = $2;

-- name: ListEpochs :many
SELECT epoch_id, start_time, daily_emission, dao_emission FROM epochs
WHERE chain_id = $1 ORDER BY epoch_id DESC LIMIT $2 OFFSET $3;

-- name: GetLatestEpoch :one
SELECT epoch_id, start_time, daily_emission, dao_emission FROM epochs
WHERE chain_id = $1 ORDER BY epoch_id DESC LIMIT 1;

-- name: UpdateEpochDAO :exec
INSERT INTO epochs (chain_id, epoch_id, start_time, daily_emission, dao_emission) VALUES ($1, $2, 0, 0, $3)
ON CONFLICT (chain_id, epoch_id) DO UPDATE SET dao_emission = EXCLUDED.dao_emission;

-- name: InsertRecipientAWPDistribution :exec
INSERT INTO recipient_awp_distributions (chain_id, epoch_id, recipient, awp_amount) VALUES ($1, $2, $3, $4)
ON CONFLICT (chain_id, epoch_id, recipient) DO NOTHING;

-- name: GetRecipientEarnings :many
SELECT epoch_id, recipient, awp_amount FROM recipient_awp_distributions
WHERE chain_id = $1 AND recipient = $2 ORDER BY epoch_id DESC LIMIT $3 OFFSET $4;

-- name: GetSubnetEarningsByID :many
SELECT r.epoch_id, r.recipient, r.awp_amount
FROM recipient_awp_distributions r
JOIN subnets s ON r.recipient = s.subnet_contract
WHERE s.subnet_id = $1
ORDER BY r.epoch_id DESC LIMIT $2 OFFSET $3;

-- name: GetEpochDistributions :many
SELECT epoch_id, recipient, awp_amount FROM recipient_awp_distributions
WHERE chain_id = $1 AND epoch_id = $2 ORDER BY recipient;

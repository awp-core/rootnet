-- name: UpsertEpoch :exec
INSERT INTO epochs (epoch_id, start_time, daily_emission, dao_emission)
VALUES ($1, $2, $3, $4)
ON CONFLICT (epoch_id) DO UPDATE SET
  start_time = EXCLUDED.start_time,
  daily_emission = EXCLUDED.daily_emission,
  dao_emission = COALESCE(EXCLUDED.dao_emission, epochs.dao_emission);

-- name: GetEpoch :one
SELECT epoch_id, start_time, daily_emission, dao_emission FROM epochs WHERE epoch_id = $1;

-- name: ListEpochs :many
SELECT epoch_id, start_time, daily_emission, dao_emission FROM epochs
ORDER BY epoch_id DESC LIMIT $1 OFFSET $2;

-- name: GetLatestEpoch :one
SELECT epoch_id, start_time, daily_emission, dao_emission FROM epochs
ORDER BY epoch_id DESC LIMIT 1;

-- name: UpdateEpochDAO :exec
UPDATE epochs SET dao_emission = $2 WHERE epoch_id = $1;

-- name: InsertRecipientAWPDistribution :exec
INSERT INTO recipient_awp_distributions (epoch_id, recipient, awp_amount) VALUES ($1, $2, $3)
ON CONFLICT (epoch_id, recipient) DO NOTHING;

-- name: GetRecipientEarnings :many
SELECT epoch_id, recipient, awp_amount FROM recipient_awp_distributions
WHERE recipient = $1 ORDER BY epoch_id DESC LIMIT $2 OFFSET $3;

-- name: GetSubnetEarningsByID :many
SELECT r.epoch_id, r.recipient, r.awp_amount
FROM recipient_awp_distributions r
JOIN subnets s ON r.recipient = s.subnet_contract
WHERE s.subnet_id = $1
ORDER BY r.epoch_id DESC LIMIT $2 OFFSET $3;

-- name: GetEpochDistributions :many
SELECT epoch_id, recipient, awp_amount FROM recipient_awp_distributions
WHERE epoch_id = $1 ORDER BY recipient;

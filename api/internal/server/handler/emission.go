package handler

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"

	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// emissionProjection is a projected emission entry (strings used to avoid float64 precision loss)
type emissionProjection struct {
	Days           int    `json:"days"`
	TotalEmission  string `json:"totalEmission"`
	FinalDailyRate string `json:"finalDailyRate"`
}

// GetCurrentEmission retrieves the current emission data from the Redis cache
func (h *Handler) GetCurrentEmission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	val, err := h.rdb.Get(ctx, "emission_current").Result()
	if err != nil {
		if err == redis.Nil {
			h.writeError(w, http.StatusServiceUnavailable, "emission data not yet available (keeper may not be running)")
			return
		}
		h.logger.Error("failed to read Redis emission_current", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get current emission data")
		return
	}

	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		h.logger.Error("failed to parse emission_current JSON", "error", err)
		h.writeError(w, http.StatusInternalServerError, "emission data format error")
		return
	}

	h.writeJSON(w, http.StatusOK, data)
}

var (
	decayFactor    = big.NewInt(996844)
	decayPrecision = big.NewInt(1000000)
	// Initial daily emission: 15,800,000 * 1e18
	initialDailyEmission, _ = new(big.Int).SetString("15800000000000000000000000", 10)
)

// GetEmissionSchedule computes emission projections for the next 30/90/365 days
// using big.Int integer arithmetic that exactly mirrors the contract logic
func (h *Handler) GetEmissionSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Use the latest epoch's emission as the baseline
	currentDaily := new(big.Int).Set(initialDailyEmission)
	latestEpoch, err := h.queries.GetLatestEpoch(ctx, h.cfg.ChainID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Error("failed to get latest epoch", "error", err)
			h.writeError(w, http.StatusInternalServerError, "failed to get emission data")
			return
		}
		// No epoch data; use initial emission
	} else if latestEpoch.DailyEmission.Valid {
		currentDaily = latestEpoch.DailyEmission.Int
	}

	// Compute emission projections for each time period
	periods := []int{30, 90, 365}
	projections := make([]emissionProjection, 0, len(periods))

	for _, days := range periods {
		total := new(big.Int)
		daily := new(big.Int).Set(currentDaily)

		for d := 0; d < days; d++ {
			total.Add(total, daily)
			// daily = daily * 996844 / 1000000 (integer decay matching the contract)
			daily.Mul(daily, decayFactor)
			daily.Div(daily, decayPrecision)
		}

		projections = append(projections, emissionProjection{
			Days:           days,
			TotalEmission:  total.String(),
			FinalDailyRate: daily.String(),
		})
	}

	h.writeJSON(w, http.StatusOK, map[string]any{
		"currentDailyEmission": currentDaily.String(),
		"projections":          projections,
	})
}

// ListEpochs returns a paginated list of epochs
func (h *Handler) ListEpochs(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePageParams(r)

	epochs, err := h.queries.ListEpochs(r.Context(), gen.ListEpochsParams{
		ChainID: h.cfg.ChainID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		h.logger.Error("failed to list epochs", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list epochs")
		return
	}

	h.writeJSON(w, http.StatusOK, epochs)
}

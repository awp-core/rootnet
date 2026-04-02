package handler

import (
	"math/big"
	"net/http"
)

// emissionProjection is a projected emission entry (strings used to avoid float64 precision loss)
type emissionProjection struct {
	Days           int    `json:"days"`
	TotalEmission  string `json:"totalEmission"`
	FinalDailyRate string `json:"finalDailyRate"`
}

var (
	decayFactor    = big.NewInt(996844)
	decayPrecision = big.NewInt(1000000)
	// Initial daily emission: 15,800,000 * 1e18
	initialDailyEmission, _ = new(big.Int).SetString("15800000000000000000000000", 10)
)

// GetCurrentEmission retrieves the current emission data from the Redis cache
func (h *Handler) GetCurrentEmission(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	data, err := h.svcGetCurrentEmission(r.Context(), chainID)
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, data)
}

// GetEmissionSchedule computes emission projections for the next 30/90/365 days
// using big.Int integer arithmetic that exactly mirrors the contract logic
func (h *Handler) GetEmissionSchedule(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	result, err := h.svcGetEmissionSchedule(r.Context(), chainID)
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// ListEpochs returns a paginated list of epochs
func (h *Handler) ListEpochs(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	limit, offset := h.parsePageParams(r)
	result, err := h.svcListEpochs(r.Context(), chainID, int32(limit), int32(offset))
	if err != nil {
		h.writeSvcError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

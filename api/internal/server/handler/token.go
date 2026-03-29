package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// GetAWPInfo retrieves AWP token info from the Redis cache
func (h *Handler) GetAWPInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	val, err := h.rdb.Get(ctx, fmt.Sprintf("awp_info:%d", h.cfg.ChainID)).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss; return empty object
			h.writeJSON(w, http.StatusOK, map[string]any{})
			return
		}
		h.logger.Error("failed to read Redis awp_info", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get AWP info")
		return
	}

	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		h.logger.Error("failed to parse awp_info JSON", "error", err)
		h.writeError(w, http.StatusInternalServerError, "AWP data format error")
		return
	}

	h.writeJSON(w, http.StatusOK, data)
}

// GetAlphaInfo retrieves subnet Alpha token info from the database
func (h *Handler) GetAlphaInfo(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	subnet, err := h.queries.GetSubnet(r.Context(), subnetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "subnet not found")
			return
		}
		h.logger.Error("failed to get subnet info", "error", err, "subnetId", subnetID)
		h.writeError(w, http.StatusInternalServerError, "failed to get subnet info")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]any{
		"subnetId":   subnet.SubnetID,
		"name":       subnet.Name,
		"symbol":     subnet.Symbol,
		"alphaToken": subnet.AlphaToken,
	})
}

// GetAlphaPrice retrieves the Alpha token price from the Redis cache
func (h *Handler) GetAlphaPrice(w http.ResponseWriter, r *http.Request) {
	subnetID, err := parseSubnetID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	key := fmt.Sprintf("alpha_price:%d", subnetID)

	val, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss; return empty object
			h.writeJSON(w, http.StatusOK, map[string]any{})
			return
		}
		h.logger.Error("failed to read Redis alpha_price", "error", err, "key", key)
		h.writeError(w, http.StatusInternalServerError, "failed to get Alpha price")
		return
	}

	var data any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		h.logger.Error("failed to parse alpha_price JSON", "error", err)
		h.writeError(w, http.StatusInternalServerError, "price data format error")
		return
	}

	h.writeJSON(w, http.StatusOK, data)
}

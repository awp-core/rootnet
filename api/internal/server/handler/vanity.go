package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"errors"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/jackc/pgx/v5"
)

// VanityHandler handles vanity address salt computation requests
type VanityHandler struct {
	factoryAddr  string
	initCodeHash string
	rule         chain.VanityRule
	timeout      time.Duration
	sem          chan struct{} // concurrency limiter semaphore
	queries      *gen.Queries  // DB queries (salt pool)
	limiter      *ratelimit.Limiter
	logger       *slog.Logger
}

// NewVanityHandler creates a VanityHandler
func NewVanityHandler(factoryAddr string, initCodeHash string, rule chain.VanityRule, queries *gen.Queries, limiter *ratelimit.Limiter, logger *slog.Logger) *VanityHandler {
	return &VanityHandler{
		factoryAddr:  factoryAddr,
		initCodeHash: initCodeHash,
		rule:         rule,
		timeout:      120 * time.Second,
		sem:          make(chan struct{}, 2),
		queries:      queries,
		limiter:      limiter,
		logger:       logger,
	}
}

func (vh *VanityHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (vh *VanityHandler) writeError(w http.ResponseWriter, status int, msg string) {
	vh.writeJSON(w, status, map[string]string{"error": msg})
}

type computeSaltResponse struct {
	Salt    string `json:"salt"`
	Address string `json:"address"`
	Source  string `json:"source"`  // "pool" or "mined"
	Elapsed string `json:"elapsed"`
}

// ComputeSalt POST /api/vanity/compute-salt — get a salt matching the vanityRule
// Tries DB salt pool first; falls back to cast create2 live mining if pool is empty
func (vh *VanityHandler) ComputeSalt(w http.ResponseWriter, r *http.Request) {
	// Rate limit (hot-updatable via Redis: HSET ratelimit:config compute_salt "20:3600")
	ip := ratelimit.GetClientIP(r)
	exceeded, _ := vh.limiter.CheckIP(r.Context(), "compute_salt", ip)
	if exceeded {
		vh.writeError(w, http.StatusTooManyRequests, vh.limiter.FormatError(r.Context(), "compute_salt"))
		return
	}
	defer vh.limiter.RecordSuccess("compute_salt", ip)

	start := time.Now()
	ctx := r.Context()

	// 1. Claim a salt from DB pool (atomic UPDATE+RETURNING with FOR UPDATE SKIP LOCKED)
	if vh.queries != nil {
		row, err := vh.queries.ClaimRandomSalt(ctx)
		if err == nil {
			vh.writeJSON(w, http.StatusOK, computeSaltResponse{
				Salt:    row.Salt,
				Address: row.Address,
				Source:  "pool",
				Elapsed: time.Since(start).Round(time.Millisecond).String(),
			})
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			vh.logger.Info("salt pool empty, falling back to cast create2")
		} else {
			vh.logger.Error("get salt from pool failed, falling back to mining", "error", err)
		}
	}

	// 2. Fallback: live mining via cast create2
	select {
	case vh.sem <- struct{}{}:
		defer func() { <-vh.sem }()
	default:
		vh.writeError(w, http.StatusTooManyRequests, "salt pool empty and too many concurrent mining requests")
		return
	}

	mineCtx, cancel := context.WithTimeout(ctx, vh.timeout)
	defer cancel()

	result, err := chain.FindVanitySalt(mineCtx, vh.factoryAddr, vh.initCodeHash, vh.rule)
	elapsed := time.Since(start)

	if err != nil {
		if mineCtx.Err() != nil {
			vh.writeError(w, http.StatusRequestTimeout, "salt pool empty and mining timed out")
			return
		}
		vh.logger.Error("vanity salt mining failed", "error", err)
		vh.writeError(w, http.StatusInternalServerError, "salt mining failed")
		return
	}

	vh.writeJSON(w, http.StatusOK, computeSaltResponse{
		Salt:    result.Salt,
		Address: result.Address,
		Source:  "mined",
		Elapsed: elapsed.Round(time.Millisecond).String(),
	})
}

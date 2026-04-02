package handler

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/cortexia/rootnet/api/internal/db/gen"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ── Request/Response types ──

type uploadSaltRequest struct {
	Salts []saltEntry `json:"salts"`
}

type saltEntry struct {
	Salt    string `json:"salt"`    // bytes32 hex (0x prefixed)
	Address string `json:"address"` // predicted Alpha token address (0x prefixed)
}

type saltCountResponse struct {
	Available int64 `json:"available"`
}

type miningParamsResponse struct {
	FactoryAddress string `json:"factoryAddress"` // AlphaTokenFactory contract address
	InitCodeHash   string `json:"initCodeHash"`   // keccak256(AlphaToken.creationCode)
	VanityRule     string `json:"vanityRule"`      // uint64 hex-encoded vanity rule
}

// GetMiningParams GET /api/vanity/mining-params — returns params needed for offline salt mining
func (h *Handler) GetMiningParams(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, miningParamsResponse{
		FactoryAddress: h.cfg.AlphaFactoryAddress,
		InitCodeHash:   h.cfg.AlphaInitCodeHash,
		VanityRule:     h.cfg.VanityRule,
	})
}

// UploadSalts POST /api/vanity/upload-salts — batch upload pre-mined salts
// Each salt is verified: 1) CREATE2 address correctness 2) vanityRule compliance
func (h *Handler) UploadSalts(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB — enforce before rate limit

	ip := ratelimit.GetClientIP(r)
	if exceeded, err := h.limiter.CheckAndIncrement(r.Context(), "upload_salts", ip); exceeded {
		h.writeError(w, http.StatusTooManyRequests, h.limiter.FormatError(r.Context(), "upload_salts"))
		return
	} else if err != nil {
		h.logger.Error("upload_salts rate limit error", "error", err)
	}

	var req uploadSaltRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Salts) == 0 {
		h.writeError(w, http.StatusBadRequest, "empty salts array")
		return
	}
	if len(req.Salts) > 1000 {
		h.writeError(w, http.StatusBadRequest, "max 1000 salts per request")
		return
	}

	if h.cfg.AlphaFactoryAddress == "" || h.cfg.AlphaInitCodeHash == "" {
		h.writeError(w, http.StatusInternalServerError, "vanity mining not configured")
		return
	}

	// Decode vanityRule for validation
	var vanityRule chain.VanityRule
	if h.cfg.VanityRule != "" {
		vanityRule, _ = chain.DecodeVanityRule(h.cfg.VanityRule)
	}

	ctx := r.Context()
	chainID := h.resolveChainID(r)
	inserted := 0
	rejected := 0
	for _, s := range req.Salts {
		if len(s.Salt) != 66 || len(s.Address) != 42 {
			rejected++
			continue
		}

		// Verify CREATE2 address correctness (using go-ethereum crypto.CreateAddress2)
		expectedAddr := computeCreate2Address(h.cfg.AlphaFactoryAddress, s.Salt, h.cfg.AlphaInitCodeHash)
		if !strings.EqualFold(expectedAddr, s.Address) {
			rejected++
			continue
		}

		// Verify vanityRule compliance
		if !vanityRule.IsEmpty() {
			addrHex := strings.TrimPrefix(strings.ToLower(expectedAddr), "0x")
			if !vanityRule.ValidateAddress(addrHex) {
				rejected++
				continue
			}
		}

		if err := h.queries.InsertVanitySalt(ctx, gen.InsertVanitySaltParams{
			ChainID: chainID,
			Salt:    s.Salt,
			Address: strings.ToLower(s.Address),
		}); err != nil {
			h.logger.Error("insert vanity salt failed", "error", err, "salt", s.Salt)
			continue
		}
		inserted++
	}

	h.writeJSON(w, http.StatusOK, map[string]int{"inserted": inserted, "rejected": rejected})
}

// ListAvailableSalts GET /api/vanity/salts — list available salts
func (h *Handler) ListAvailableSalts(w http.ResponseWriter, r *http.Request) {
	limit, _ := h.parsePageParams(r)
	chainID := h.resolveChainID(r)

	salts, err := h.queries.ListAvailableSalts(r.Context(), gen.ListAvailableSaltsParams{
		ChainID: chainID,
		Limit:   int32(limit),
	})
	if err != nil {
		h.logger.Error("list available salts failed", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list salts")
		return
	}

	entries := make([]saltEntry, len(salts))
	for i, s := range salts {
		entries[i] = saltEntry{Salt: s.Salt, Address: s.Address}
	}
	h.writeJSON(w, http.StatusOK, entries)
}

// CountAvailableSalts GET /api/vanity/salts/count — count available salts
func (h *Handler) CountAvailableSalts(w http.ResponseWriter, r *http.Request) {
	chainID := h.resolveChainID(r)
	count, err := h.queries.CountAvailableSalts(r.Context(), chainID)
	if err != nil {
		h.logger.Error("count available salts failed", "error", err)
		h.writeError(w, http.StatusInternalServerError, "failed to count salts")
		return
	}
	h.writeJSON(w, http.StatusOK, saltCountResponse{Available: count})
}

// ── Helpers ──

// computeCreate2Address computes CREATE2 address using go-ethereum library
func computeCreate2Address(factoryHex, saltHex, initCodeHashHex string) string {
	factory := common.HexToAddress(factoryHex)

	saltBytes, err := hex.DecodeString(strings.TrimPrefix(saltHex, "0x"))
	if err != nil || len(saltBytes) != 32 {
		return ""
	}
	var salt [32]byte
	copy(salt[:], saltBytes)

	hashBytes, err := hex.DecodeString(strings.TrimPrefix(initCodeHashHex, "0x"))
	if err != nil || len(hashBytes) != 32 {
		return ""
	}

	addr := crypto.CreateAddress2(factory, salt, hashBytes)
	return strings.ToLower(addr.Hex())
}

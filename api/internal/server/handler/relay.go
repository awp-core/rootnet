package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
)

// RelayHandler handles gasless relay transaction requests
type RelayHandler struct {
	relayer *chain.Relayer
	rdb     *redis.Client
	logger  *slog.Logger

	// IP rate limit params
	rateLimit  int           // max requests per window
	rateWindow time.Duration // window duration
}

// NewRelayHandler creates a RelayHandler
func NewRelayHandler(relayer *chain.Relayer, rdb *redis.Client, logger *slog.Logger) *RelayHandler {
	return &RelayHandler{
		relayer:    relayer,
		rdb:        rdb,
		logger:     logger,
		rateLimit:  5,
		rateWindow: 4 * time.Hour,
	}
}

// checkRateLimit checks IP rate limit, returns true if exceeded
func (rh *RelayHandler) checkRateLimit(r *http.Request) (bool, error) {
	// chi middleware.RealIP sets RemoteAddr to real IP (may include port)
	ip := r.RemoteAddr
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}

	key := "relay_ratelimit:" + ip
	ctx := r.Context()

	// Atomic Lua script: INCR + conditional EXPIRE, avoids TOCTOU race
	luaScript := redis.NewScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return count
	`)
	count, err := luaScript.Run(ctx, rh.rdb, []string{key}, int(rh.rateWindow.Seconds())).Int64()
	if err != nil {
		return false, fmt.Errorf("redis rate limit: %w", err)
	}

	return count > int64(rh.rateLimit), nil
}

func (rh *RelayHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (rh *RelayHandler) writeError(w http.ResponseWriter, status int, msg string) {
	rh.writeJSON(w, status, map[string]string{"error": msg})
}

// parseSignature splits a 0x-prefixed 65-byte hex signature into v, r, s
func parseSignature(sigHex string) (v uint8, r [32]byte, s [32]byte, err error) {
	sigHex = strings.TrimPrefix(sigHex, "0x")
	sig, err := hex.DecodeString(sigHex)
	if err != nil || len(sig) != 65 {
		return 0, r, s, fmt.Errorf("invalid signature: expected 65 bytes hex")
	}
	copy(r[:], sig[0:32])
	copy(s[:], sig[32:64])
	v = sig[64]
	// EIP-155: normalize v from 0/1 to 27/28
	if v < 27 {
		v += 27
	}
	return v, r, s, nil
}

// ── Request types ──

type relayRegisterRequest struct {
	User      string `json:"user"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayBindRequest struct {
	Agent     string `json:"agent"`
	Principal string `json:"principal"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayRegisterSubnetRequest struct {
	User              string `json:"user"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	MetadataURI       string `json:"metadataURI"`
	SubnetManager     string `json:"subnetManager"`     // "0x0...0" or "" = auto-deploy SubnetManager
	CoordinatorURL    string `json:"coordinatorURL"`
	Salt              string `json:"salt"`               // bytes32 hex, "0x00...00" = use subnetId
	MinStake          string `json:"minStake"`           // minimum stake wei string (0 = no minimum)
	Deadline          uint64 `json:"deadline"`
	PermitSignature   string `json:"permitSignature"`    // ERC-2612 permit signature (AWP approval)
	RegisterSignature string `json:"registerSignature"`  // EIP-712 registerSubnet signature
}

type relayResponse struct {
	TxHash string `json:"txHash"`
}

// RelayRegister POST /api/relay/register — relay registerFor transaction
func (rh *RelayHandler) RelayRegister(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req relayRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !common.IsHexAddress(req.User) {
		rh.writeError(w, http.StatusBadRequest, "invalid user address")
		return
	}
	if req.Deadline == 0 || int64(req.Deadline) <= time.Now().Unix() {
		rh.writeError(w, http.StatusBadRequest, "deadline is missing or expired")
		return
	}
	if req.Signature == "" {
		rh.writeError(w, http.StatusBadRequest, "missing signature")
		return
	}

	v, rs, ss, err := parseSignature(req.Signature)
	if err != nil {
		rh.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Rate limit check (after validation to avoid invalid requests consuming quota)
	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, "rate limit exceeded: max 5 requests per 4 hours")
		return
	}

	user := common.HexToAddress(req.User)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := rh.relayer.RelayRegister(r.Context(), user, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay registerFor failed", "error", err, "user", req.User)
		rh.writeError(w, http.StatusInternalServerError, "relay transaction failed")
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayBind POST /api/relay/bind — relay bindFor transaction
func (rh *RelayHandler) RelayBind(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req relayBindRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !common.IsHexAddress(req.Agent) {
		rh.writeError(w, http.StatusBadRequest, "invalid agent address")
		return
	}
	if !common.IsHexAddress(req.Principal) {
		rh.writeError(w, http.StatusBadRequest, "invalid principal address")
		return
	}
	if req.Deadline == 0 || int64(req.Deadline) <= time.Now().Unix() {
		rh.writeError(w, http.StatusBadRequest, "deadline is missing or expired")
		return
	}
	if req.Signature == "" {
		rh.writeError(w, http.StatusBadRequest, "missing signature")
		return
	}

	v, rs, ss, err := parseSignature(req.Signature)
	if err != nil {
		rh.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Rate limit check (after validation)
	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, "rate limit exceeded: max 5 requests per 4 hours")
		return
	}

	agent := common.HexToAddress(req.Agent)
	principal := common.HexToAddress(req.Principal)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := rh.relayer.RelayBind(r.Context(), agent, principal, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay bindFor failed", "error", err, "agent", req.Agent, "principal", req.Principal)
		rh.writeError(w, http.StatusInternalServerError, "relay transaction failed")
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayRegisterSubnet POST /api/relay/register-subnet — relay registerSubnetFor transaction
func (rh *RelayHandler) RelayRegisterSubnet(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)

	var req relayRegisterSubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !common.IsHexAddress(req.User) {
		rh.writeError(w, http.StatusBadRequest, "invalid user address")
		return
	}
	if req.Name == "" || req.Symbol == "" {
		rh.writeError(w, http.StatusBadRequest, "name and symbol are required")
		return
	}
	if req.Deadline == 0 || int64(req.Deadline) <= time.Now().Unix() {
		rh.writeError(w, http.StatusBadRequest, "deadline is missing or expired")
		return
	}
	if req.PermitSignature == "" || req.RegisterSignature == "" {
		rh.writeError(w, http.StatusBadRequest, "both permitSignature and registerSignature are required")
		return
	}

	permitV, permitR, permitS, err := parseSignature(req.PermitSignature)
	if err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid permitSignature: "+err.Error())
		return
	}
	registerV, registerR, registerS, err := parseSignature(req.RegisterSignature)
	if err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid registerSignature: "+err.Error())
		return
	}

	// Rate limit
	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, "rate limit exceeded: max 5 requests per 4 hours")
		return
	}

	// Build SubnetParams
	var salt [32]byte
	if req.Salt != "" {
		saltHex := strings.TrimPrefix(req.Salt, "0x")
		if b, decErr := hex.DecodeString(saltHex); decErr == nil && len(b) == 32 {
			copy(salt[:], b)
		}
	}

	var subnetMgr common.Address
	if req.SubnetManager != "" && common.IsHexAddress(req.SubnetManager) {
		subnetMgr = common.HexToAddress(req.SubnetManager)
	}

	minStake, _ := new(big.Int).SetString(req.MinStake, 10)
	if minStake == nil {
		minStake = big.NewInt(0)
	}

	params := bindings.IRootNetSubnetParams{
		Name:           req.Name,
		Symbol:         req.Symbol,
		MetadataURI:    req.MetadataURI,
		SubnetManager:  subnetMgr,
		CoordinatorURL: req.CoordinatorURL,
		Salt:           salt,
		MinStake:       minStake,
	}

	user := common.HexToAddress(req.User)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := rh.relayer.RelayRegisterSubnet(
		r.Context(), user, params, deadline,
		permitV, permitR, permitS,
		registerV, registerR, registerS,
	)
	if err != nil {
		rh.logger.Error("relay registerSubnetFor failed", "error", err, "user", req.User, "name", req.Name)
		rh.writeError(w, http.StatusInternalServerError, "relay transaction failed")
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

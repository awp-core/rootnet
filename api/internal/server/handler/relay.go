package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/cortexia/rootnet/api/internal/chain"
	"github.com/cortexia/rootnet/api/internal/chain/bindings"
	"github.com/cortexia/rootnet/api/internal/ratelimit"
	"github.com/ethereum/go-ethereum/common"
)

// revertErrors maps 4-byte Solidity error selectors to user-friendly messages
var revertErrors = map[string]string{
	"0x8baa579f": "invalid signature",
	"0x3a81d6fc": "user already registered",
	"0xaba47339": "user not registered",
	"0x682a9065": "agent already bound",
	"0x179435b3": "agent not bound",
	"0xdf4cc36d": "signature expired",
	"0x7991d6f7": "subnet manager address required (auto-deploy not available)",
	"0x6a6a9712": "invalid subnet params (name 1-64 bytes, symbol 1-16 bytes)",
	"0x00476ad8": "subnet not found",
	"0x2783bd34": "not subnet owner",
	"0x6ce3ac0e": "invalid subnet status for this operation",
	"0x72b5dd4b": "subnet immunity period still active",
	"0xab59c60f": "max active subnets reached",
	"0x0f38cabd": "allocation below subnet minimum stake",
	"0x9e87fac8": "contract is paused",
	// OpenZeppelin ECDSA errors
	"0xf645eedf": "invalid signature",
	"0xfce698f7": "invalid signature length",
	"0xd78bce0c": "invalid signature S value",
	// ERC20 errors
	"0xfb8f41b2": "insufficient AWP allowance",
	"0xe450d38c": "insufficient AWP balance",
	// AccessManager errors
	"0xf48904b2": "cannot bind to self",
	"0x05d14037": "not owner or manager",
	"0x27f5ce6b": "unknown address",
	"0x5e03d55f": "cannot remove self",
	// StakeNFT errors
	"0x6855a802": "lock not expired",
	"0x2bff29a6": "position expired",
	"0xd247d121": "insufficient unallocated stake",
	// Staking errors
	"0xc18316bf": "subnet not active",
	"0x0baf7432": "invalid allocation",
	"0x78838d16": "subnet ID cannot be zero",
	// Subnet lifecycle errors
	"0x84e3b93f": "immunity period not expired",
	"0x0dc149f0": "already initialized",
	// AlphaToken / Factory errors
	"0xc30436e9": "exceeds AWP max supply",
	"0x1c04203f": "exceeds mintable limit",
	"0x16a1ae75": "invalid vanity address (EIP-55 mismatch)",
	"0x00af5596": "not authorized (onlyRootNet)",
	"0x03119322": "LP pool already exists",
	// AWPEmission errors
	"0xc2cf00fc": "emission mining complete",
	"0xc7d141a8": "settlement in progress",
	"0x11631a24": "must be future epoch",
}

// decodeRelayError extracts a user-friendly message from an on-chain revert error
func decodeRelayError(err error) string {
	s := err.Error()
	// go-ethereum wraps revert data as "execution reverted: 0x{selector}"
	if idx := strings.Index(s, "0x"); idx >= 0 && len(s) >= idx+10 {
		selector := s[idx : idx+10]
		if msg, ok := revertErrors[selector]; ok {
			return msg
		}
		// Unknown selector — return it so the caller can debug
		return "on-chain revert: " + selector
	}
	// Non-revert errors (RPC timeout, nonce issues, gas estimation, etc.)
	if strings.Contains(s, "nonce") {
		return "relay nonce conflict, please retry"
	}
	if strings.Contains(s, "insufficient funds") {
		return "relayer has insufficient BNB for gas"
	}
	if strings.Contains(s, "timeout") || strings.Contains(s, "deadline") {
		return "relay RPC timeout, please retry"
	}
	return "relay request failed, please try again"
}

// RelayHandler handles gasless relay transaction requests
type RelayHandler struct {
	relayer *chain.Relayer
	limiter *ratelimit.Limiter
	logger  *slog.Logger
}

// NewRelayHandler creates a RelayHandler
func NewRelayHandler(relayer *chain.Relayer, limiter *ratelimit.Limiter, logger *slog.Logger) *RelayHandler {
	return &RelayHandler{
		relayer: relayer,
		limiter: limiter,
		logger:  logger,
	}
}

// checkRateLimit checks relay IP rate limit (read-only pre-check), returns true if exceeded.
func (rh *RelayHandler) checkRateLimit(r *http.Request) (bool, error) {
	return rh.limiter.CheckIP(r.Context(), "relay", ratelimit.GetClientIP(r))
}

// recordSuccessfulRelay increments the relay rate limit counter after a successful transaction.
func (rh *RelayHandler) recordSuccessfulRelay(r *http.Request) {
	rh.limiter.RecordSuccess("relay", ratelimit.GetClientIP(r))
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
	SubnetManager     string `json:"subnetManager"`     // "0x0...0" or "" = auto-deploy SubnetManager
	Salt              string `json:"salt"`               // bytes32 hex, "0x00...00" = use subnetId
	MinStake          string `json:"minStake"`           // minimum stake wei string (0 = no minimum)
	SkillsURI         string `json:"skillsUri"`          // skills description URI
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
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	user := common.HexToAddress(req.User)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := rh.relayer.RelayRegister(r.Context(), user, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay registerFor failed", "error", err, "user", req.User)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.recordSuccessfulRelay(r)
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
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	agent := common.HexToAddress(req.Agent)
	principal := common.HexToAddress(req.Principal)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := rh.relayer.RelayBind(r.Context(), agent, principal, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay bindFor failed", "error", err, "agent", req.Agent, "principal", req.Principal)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.recordSuccessfulRelay(r)
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
	if len(req.Name) > 64 {
		rh.writeError(w, http.StatusBadRequest, "name exceeds 64 bytes")
		return
	}
	if len(req.Symbol) > 16 {
		rh.writeError(w, http.StatusBadRequest, "symbol exceeds 16 bytes")
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
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
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
		Name:          req.Name,
		Symbol:        req.Symbol,
		SubnetManager: subnetMgr,
		Salt:          salt,
		MinStake:      minStake,
		SkillsURI:     req.SkillsURI,
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
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.recordSuccessfulRelay(r)
	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

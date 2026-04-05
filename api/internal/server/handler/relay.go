package handler

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-chi/chi/v5"
)

// revertErrors maps 4-byte Solidity error selectors to user-friendly messages.
// Selectors verified against contract source via cast sig on 2026-04-02.
var revertErrors = map[string]string{
	// ── AWPRegistry errors ──
	"0xd92e233d": "zero address",
	"0xf48904b2": "self-bind not allowed",
	"0x8b906c97": "caller is not the deployer",
	"0x0dc149f0": "already initialized",
	"0xcad4da2a": "caller is not the Timelock",
	"0xef6d0f02": "caller is not the Guardian",
	"0x33b4f8bf": "invalid worknet name (1-64 bytes)",
	"0xa581811c": "invalid worknet symbol (1-16 bytes)",
	"0xfdc63ea7": "computed LP AWP amount is zero",
	"0xccfad018": "mint amount must be > 0",
	"0x158dc889": "immunity period too short (min 7 days)",
	"0x4ff2aff2": "worknet manager address required (auto-deploy not available)",
	"0x6baa309b": "invalid worknet status for this operation",
	"0x84e3b93f": "immunity period not expired",
	"0x30cd7471": "not the worknet owner",
	"0xdbbbe822": "price too low",
	"0x24fe1192": "price too high",
	"0xdf4cc36d": "signature deadline expired",
	"0x8baa579f": "invalid EIP-712 signature",
	"0xfb92bbfa": "max active worknets reached (10000)",
	"0x9146c723": "cycle detected in binding tree",
	"0xd93ed101": "binding chain too deep (max 256)",
	"0xa18bfb23": "name or symbol contains invalid characters (\" or \\)",
	"0xea8e4eb5": "not authorized (not staker or delegate)",
	"0x373d7529": "cannot revoke own delegation",

	// ── OZ ERC20 errors ──
	"0xe450d38c": "insufficient AWP balance",
	"0xfb8f41b2": "insufficient AWP allowance",

	// ── OZ ECDSA errors ──
	"0xf645eedf": "invalid ECDSA signature",
	"0xfce698f7": "invalid signature length",
	"0xd78bce0c": "invalid signature S value",

	// ── OZ common errors ──
	"0xd93c0665": "contract is paused",
	"0x8dfc202b": "contract is not paused",
	"0x3ee5aeb5": "reentrancy guard: reentrant call",
	"0xf92ee8a9": "contract already initialized",
	"0xd7e6bcf8": "contract not initializing",
	"0x5274afe7": "SafeERC20: token operation failed",

	// ── UUPS errors ──
	"0xe07c8dba": "unauthorized UUPS upgrade caller",
	"0xaa1d49a4": "unsupported proxiable UUID",

	// ── veAWP errors ──
	"0x2c5211c6": "invalid amount",
	"0xd4005715": "lock duration too short",
	"0x59dc379f": "not the token owner",
	"0x6855a802": "lock not expired",
	"0x2bff29a6": "position lock expired (cannot add to expired position)",
	"0xd247d121": "insufficient unallocated stake",
	"0x2995fef1": "nothing to update",
	"0x401640ae": "lock duration cannot be shortened",
	"0x481ffa6a": "lock must exceed current time",
	"0xd0bfc8d2": "not authorized (onlyAWPRegistry)",
	"0xa13bdd4f": "cannot rescue staked AWP token",

	// ── AWPAllocator errors ──
	"0xdf2d4774": "insufficient allocation",
	"0x1f2a2005": "amount cannot be zero",
	"0x54ada055": "amount exceeds uint128 max",
	"0xd76d9a3d": "worknet ID cannot be zero",
	"0x2fe6f8a8": "allocation would overflow uint128",
	"0xa741a045": "already set",

	// ── WorknetToken errors ──
	"0xc30436e9": "exceeds AWP max supply",
	"0x1c04203f": "exceeds worknet token mintable limit",
	"0x7bfa4b9f": "not admin",
	"0xf8d2906c": "not minter",
	"0x69b757d8": "minter paused",
	"0x815eb757": "minters locked",
	"0xf7a632f5": "invalid callback",
	"0xddefae28": "initial mint already called",

	// ── WorknetTokenFactory errors ──
	"0x16a1ae75": "invalid vanity address (EIP-55 mismatch)",

	// ── LPManager errors ──
	"0x03119322": "LP pool already exists",
	"0xf591b277": "no LP pool exists for this token",
	"0x5a916d92": "amount exceeds Permit2 uint160 limit",

	// ── AWPEmission errors ──
	"0x9c8d2cd2": "invalid recipient",
	"0xac4258ee": "epoch not ready for settlement",
	"0xc2cf00fc": "emission mining complete (MAX_SUPPLY reached)",
	"0xc7d141a8": "settlement in progress",
	"0x19a2a9bd": "zero weight not allowed",
	"0x3971ddf8": "too many recipients",
	"0x5ae6162f": "invalid decay factor",
	"0x96cb8a7f": "limit cannot be zero",
	"0x8b8a8a93": "genesis time not reached yet",
	"0x66128a97": "epoch duration cannot be zero",
	"0x11631a24": "must be future epoch",

	"0xb7c39e25": "invalid resume time (must be in the future)",

	// ── WorknetManager errors ──
	"0x334ee9a1": "invalid slippage tolerance (1-5000 bps)",
	"0xf4742923": "strategy is paused",
	"0x646cf558": "already claimed",
	"0x09bde339": "invalid Merkle proof",
	"0xb466ddbf": "Merkle root already set for this epoch",
	"0x8da6b984": "no Merkle root for this epoch",
	"0xb263ae73": "zero Merkle root",

	// ── WorknetManagerUni errors ──
	"0xae18210a": "not the pool manager",
	"0x8199f5f3": "slippage exceeded",

	// ── AWPDAO errors ──
	"0xdf957883": "no tokens provided",
	"0xa1f0d74b": "token already voted in this proposal",
	"0xf6fafba0": "lock expired",
	"0x467b1124": "NFT minted after proposal creation",
	"0x376ef12e": "use proposeWithTokens instead",
	"0xc10b42c7": "use castVoteWithParams instead",
	"0xbf5bbc9c": "invalid quorum percent",
	"0x2721b57b": "zero total voting power",

	// ── AWPWorkNet errors ──
	"0x44943622": "token does not exist",
}

// decodeRelayError extracts a user-friendly message from an on-chain revert error
func decodeRelayError(err error) string {
	// First, try to extract revert data from go-ethereum's rpc.DataError interface.
	// In go-ethereum v1.17+, EstimateGas returns jsonError where ErrorData()
	// contains the actual revert data (e.g. "0x8baa579f") while Error() only
	// has the message ("execution reverted").
	var dataErr rpc.DataError
	if errors.As(err, &dataErr) {
		if data, ok := dataErr.ErrorData().(string); ok && strings.HasPrefix(data, "0x") && len(data) >= 10 {
			selector := data[:10]
			if msg, ok := revertErrors[selector]; ok {
				return msg
			}
			return "on-chain revert: " + selector
		}
	}

	s := err.Error()

	// Fallback: try to extract 4-byte error selector from error string formats:
	// - "execution reverted: 0x{selector}{args...}"
	// - "error code 3: execution reverted, data: \"0x{selector}{args...}\""
	// - "VM Exception: revert 0x{selector}"
	for _, prefix := range []string{"data: \"0x", "data: 0x", "reverted: 0x", "revert 0x"} {
		if idx := strings.Index(s, prefix); idx >= 0 {
			hexStart := idx + len(prefix) - 2 // point to "0x"
			if len(s) >= hexStart+10 {
				selector := s[hexStart : hexStart+10]
				if msg, ok := revertErrors[selector]; ok {
					return msg
				}
				return "on-chain revert: " + selector
			}
		}
	}
	// Fallback: find a valid 4-byte selector (0x + exactly 8 hex chars)
	if idx := strings.Index(s, "0x"); idx >= 0 && len(s) >= idx+10 {
		selector := s[idx : idx+10]
		// Validate it's a proper hex selector
		isHex := true
		for _, c := range selector[2:] {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				isHex = false
				break
			}
		}
		if isHex {
			if msg, ok := revertErrors[selector]; ok {
				return msg
			}
			return "on-chain revert: " + selector
		}
	}
	// Non-revert errors — match common patterns without leaking internal details
	if strings.Contains(s, "nonce") {
		return "relay nonce conflict, please retry"
	}
	if strings.Contains(s, "insufficient funds") {
		return "relayer has insufficient BNB for gas"
	}
	if strings.Contains(s, "timeout") || strings.Contains(s, "context deadline") {
		return "relay RPC timeout, please retry"
	}
	if strings.Contains(s, "connection refused") || strings.Contains(s, "dial tcp") {
		return "relay RPC connection failed, please retry later"
	}
	if strings.Contains(s, "gas required exceeds") || strings.Contains(s, "gas limit") {
		return "transaction gas estimation failed (check parameters)"
	}
	if strings.Contains(s, "already known") {
		return "transaction already pending, please wait"
	}
	if strings.Contains(s, "replacement transaction") {
		return "relay nonce conflict (replacement tx), please retry"
	}
	if strings.Contains(s, "execution reverted") {
		return "on-chain execution reverted (unrecognized error selector)"
	}
	// Final fallback — include a sanitized hint from the error type
	return "relay failed: internal error, please try again"
}

// RelayHandler handles gasless relay transaction requests (multi-chain)
type RelayHandler struct {
	relayers map[int64]*chain.Relayer // chainId → relayer
	limiter  *ratelimit.Limiter
	logger   *slog.Logger
}

// NewRelayHandler creates a multi-chain RelayHandler
func NewRelayHandler(relayers map[int64]*chain.Relayer, limiter *ratelimit.Limiter, logger *slog.Logger) *RelayHandler {
	return &RelayHandler{
		relayers: relayers,
		limiter:  limiter,
		logger:   logger,
	}
}

// resolveRelayer returns the relayer for the given chainId, or the single relayer if only one exists
func (rh *RelayHandler) resolveRelayer(chainID int64) (*chain.Relayer, error) {
	// Single-chain mode: directly return the only relayer
	if len(rh.relayers) == 1 && chainID == 0 {
		for _, r := range rh.relayers {
			return r, nil
		}
	}
	r, ok := rh.relayers[chainID]
	if !ok {
		return nil, fmt.Errorf("unsupported chainId: %d", chainID)
	}
	return r, nil
}

// checkRateLimit atomically checks and increments the relay IP rate limit.
// Returns true if exceeded (counter NOT incremented when exceeded).
func (rh *RelayHandler) checkRateLimit(r *http.Request) (bool, error) {
	return rh.limiter.CheckAndIncrement(r.Context(), "relay", ratelimit.GetClientIP(r))
}


func (rh *RelayHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		rh.logger.Error("failed to encode JSON response", "error", err)
	}
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
	ChainID   int64  `json:"chainId"`
	User      string `json:"user"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relaySetRecipientRequest struct {
	ChainID   int64  `json:"chainId"`
	User      string `json:"user"`
	Recipient string `json:"recipient"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayBindRequest struct {
	ChainID   int64  `json:"chainId"`
	Agent     string `json:"agent"`
	Target    string `json:"target"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayAllocateRequest struct {
	ChainID   int64  `json:"chainId"`
	Staker    string `json:"staker"`
	Agent     string `json:"agent"`
	SubnetID  string `json:"worknetId"`
	Amount    string `json:"amount"` // wei string
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayDeallocateRequest struct {
	ChainID   int64  `json:"chainId"`
	Staker    string `json:"staker"`
	Agent     string `json:"agent"`
	SubnetID  string `json:"worknetId"`
	Amount    string `json:"amount"` // wei string
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayActivateSubnetRequest struct {
	ChainID   int64  `json:"chainId"`
	User      string `json:"user"`
	SubnetID  string `json:"worknetId"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayRegisterSubnetRequest struct {
	ChainID   int64  `json:"chainId"`
	User              string `json:"user"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	WorknetManager     string `json:"worknetManager"`    // "0x0...0" or "" = auto-deploy WorknetManager
	Salt              string `json:"salt"`              // bytes32 hex, "0x00...00" = use worknetId
	MinStake          string `json:"minStake"`          // minimum stake wei string (0 = no minimum)
	SkillsURI         string `json:"skillsUri"`         // skills description URI
	Deadline          uint64 `json:"deadline"`
	PermitSignature   string `json:"permitSignature"`   // ERC-2612 permit signature (AWP approval)
	RegisterSignature string `json:"registerSignature"` // EIP-712 registerSubnet signature
}

type relayResponse struct {
	TxHash string `json:"txHash"`
}

// RelayRegister POST /api/relay/register — gasless registration (calls setRecipientFor(user, user))
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

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	user := common.HexToAddress(req.User)
	deadline := new(big.Int).SetUint64(req.Deadline)

	// register = setRecipientFor(user, user, ...) — sets recipient to self
	txHash, err := relayer.RelaySetRecipient(r.Context(), user, user, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay register failed", "error", err, "user", req.User)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayBind POST /api/relay/bind — relay bindFor transaction (V2: agent binds to target)
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
	if !common.IsHexAddress(req.Target) {
		rh.writeError(w, http.StatusBadRequest, "invalid target address")
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

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	agent := common.HexToAddress(req.Agent)
	target := common.HexToAddress(req.Target)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := relayer.RelayBind(r.Context(), agent, target, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay bindFor failed", "error", err, "agent", req.Agent, "target", req.Target)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelaySetRecipient POST /api/relay/set-recipient — relay setRecipientFor transaction (V2)
func (rh *RelayHandler) RelaySetRecipient(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req relaySetRecipientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !common.IsHexAddress(req.User) {
		rh.writeError(w, http.StatusBadRequest, "invalid user address")
		return
	}
	if !common.IsHexAddress(req.Recipient) {
		rh.writeError(w, http.StatusBadRequest, "invalid recipient address")
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

	// Rate limit check
	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	user := common.HexToAddress(req.User)
	recipient := common.HexToAddress(req.Recipient)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := relayer.RelaySetRecipient(r.Context(), user, recipient, deadline, v, rs, ss)
	if err != nil {
		rh.logger.Error("relay setRecipientFor failed", "error", err, "user", req.User, "recipient", req.Recipient)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayAllocate POST /api/relay/allocate — relay allocateFor transaction
func (rh *RelayHandler) RelayAllocate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var req relayAllocateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !common.IsHexAddress(req.Staker) { rh.writeError(w, http.StatusBadRequest, "invalid staker address"); return }
	if !common.IsHexAddress(req.Agent) { rh.writeError(w, http.StatusBadRequest, "invalid agent address"); return }
	if req.SubnetID == "" { rh.writeError(w, http.StatusBadRequest, "worknetId is required"); return }
	if req.Amount == "" { rh.writeError(w, http.StatusBadRequest, "amount is required"); return }
	if req.Deadline == 0 || int64(req.Deadline) <= time.Now().Unix() { rh.writeError(w, http.StatusBadRequest, "deadline is missing or expired"); return }
	if req.Signature == "" { rh.writeError(w, http.StatusBadRequest, "missing signature"); return }

	v, rs, ss, err := parseSignature(req.Signature)
	if err != nil { rh.writeError(w, http.StatusBadRequest, err.Error()); return }

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil { rh.writeError(w, http.StatusInternalServerError, "rate limit check failed"); return }
	if exceeded { rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay")); return }

	worknetId, ok := new(big.Int).SetString(req.SubnetID, 10)
	if !ok || worknetId.Sign() <= 0 { rh.writeError(w, http.StatusBadRequest, "invalid worknetId"); return }

	amount, amtOk := new(big.Int).SetString(req.Amount, 10)
	if !amtOk || amount.Sign() <= 0 { rh.writeError(w, http.StatusBadRequest, "invalid amount"); return }

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	txHash, err := relayer.RelayAllocate(r.Context(),
		common.HexToAddress(req.Staker), common.HexToAddress(req.Agent),
		worknetId, amount,
		new(big.Int).SetUint64(req.Deadline), v, rs, ss)
	if err != nil {
		rh.logger.Error("relay allocateFor failed", "error", err, "staker", req.Staker)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}
	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayDeallocate POST /api/relay/deallocate — relay deallocateFor transaction
func (rh *RelayHandler) RelayDeallocate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var req relayDeallocateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !common.IsHexAddress(req.Staker) { rh.writeError(w, http.StatusBadRequest, "invalid staker address"); return }
	if !common.IsHexAddress(req.Agent) { rh.writeError(w, http.StatusBadRequest, "invalid agent address"); return }
	if req.SubnetID == "" { rh.writeError(w, http.StatusBadRequest, "worknetId is required"); return }
	if req.Amount == "" { rh.writeError(w, http.StatusBadRequest, "amount is required"); return }
	if req.Deadline == 0 || int64(req.Deadline) <= time.Now().Unix() { rh.writeError(w, http.StatusBadRequest, "deadline is missing or expired"); return }
	if req.Signature == "" { rh.writeError(w, http.StatusBadRequest, "missing signature"); return }

	v, rs, ss, err := parseSignature(req.Signature)
	if err != nil { rh.writeError(w, http.StatusBadRequest, err.Error()); return }

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil { rh.writeError(w, http.StatusInternalServerError, "rate limit check failed"); return }
	if exceeded { rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay")); return }

	worknetId, ok := new(big.Int).SetString(req.SubnetID, 10)
	if !ok || worknetId.Sign() <= 0 { rh.writeError(w, http.StatusBadRequest, "invalid worknetId"); return }

	amount, amtOk := new(big.Int).SetString(req.Amount, 10)
	if !amtOk || amount.Sign() <= 0 { rh.writeError(w, http.StatusBadRequest, "invalid amount"); return }

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	txHash, err := relayer.RelayDeallocate(r.Context(),
		common.HexToAddress(req.Staker), common.HexToAddress(req.Agent),
		worknetId, amount,
		new(big.Int).SetUint64(req.Deadline), v, rs, ss)
	if err != nil {
		rh.logger.Error("relay deallocateFor failed", "error", err, "staker", req.Staker)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}
	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayActivateSubnet POST /api/relay/activate-subnet — relay activateSubnetFor transaction
func (rh *RelayHandler) RelayActivateSubnet(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var req relayActivateSubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !common.IsHexAddress(req.User) { rh.writeError(w, http.StatusBadRequest, "invalid user address"); return }
	if req.SubnetID == "" { rh.writeError(w, http.StatusBadRequest, "worknetId is required"); return }
	if req.Deadline == 0 || int64(req.Deadline) <= time.Now().Unix() { rh.writeError(w, http.StatusBadRequest, "deadline is missing or expired"); return }
	if req.Signature == "" { rh.writeError(w, http.StatusBadRequest, "missing signature"); return }

	v, rs, ss, err := parseSignature(req.Signature)
	if err != nil { rh.writeError(w, http.StatusBadRequest, err.Error()); return }

	worknetId, ok := new(big.Int).SetString(req.SubnetID, 10)
	if !ok || worknetId.Sign() <= 0 {
		rh.writeError(w, http.StatusBadRequest, "invalid worknetId")
		return
	}

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil { rh.writeError(w, http.StatusInternalServerError, "rate limit check failed"); return }
	if exceeded { rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay")); return }

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	txHash, err := relayer.RelayActivateSubnet(r.Context(),
		common.HexToAddress(req.User), worknetId,
		new(big.Int).SetUint64(req.Deadline), v, rs, ss)
	if err != nil {
		rh.logger.Error("relay activateSubnetFor failed", "error", err, "user", req.User, "worknetId", req.SubnetID)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
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
	if len(req.Name) > 64 {
		rh.writeError(w, http.StatusBadRequest, "name exceeds 64 bytes")
		return
	}
	if len(req.Symbol) > 16 {
		rh.writeError(w, http.StatusBadRequest, "symbol exceeds 16 bytes")
		return
	}
	if len(req.SkillsURI) > 2048 {
		rh.writeError(w, http.StatusBadRequest, "skillsUri exceeds 2048 bytes")
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
	if req.WorknetManager != "" && common.IsHexAddress(req.WorknetManager) {
		subnetMgr = common.HexToAddress(req.WorknetManager)
	}

	minStake, _ := new(big.Int).SetString(req.MinStake, 10)
	if minStake == nil {
		minStake = big.NewInt(0)
	}

	params := bindings.IAWPRegistryWorknetParams{
		Name:          req.Name,
		Symbol:        req.Symbol,
		WorknetManager: subnetMgr,
		Salt:          salt,
		MinStake:      minStake,
		SkillsURI:     req.SkillsURI,
	}

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	user := common.HexToAddress(req.User)
	deadline := new(big.Int).SetUint64(req.Deadline)

	txHash, err := relayer.RelayRegisterSubnet(
		r.Context(), user, params, deadline,
		permitV, permitR, permitS,
		registerV, registerR, registerS,
	)
	if err != nil {
		rh.logger.Error("relay registerSubnetFor failed", "error", err, "user", req.User, "name", req.Name)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// GetRelayStatus GET /api/relay/status/{txHash} — query on-chain confirmation status of a relay transaction
func (rh *RelayHandler) GetRelayStatus(w http.ResponseWriter, r *http.Request) {
	txHash := chi.URLParam(r, "txHash")
	if txHash == "" || len(txHash) != 66 {
		rh.writeError(w, http.StatusBadRequest, "invalid txHash")
		return
	}

	// Tx status is in Redis (shared) — use any relayer
	var anyRelayer *chain.Relayer
	for _, rl := range rh.relayers {
		anyRelayer = rl
		break
	}
	if anyRelayer == nil {
		rh.writeError(w, http.StatusServiceUnavailable, "no relayer configured")
		return
	}
	status, err := anyRelayer.GetTxStatus(r.Context(), txHash)
	if err != nil {
		rh.writeJSON(w, http.StatusOK, map[string]string{"status": "unknown", "txHash": txHash})
		return
	}

	rh.writeJSON(w, http.StatusOK, status)
}

// ── Grant/Revoke Delegate + Unbind relay endpoints ──

type relayDelegateRequest struct {
	ChainID   int64  `json:"chainId"`
	User      string `json:"user"`
	Delegate  string `json:"delegate"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

type relayUnbindRequest struct {
	ChainID   int64  `json:"chainId"`
	User      string `json:"user"`
	Deadline  uint64 `json:"deadline"`
	Signature string `json:"signature"`
}

// RelayGrantDelegate POST /api/relay/grant-delegate
func (rh *RelayHandler) RelayGrantDelegate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req relayDelegateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !common.IsHexAddress(req.User) || !common.IsHexAddress(req.Delegate) {
		rh.writeError(w, http.StatusBadRequest, "invalid address")
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

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	txHash, err := relayer.RelayGrantDelegate(r.Context(), common.HexToAddress(req.User), common.HexToAddress(req.Delegate), new(big.Int).SetUint64(req.Deadline), v, rs, ss)
	if err != nil {
		rh.logger.Error("relay grantDelegateFor failed", "error", err, "user", req.User)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayRevokeDelegate POST /api/relay/revoke-delegate
func (rh *RelayHandler) RelayRevokeDelegate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req relayDelegateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !common.IsHexAddress(req.User) || !common.IsHexAddress(req.Delegate) {
		rh.writeError(w, http.StatusBadRequest, "invalid address")
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

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	txHash, err := relayer.RelayRevokeDelegate(r.Context(), common.HexToAddress(req.User), common.HexToAddress(req.Delegate), new(big.Int).SetUint64(req.Deadline), v, rs, ss)
	if err != nil {
		rh.logger.Error("relay revokeDelegateFor failed", "error", err, "user", req.User)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

// RelayUnbind POST /api/relay/unbind
func (rh *RelayHandler) RelayUnbind(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req relayUnbindRequest
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

	exceeded, rateLimitErr := rh.checkRateLimit(r)
	if rateLimitErr != nil {
		rh.writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if exceeded {
		rh.writeError(w, http.StatusTooManyRequests, rh.limiter.FormatError(r.Context(), "relay"))
		return
	}

	relayer, resolveErr := rh.resolveRelayer(req.ChainID)
	if resolveErr != nil {
		rh.writeError(w, http.StatusBadRequest, resolveErr.Error())
		return
	}

	txHash, err := relayer.RelayUnbind(r.Context(), common.HexToAddress(req.User), new(big.Int).SetUint64(req.Deadline), v, rs, ss)
	if err != nil {
		rh.logger.Error("relay unbindFor failed", "error", err, "user", req.User)
		rh.writeError(w, http.StatusBadRequest, decodeRelayError(err))
		return
	}

	rh.writeJSON(w, http.StatusOK, relayResponse{TxHash: txHash})
}

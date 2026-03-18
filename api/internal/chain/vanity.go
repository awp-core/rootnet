package chain

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/crypto/sha3"
)

// VanityResult is the result of a cast create2 salt mining operation
type VanityResult struct {
	Salt    string // bytes32 hex (0x-prefixed)
	Address string // EIP-55 checksummed address
}

// VanityRule decoded vanity rule (8 position constraints)
// Decoded from the factory uint64 vanityRule
type VanityRule struct {
	Prefix [4]uint8 // prefix[0..3], maps to address hex[0..3]
	Suffix [4]uint8 // suffix[0..3], maps to address hex[36..39]
}

// DecodeVanityRule decodes a vanity rule from uint64 hex string
// e.g. "0x1001FFFF0C0A0F0E"
func DecodeVanityRule(ruleHex string) (VanityRule, error) {
	ruleHex = strings.TrimPrefix(ruleHex, "0x")
	if len(ruleHex) > 16 {
		return VanityRule{}, fmt.Errorf("vanity rule too long: %s", ruleHex)
	}
	// Left-pad with zeros to 16 characters
	for len(ruleHex) < 16 {
		ruleHex = "0" + ruleHex
	}

	b, err := hex.DecodeString(ruleHex)
	if err != nil {
		return VanityRule{}, fmt.Errorf("invalid vanity rule hex: %w", err)
	}

	var vr VanityRule
	// bytes [7..4] = prefix[0..3], bytes [3..0] = suffix[0..3]
	vr.Prefix[0] = b[0]
	vr.Prefix[1] = b[1]
	vr.Prefix[2] = b[2]
	vr.Prefix[3] = b[3]
	vr.Suffix[0] = b[4]
	vr.Suffix[1] = b[5]
	vr.Suffix[2] = b[6]
	vr.Suffix[3] = b[7]
	return vr, nil
}

// IsEmpty returns true if all positions are wildcards (no constraint)
func (vr VanityRule) IsEmpty() bool {
	for _, v := range vr.Prefix {
		if v < 22 {
			return false
		}
	}
	for _, v := range vr.Suffix {
		if v < 22 {
			return false
		}
	}
	return true
}

// positionToChar converts a rule value to a cast create2 match character
// 0-9 → '0'-'9', 10-15 → 'a'-'f', 16-21 → 'A'-'F', >=22 → wildcard
func positionToChar(v uint8) (byte, bool) {
	if v <= 9 {
		return '0' + v, true
	} else if v <= 15 {
		return 'a' + (v - 10), true
	} else if v <= 21 {
		return 'A' + (v - 16), true
	}
	return 0, false // wildcard
}

// buildContiguousPattern extracts the longest contiguous non-wildcard prefix for cast create2
// Returns (leading contiguous string, whether gap constraints need post-validation)
func buildContiguousPattern(positions [4]uint8) (string, bool) {
	var pattern []byte
	hasGap := false
	for _, v := range positions {
		if c, ok := positionToChar(v); ok {
			if hasGap {
				// wildcard followed by constraint — gap exists
				break
			}
			pattern = append(pattern, c)
		} else {
			if len(pattern) > 0 {
				hasGap = true
			}
		}
	}
	// Check for gap constraints (non-wildcard after wildcard)
	needPostValidation := false
	wildcardSeen := false
	for _, v := range positions {
		if v >= 22 {
			wildcardSeen = true
		} else if wildcardSeen {
			needPostValidation = true
			break
		}
	}
	return string(pattern), needPostValidation
}

// buildContiguousSuffix extracts the longest contiguous non-wildcard suffix (right to left)
func buildContiguousSuffix(positions [4]uint8) (string, bool) {
	var pattern []byte
	for i := 3; i >= 0; i-- {
		if c, ok := positionToChar(positions[i]); ok {
			pattern = append([]byte{c}, pattern...)
		} else {
			break
		}
	}
	// Check for gaps
	needPostValidation := false
	wildcardSeen := false
	for i := 3; i >= 0; i-- {
		if positions[i] >= 22 {
			wildcardSeen = true
		} else if wildcardSeen {
			needPostValidation = true
			break
		}
	}
	return string(pattern), needPostValidation
}

// hasCaseSensitive returns true if the rule contains letter constraints (requires --case-sensitive)
func (vr VanityRule) hasCaseSensitive() bool {
	for _, v := range vr.Prefix {
		if v >= 10 && v <= 21 {
			return true
		}
	}
	for _, v := range vr.Suffix {
		if v >= 10 && v <= 21 {
			return true
		}
	}
	return false
}

// ValidateAddress verifies an address satisfies the full vanity rule (including gap positions and EIP-55)
// addrHex: 40-char lowercase hex (no 0x prefix)
func (vr VanityRule) ValidateAddress(addrHex string) bool {
	if len(addrHex) != 40 {
		return false
	}
	addrHex = strings.ToLower(addrHex)

	// Compute EIP-55 hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(addrHex))
	hash := hasher.Sum(nil)

	// Check 4 prefix positions
	for i := 0; i < 4; i++ {
		if !checkPosition(addrHex, hash, i, vr.Prefix[i]) {
			return false
		}
	}
	// Check 4 suffix positions
	for i := 0; i < 4; i++ {
		if !checkPosition(addrHex, hash, 36+i, vr.Suffix[i]) {
			return false
		}
	}
	return true
}

// checkPosition mirrors the on-chain _checkPosition logic exactly
func checkPosition(hex40 string, hash []byte, pos int, expected uint8) bool {
	if expected >= 22 {
		return true // wildcard
	}
	c := hex40[pos]

	if expected <= 9 {
		return c == ('0' + expected)
	} else if expected <= 15 {
		letter := expected - 10
		if c != ('a' + letter) {
			return false
		}
		nibble := getHashNibble(hash, pos)
		return nibble < 8 // EIP-55 must stay lowercase
	} else {
		// 16-21: uppercase
		letter := expected - 16
		if c != ('a' + letter) {
			return false
		}
		nibble := getHashNibble(hash, pos)
		return nibble >= 8 // EIP-55 must be uppercase
	}
}

func getHashNibble(hash []byte, pos int) uint8 {
	b := hash[pos/2]
	if pos%2 == 0 {
		return b >> 4
	}
	return b & 0x0f
}

// FindVanitySalt mines a CREATE2 salt matching the vanityRule using Foundry cast create2
// For rules with gap wildcards, automatically retries until a full match is found
func FindVanitySalt(
	ctx context.Context,
	factoryAddr string,
	initCodeHash string,
	rule VanityRule,
) (*VanityResult, error) {
	if rule.IsEmpty() {
		return nil, fmt.Errorf("vanity rule is all wildcards, no mining needed")
	}

	prefix, prefixNeedPost := buildContiguousPattern(rule.Prefix)
	suffix, suffixNeedPost := buildContiguousSuffix(rule.Suffix)
	needPostValidation := prefixNeedPost || suffixNeedPost

	if prefix == "" && suffix == "" {
		// All constraints are in gap positions — cannot use cast starts-with/ends-with
		// Requires brute-force search (extremely rare case)
		return nil, fmt.Errorf("vanity rule has no contiguous prefix or suffix, pattern too complex for cast create2")
	}

	args := []string{
		"create2",
		"--deployer", factoryAddr,
		"--init-code-hash", initCodeHash,
	}
	if prefix != "" {
		args = append(args, "--starts-with", prefix)
	}
	if suffix != "" {
		args = append(args, "--ends-with", suffix)
	}
	if rule.hasCaseSensitive() {
		args = append(args, "--case-sensitive")
	}

	// If gap constraints exist, retry until full match
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		cmd := exec.CommandContext(ctx, "cast", args...)
		out, err := cmd.Output()
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			if exitErr, ok := err.(*exec.ExitError); ok {
				return nil, fmt.Errorf("cast create2 failed: %s", string(exitErr.Stderr))
			}
			return nil, fmt.Errorf("cast create2: %w", err)
		}

		result, err := parseCastOutput(string(out))
		if err != nil {
			return nil, err
		}

		if !needPostValidation {
			return result, nil
		}

		// Post-validate full rule
		addrHex := strings.TrimPrefix(strings.ToLower(result.Address), "0x")
		if rule.ValidateAddress(addrHex) {
			return result, nil
		}
		// Gap position constraint not met, retry
	}
}

// parseCastOutput parses cast create2 output
func parseCastOutput(output string) (*VanityResult, error) {
	var result VanityResult
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Address: ") {
			result.Address = strings.TrimPrefix(line, "Address: ")
		} else if strings.HasPrefix(line, "Salt: ") {
			saltPart := strings.TrimPrefix(line, "Salt: ")
			if idx := strings.Index(saltPart, " "); idx != -1 {
				saltPart = saltPart[:idx]
			}
			result.Salt = saltPart
		}
	}
	if result.Address == "" || result.Salt == "" {
		return nil, fmt.Errorf("failed to parse cast create2 output:\n%s", output)
	}
	return &result, nil
}

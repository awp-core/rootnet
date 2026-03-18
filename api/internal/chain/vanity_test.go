package chain

import (
	"context"
	"strings"
	"testing"
	"time"
)

const (
	testFactory      = "0x4200000000000000000000000000000000000042"
	testInitCodeHash = "0xec7670c4271d07eeacbc9724f60ae43975fc32ddb548a3e2a85acbe1c9dcf901"
)

func TestDecodeVanityRule(t *testing.T) {
	// "A1????cafe" → 0x1001FFFF0C0A0F0E
	rule, err := DecodeVanityRule("0x1001FFFF0C0A0F0E")
	if err != nil {
		t.Fatal(err)
	}

	// prefix: A(16), 1(1), wildcard(0xFF), wildcard(0xFF)
	if rule.Prefix[0] != 16 || rule.Prefix[1] != 1 || rule.Prefix[2] != 0xFF || rule.Prefix[3] != 0xFF {
		t.Errorf("unexpected prefix: %v", rule.Prefix)
	}
	// suffix: c(12), a(10), f(15), e(14)
	if rule.Suffix[0] != 12 || rule.Suffix[1] != 10 || rule.Suffix[2] != 15 || rule.Suffix[3] != 14 {
		t.Errorf("unexpected suffix: %v", rule.Suffix)
	}
	if rule.IsEmpty() {
		t.Error("rule should not be empty")
	}
}

func TestDecodeVanityRule_AllWildcards(t *testing.T) {
	rule, err := DecodeVanityRule("0xFFFFFFFFFFFFFFFF")
	if err != nil {
		t.Fatal(err)
	}
	if !rule.IsEmpty() {
		t.Error("rule should be empty (all wildcards)")
	}
}

func TestFindVanitySalt_DigitPrefix(t *testing.T) {
	// Rule: prefix "00", no suffix constraint
	rule := VanityRule{
		Prefix: [4]uint8{0, 0, 0xFF, 0xFF}, // "00??"
		Suffix: [4]uint8{0xFF, 0xFF, 0xFF, 0xFF},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := FindVanitySalt(ctx, testFactory, testInitCodeHash, rule)
	if err != nil {
		t.Fatal(err)
	}

	addr := strings.ToLower(result.Address)
	if !strings.HasPrefix(addr, "0x00") {
		t.Errorf("expected prefix '00', got %s", result.Address)
	}
	t.Logf("salt=%s address=%s", result.Salt, result.Address)
}

func TestFindVanitySalt_CaseSensitiveSuffix(t *testing.T) {
	// Rule: suffix "cafe" (lowercase c=12, a=10, f=15, e=14)
	rule := VanityRule{
		Prefix: [4]uint8{0xFF, 0xFF, 0xFF, 0xFF},
		Suffix: [4]uint8{12, 10, 15, 14}, // lowercase "cafe"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := FindVanitySalt(ctx, testFactory, testInitCodeHash, rule)
	if err != nil {
		t.Fatal(err)
	}

	// Post-validation
	addrHex := strings.TrimPrefix(strings.ToLower(result.Address), "0x")
	if !rule.ValidateAddress(addrHex) {
		t.Errorf("address %s does not pass full vanity rule validation", result.Address)
	}
	t.Logf("salt=%s address=%s", result.Salt, result.Address)
}

func TestFindVanitySalt_FullRule(t *testing.T) {
	// "A1????cafe" — matches contract example
	rule, _ := DecodeVanityRule("0x1001FFFF0C0A0F0E")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	result, err := FindVanitySalt(ctx, testFactory, testInitCodeHash, rule)
	if err != nil {
		t.Fatal(err)
	}

	addrHex := strings.TrimPrefix(strings.ToLower(result.Address), "0x")
	if !rule.ValidateAddress(addrHex) {
		t.Errorf("address %s does not pass vanity rule validation", result.Address)
	}
	t.Logf("salt=%s address=%s", result.Salt, result.Address)
}

func TestValidateAddress(t *testing.T) {
	rule, _ := DecodeVanityRule("0x1001FFFF0C0A0F0E")

	// Zero address should not pass A1????cafe rule (first char is not a)
	if rule.ValidateAddress("0000000000000000000000000000000000000000") {
		t.Error("zero address should not pass A1????cafe rule")
	}

	// All-wildcard rule should accept any address
	emptyRule := VanityRule{
		Prefix: [4]uint8{0xFF, 0xFF, 0xFF, 0xFF},
		Suffix: [4]uint8{0xFF, 0xFF, 0xFF, 0xFF},
	}
	if !emptyRule.ValidateAddress("0000000000000000000000000000000000000000") {
		t.Error("empty rule should accept any address")
	}
}

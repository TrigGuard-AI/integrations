// Command verify_apply blocks terraform apply unless an OER v1 receipt verifies offline for the plan JSON.
// See docs/integrations/terraform.md and docs/specs/VERIFIER_SDK.md.
package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	oer "github.com/TrigGuard-AI/TrigGuard/tools/oer-verifier-go"
)

func main() {
	os.Exit(run())
}

func run() int {
	flagReceipt := flag.String("receipt", "", "OER wire string (or env TG_EXECUTION_RECEIPT)")
	flagSurface := flag.String("surface", "", "Surface id sid (or env TG_SURFACE)")
	flagAction := flag.String("action-file", "", "Path to UTF-8 Terraform plan JSON (or env TG_ACTION_FILE, default terraform-plan.json)")
	flagPub := flag.String("public-key-hex", "", "Ed25519 public key as 64 hex chars (or env TG_PUBLIC_KEY_HEX)")
	flagNow := flag.String("now-unix", "", "Unix seconds for exp check (or env TG_NOW_UNIX; default now)")
	flag.Parse()

	receipt := firstNonEmpty(*flagReceipt, os.Getenv("TG_EXECUTION_RECEIPT"))
	surface := firstNonEmpty(*flagSurface, os.Getenv("TG_SURFACE"))
	actionPath := firstNonEmpty(*flagAction, os.Getenv("TG_ACTION_FILE"))
	if actionPath == "" {
		actionPath = "terraform-plan.json"
	}
	pubHex := firstNonEmpty(*flagPub, os.Getenv("TG_PUBLIC_KEY_HEX"))
	nowStr := firstNonEmpty(*flagNow, os.Getenv("TG_NOW_UNIX"))

	if receipt == "" {
		fmt.Fprintln(os.Stderr, "verify_apply: missing receipt (TG_EXECUTION_RECEIPT or -receipt)")
		return 1
	}
	if surface == "" {
		fmt.Fprintln(os.Stderr, "verify_apply: missing surface (TG_SURFACE or -surface)")
		return 1
	}
	if pubHex == "" {
		fmt.Fprintln(os.Stderr, "verify_apply: missing public key (TG_PUBLIC_KEY_HEX or -public-key-hex)")
		return 1
	}

	actionPath = resolveActionFile(actionPath)
	actionJSON, err := os.ReadFile(actionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify_apply: read action file: %v\n", err)
		return 1
	}

	pub, err := parsePublicKeyHex(pubHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify_apply: %v\n", err)
		return 1
	}

	var nowUnix int64
	if nowStr == "" {
		nowUnix = time.Now().Unix()
	} else {
		n, err := strconv.ParseInt(nowStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "verify_apply: invalid now-unix: %v\n", err)
			return 1
		}
		nowUnix = n
	}

	if err := oer.Verify(receipt, actionJSON, surface, pub, nowUnix); err != nil {
		fmt.Fprintf(os.Stderr, "verify_apply: %v\n", err)
		return 1
	}
	return 0
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func resolveActionFile(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if ws := os.Getenv("GITHUB_WORKSPACE"); ws != "" {
		return filepath.Join(ws, p)
	}
	if cd := os.Getenv("TG_CALLER_DIR"); cd != "" {
		return filepath.Join(cd, p)
	}
	return filepath.Clean(p)
}

func parsePublicKeyHex(s string) (ed25519.PublicKey, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "0x")
	if len(s) != 64 {
		return nil, fmt.Errorf("public key must be 64 hex characters (32 bytes)")
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("public key hex: %w", err)
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("public key must decode to 32 bytes")
	}
	return ed25519.PublicKey(raw), nil
}

// gen-oer-fixtures writes conformance/oer-v1/fixtures/vectors.json with deterministic OER v1 vectors.
// Run from repo root: go run ./tools/oer-verifier-go/cmd/gen-oer-fixtures
package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	oer "github.com/TrigGuard-AI/TrigGuard/tools/oer-verifier-go"
)

// payload fields in lexicographic JSON key order (matches canonical expectations).
type payload struct {
	Ahsh string `json:"ahsh"`
	Ctx  string `json:"ctx"`
	Dcsn string `json:"dcsn"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
	Pid  string `json:"pid"`
	Rid  string `json:"rid"`
	Sid  string `json:"sid"`
	V    int    `json:"v"`
}

func marshalPayload(p payload) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

func wireReceipt(payloadBytes, sig []byte) string {
	p1 := base64.RawURLEncoding.EncodeToString(payloadBytes)
	p2 := base64.RawURLEncoding.EncodeToString(sig)
	return p1 + "." + p2
}

func sign(priv ed25519.PrivateKey, payloadBytes []byte) []byte {
	return ed25519.Sign(priv, payloadBytes)
}

func main() {
	root := findRepoRoot()
	outPath := filepath.Join(root, "conformance", "oer-v1", "fixtures", "vectors.json")

	seed := sha256.Sum256([]byte("trigguard-oer-conformance-v1-fixed-seed"))
	priv := ed25519.NewKeyFromSeed(seed[:32])
	pub := priv.Public().(ed25519.PublicKey)

	actionJSON := []byte(`{"amount":100,"currency":"USD"}`)
	ahsh, err := oer.ActionHashHex(actionJSON)
	if err != nil {
		exitErr(err)
	}
	pid := fmt.Sprintf("%064x", sha256.Sum256([]byte("pid-placeholder")))

	base := payload{
		Ahsh: ahsh,
		Ctx:  "oer_conf",
		Dcsn: "PERMIT",
		Exp:  4000000000,
		Iat:  2000000000,
		Pid:  pid,
		Rid:  "conf_vec_1",
		Sid:  "fin.transfer",
		V:    1,
	}

	baseBytes, err := marshalPayload(base)
	if err != nil {
		exitErr(err)
	}
	sigOK := sign(priv, baseBytes)
	validReceipt := wireReceipt(baseBytes, sigOK)

	// invalid_signature: flip first signature byte (still 64 bytes)
	sigBad := append([]byte(nil), sigOK...)
	sigBad[0] ^= 0xff
	invalidSigReceipt := wireReceipt(baseBytes, sigBad)

	// expired_receipt
	expired := base
	expired.Exp = 2000000100
	expired.Rid = "conf_vec_expired"
	expiredBytes, err := marshalPayload(expired)
	if err != nil {
		exitErr(err)
	}
	expiredReceipt := wireReceipt(expiredBytes, sign(priv, expiredBytes))

	// action_hash_mismatch: wrong ahsh, still signed over that payload
	badAh := payload{
		Ahsh: "0000000000000000000000000000000000000000000000000000000000000000",
		Ctx:  "oer_conf",
		Dcsn: "PERMIT",
		Exp:  4000000000,
		Iat:  2000000000,
		Pid:  pid,
		Rid:  "conf_vec_ahsh",
		Sid:  "fin.transfer",
		V:    1,
	}
	badAhBytes, err := marshalPayload(badAh)
	if err != nil {
		exitErr(err)
	}
	ahshMismatchReceipt := wireReceipt(badAhBytes, sign(priv, badAhBytes))

	// action_hash_case_mismatch: uppercase ahsh should fail exact case-sensitive comparison
	caseMismatch := base
	caseMismatch.Ahsh = strings.ToUpper(ahsh)
	caseMismatch.Rid = "conf_vec_ahsh_case"
	caseMismatchBytes, err := marshalPayload(caseMismatch)
	if err != nil {
		exitErr(err)
	}
	ahshCaseMismatchReceipt := wireReceipt(caseMismatchBytes, sign(priv, caseMismatchBytes))

	// unknown_protocol_version
	v2 := base
	v2.V = 2
	v2.Rid = "conf_vec_v2"
	v2Bytes, err := marshalPayload(v2)
	if err != nil {
		exitErr(err)
	}
	v2Receipt := wireReceipt(v2Bytes, sign(priv, v2Bytes))

	out := map[string]interface{}{
		"suite":          "OER-1",
		"description":    "Deterministic OER v1 conformance vectors (Ed25519 seed fixed in gen-oer-fixtures).",
		"public_key_hex": hex.EncodeToString(pub),
		"action_json_b64": base64.StdEncoding.EncodeToString(actionJSON),
		"vectors": []map[string]interface{}{
			{
				"name":             "valid_receipt",
				"expected":         "PASS",
				"receipt":          validReceipt,
				"expected_surface": "fin.transfer",
				"now_unix":         int64(2000001000),
			},
			{
				"name":             "invalid_signature",
				"expected":         "FAIL",
				"receipt":          invalidSigReceipt,
				"expected_surface": "fin.transfer",
				"now_unix":         int64(2000001000),
			},
			{
				"name":             "expired_receipt",
				"expected":         "FAIL",
				"receipt":          expiredReceipt,
				"expected_surface": "fin.transfer",
				"now_unix":         int64(2000000200),
			},
			{
				"name":             "surface_mismatch",
				"expected":         "FAIL",
				"receipt":          validReceipt,
				"expected_surface": "other.surface",
				"now_unix":         int64(2000001000),
			},
			{
				"name":             "action_hash_mismatch",
				"expected":         "FAIL",
				"receipt":          ahshMismatchReceipt,
				"expected_surface": "fin.transfer",
				"now_unix":         int64(2000001000),
			},
			{
				"name":             "action_hash_case_mismatch",
				"expected":         "FAIL",
				"receipt":          ahshCaseMismatchReceipt,
				"expected_surface": "fin.transfer",
				"now_unix":         int64(2000001000),
			},
			{
				"name":             "unknown_protocol_version",
				"expected":         "FAIL",
				"receipt":          v2Receipt,
				"expected_surface": "fin.transfer",
				"now_unix":         int64(2000001000),
			},
		},
	}

	raw, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		exitErr(err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		exitErr(err)
	}
	if err := os.WriteFile(outPath, raw, 0o644); err != nil {
		exitErr(err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", outPath)
}

func findRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		exitErr(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "conformance", "oer-v1")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			exitErr(fmt.Errorf("could not find repo root (conformance/oer-v1) from cwd"))
		}
		dir = parent
	}
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

package oer

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type vectorFile struct {
	Suite          string `json:"suite"`
	PublicKeyHex   string `json:"public_key_hex"`
	ActionJSONB64  string `json:"action_json_b64"`
	Vectors        []struct {
		Name            string `json:"name"`
		Expected        string `json:"expected"`
		Receipt         string `json:"receipt"`
		ExpectedSurface string `json:"expected_surface"`
		NowUnix         int64  `json:"now_unix"`
	} `json:"vectors"`
}

func TestOERConformanceVectors(t *testing.T) {
	path := vectorsJSONPath(t)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}
	var vf vectorFile
	if err := json.Unmarshal(raw, &vf); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}
	if vf.Suite != "OER-1" {
		t.Fatalf("suite: got %q", vf.Suite)
	}
	pubBytes, err := hex.DecodeString(vf.PublicKeyHex)
	if err != nil || len(pubBytes) != ed25519.PublicKeySize {
		t.Fatalf("public key: %v", err)
	}
	pub := ed25519.PublicKey(pubBytes)
	actionJSON, err := base64.StdEncoding.DecodeString(vf.ActionJSONB64)
	if err != nil {
		t.Fatalf("action b64: %v", err)
	}

	for _, c := range vf.Vectors {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			err := Verify(c.Receipt, actionJSON, c.ExpectedSurface, pub, c.NowUnix)
			switch c.Expected {
			case "PASS":
				if err != nil {
					t.Fatalf("expected PASS, got %v", err)
				}
			case "FAIL":
				if err == nil {
					t.Fatal("expected FAIL, got nil")
				}
			default:
				t.Fatalf("unknown expected %q", c.Expected)
			}
		})
	}
}

func vectorsJSONPath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	dir := filepath.Dir(thisFile)
	// tools/oer-verifier-go -> repo root
	root := filepath.Join(dir, "..", "..")
	path := filepath.Join(root, "conformance", "oer-v1", "fixtures", "vectors.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("vectors file %s: %v", path, err)
	}
	return path
}

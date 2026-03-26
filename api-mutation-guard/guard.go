// Package apimutationguard provides HTTP middleware that blocks requests unless an OER v1 receipt verifies offline.
// See docs/integrations/api_mutation_guard.md and docs/specs/VERIFIER_SDK.md.
package apimutationguard

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	oer "github.com/TrigGuard-AI/TrigGuard/tools/oer-verifier-go"
)

// MaxBodyBytes caps request body read size for hashing / verification (fail-closed if exceeded).
const MaxBodyBytes = 1 << 20

const blockedMsg = "OER authorization failed — mutation blocked"

var receiptReplayCache = struct {
	sync.Mutex
	used map[string]int64
}{
	used: make(map[string]int64),
}

// Middleware wraps next and requires a valid OER for the raw request body bytes (action JSON) and headers:
//   TG-Execution-Receipt — wire receipt
//   TG-Surface — must match payload sid
//   TG-Public-Key — Ed25519 public key as 64 hex characters (32 bytes)
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receipt := r.Header.Get("TG-Execution-Receipt")
		surface := r.Header.Get("TG-Surface")
		pubHex := r.Header.Get("TG-Public-Key")

		if receipt == "" || surface == "" || pubHex == "" {
			http.Error(w, blockedMsg, http.StatusForbidden)
			return
		}
		h := sha256.Sum256([]byte(receipt))
		receiptID := hex.EncodeToString(h[:])

		pub, err := parsePublicKeyHex(pubHex)
		if err != nil {
			http.Error(w, blockedMsg, http.StatusForbidden)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, blockedMsg, http.StatusForbidden)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		receiptReplayCache.Lock()
		if ts, ok := receiptReplayCache.used[receiptID]; ok {
			if time.Now().Unix()-ts < 60 {
				receiptReplayCache.Unlock()
				http.Error(w, "OER replay detected — mutation blocked", http.StatusForbidden)
				return
			}
		}
		receiptReplayCache.used[receiptID] = time.Now().Unix()
		receiptReplayCache.Unlock()

		if err := oer.Verify(receipt, body, surface, pub, time.Now().Unix()); err != nil {
			http.Error(w, blockedMsg, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parsePublicKeyHex(s string) (ed25519.PublicKey, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "0x")
	if len(s) != 64 {
		return nil, fmt.Errorf("public key must be 64 hex characters")
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("public key must decode to 32 bytes")
	}
	return ed25519.PublicKey(raw), nil
}

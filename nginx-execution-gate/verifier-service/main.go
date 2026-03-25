package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	oer "github.com/TrigGuard-AI/TrigGuard/tools/oer-verifier-go"
)

type verifyRequest struct {
	Receipt    string `json:"receipt"`
	Surface    string `json:"surface"`
	Action     string `json:"action"`
	PublicKey  string `json:"public_key"`
}

type verifyResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/verify", handleVerify)

	addr := ":8081"
	log.Printf("verifier-service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, verifyResponse{OK: false, Error: "method not allowed"})
		return
	}

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	var req verifyRequest
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusForbidden, verifyResponse{OK: false, Error: "bad request"})
		return
	}
	if req.Receipt == "" || req.Surface == "" || req.Action == "" || req.PublicKey == "" {
		writeJSON(w, http.StatusForbidden, verifyResponse{OK: false, Error: "missing fields"})
		return
	}

	actionBytes, err := decodeAction(req.Action)
	if err != nil {
		writeJSON(w, http.StatusForbidden, verifyResponse{OK: false, Error: "bad action"})
		return
	}

	pub, err := parsePublicKeyHex(req.PublicKey)
	if err != nil {
		writeJSON(w, http.StatusForbidden, verifyResponse{OK: false, Error: "bad public_key"})
		return
	}

	if err := oer.Verify(req.Receipt, actionBytes, req.Surface, pub, time.Now().Unix()); err != nil {
		writeJSON(w, http.StatusForbidden, verifyResponse{OK: false, Error: "verification failed"})
		return
	}

	writeJSON(w, http.StatusOK, verifyResponse{OK: true})
}

func decodeAction(s string) ([]byte, error) {
	// For proxy-friendly transport we accept either:
	// - base64 std encoding of the UTF-8 JSON bytes, or
	// - raw JSON string bytes (as-is)
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return []byte(s), nil
}

func parsePublicKeyHex(s string) (ed25519.PublicKey, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "0x")
	if len(s) != 64 {
		return nil, fmt.Errorf("public key must be 64 hex chars")
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("public key must be 32 bytes")
	}
	return ed25519.PublicKey(raw), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

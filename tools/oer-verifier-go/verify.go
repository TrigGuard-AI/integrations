// Package oer implements offline verification for Open Execution Receipts (OER v1).
// See docs/specs/OER_PROTOCOL.md and docs/specs/CANONICAL_HASHING.md.
// No network I/O.
package oer

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

const supportedVersion = 1

// ErrVerify is returned when verification fails (signature, binding, expiry, or version).
var ErrVerify = errors.New("oer: verification failed")

// Verify checks an OER wire string: base64url(payload).base64url(signature).
// actionJSON is the raw JSON bytes for the action to hash (canonical per spec).
// expectedSurface must match payload sid. nowUnix is current Unix seconds (caller clock).
// publicKey must be ed25519.PublicKey (32 bytes).
func Verify(receipt string, actionJSON []byte, expectedSurface string, publicKey ed25519.PublicKey, nowUnix int64) error {
	if len(publicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("%w: invalid public key size", ErrVerify)
	}
	parts := strings.Split(receipt, ".")
	if len(parts) != 2 {
		return fmt.Errorf("%w: expected two segments", ErrVerify)
	}
	payloadBytes, err := b64URLDecode(parts[0])
	if err != nil {
		return fmt.Errorf("%w: payload decode: %v", ErrVerify, err)
	}
	sig, err := b64URLDecode(parts[1])
	if err != nil {
		return fmt.Errorf("%w: signature decode: %v", ErrVerify, err)
	}
	if len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("%w: bad signature length", ErrVerify)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return fmt.Errorf("%w: payload JSON: %v", ErrVerify, err)
	}

	v, ok := payload["v"]
	if !ok {
		return fmt.Errorf("%w: missing v", ErrVerify)
	}
	vf, ok := toFloat(v)
	if !ok || vf != float64(supportedVersion) {
		return fmt.Errorf("%w: unsupported protocol version", ErrVerify)
	}

	if !ed25519.Verify(publicKey, payloadBytes, sig) {
		return fmt.Errorf("%w: Ed25519 signature", ErrVerify)
	}

	exp, ok := toFloat(payload["exp"])
	if !ok {
		return fmt.Errorf("%w: exp", ErrVerify)
	}
	if nowUnix >= int64(exp) {
		return fmt.Errorf("%w: expired", ErrVerify)
	}

	sid, _ := payload["sid"].(string)
	if sid != expectedSurface {
		return fmt.Errorf("%w: sid mismatch", ErrVerify)
	}

	wantAH, _ := payload["ahsh"].(string)
	got, err := ActionHashHex(actionJSON)
	if err != nil {
		return fmt.Errorf("%w: action hash: %v", ErrVerify, err)
	}
	if wantAH != got {
		return fmt.Errorf("%w: ahsh mismatch", ErrVerify)
	}
	return nil
}

func toFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func b64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// ActionHashHex returns lowercase hex SHA-256 of canonical JSON for actionJSON.
func ActionHashHex(actionJSON []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(actionJSON))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return "", err
	}
	cj, err := canonicalJSON(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(cj))
	return hex.EncodeToString(sum[:]), nil
}

func jsonMarshalNoHTMLEscape(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

func canonicalJSON(v interface{}) (string, error) {
	switch x := v.(type) {
	case nil:
		return "null", nil
	case bool:
		if x {
			return "true", nil
		}
		return "false", nil
	case json.Number:
		return x.String(), nil
	case string:
		b, err := jsonMarshalNoHTMLEscape(x)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case []interface{}:
		if len(x) == 0 {
			return "[]", nil
		}
		parts := make([]string, len(x))
		for i, e := range x {
			s, err := canonicalJSON(e)
			if err != nil {
				return "", err
			}
			parts[i] = s
		}
		return "[" + strings.Join(parts, ",") + "]", nil
	case map[string]interface{}:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			kb, err := jsonMarshalNoHTMLEscape(k)
			if err != nil {
				return "", err
			}
			vb, err := canonicalJSON(x[k])
			if err != nil {
				return "", err
			}
			parts = append(parts, string(kb)+":"+vb)
		}
		return "{" + strings.Join(parts, ",") + "}", nil
	default:
		return "", fmt.Errorf("unsupported JSON type %T", x)
	}
}

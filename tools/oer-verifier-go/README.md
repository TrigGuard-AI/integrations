# OER verifier (Go, reference)

Minimal **offline** verifier for **Open Execution Receipt (OER) v1** per [`docs/specs/OER_PROTOCOL.md`](../../docs/specs/OER_PROTOCOL.md) and [`docs/specs/CANONICAL_HASHING.md`](../../docs/specs/CANONICAL_HASHING.md).

- **No network calls**
- **&lt;300 lines** in `verify.go`

Not a substitute for conformance tests; copy or vendor for integration as needed.

## Verification (OER v1)

From a wire `receipt` and the live `actionJSON` / `expectedSurface`, `Verify` runs these checks in order:

| Step | What it verifies |
|------|------------------|
| 1 | **Parse payload JSON** — decode wire payload bytes and unmarshal JSON |
| 2 | **`v`** — payload field is **1** (OER v1) |
| 3 | **Signature** — Ed25519 over the exact UTF-8 payload bytes |
| 4 | **`exp`** — not expired at `nowUnix` |
| 5 | **`sid`** — equals `expectedSurface` |
| 6 | **`ahsh`** — exact (case-sensitive) match with SHA-256 of canonical JSON of `actionJSON` ([CANONICAL_HASHING.md](../../docs/specs/CANONICAL_HASHING.md)) |

## API

```go
func Verify(receipt string, actionJSON []byte, expectedSurface string, publicKey ed25519.PublicKey, nowUnix int64) error
```

- **`receipt`**: wire string `base64url(payload).base64url(signature)`.
- **`actionJSON`**: UTF-8 JSON bytes of the action to hash (must match what was hashed at issuance).
- **`expectedSurface`**: must equal payload `sid`.
- **`publicKey`**: 32-byte Ed25519 public key.
- **`nowUnix`**: current time in Unix seconds (for `exp` check).

Returns `nil` on success; `ErrVerify` or wrapped errors on failure.

## Example

```go
package main

import (
	"crypto/ed25519"
	"fmt"

	oer "github.com/TrigGuard-AI/TrigGuard/tools/oer-verifier-go"
)

func main() {
	receipt := "<from issuer>"
	action := []byte(`{"amount":100,"currency":"USD"}`)
	surface := "fin.transfer"
	pub := ed25519.PublicKey(/* 32 bytes from /.well-known/trigguard or pin */)
	var now int64 = 1742812005

	if err := oer.Verify(receipt, action, surface, pub, now); err != nil {
		fmt.Println("blocked:", err)
		return
	}
	fmt.Println("ok")
}
```

Build from repo root:

```bash
cd tools/oer-verifier-go && go build -o /dev/null .
```

Package import path: `github.com/TrigGuard-AI/TrigGuard/tools/oer-verifier-go` with package name **`oer`** — use `go.work` or run inside this module for local dev.

From this directory:

```bash
go test -c . 2>/dev/null || go build .
```

## Limitations

- Does not validate `dcsn` (e.g. PERMIT vs DENY) for business logic; enforcement layers should still refuse execution on DENY/SILENCE per deployment rules.

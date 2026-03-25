# Examples

Runnable and copy-paste-friendly samples for integrating TrigGuard-style **execution gates** and policy flows. These are **not** a substitute for the normative protocol docs — see [`docs/protocol/TRIGGUARD_EXECUTION_PROTOCOL.md`](../docs/protocol/TRIGGUARD_EXECUTION_PROTOCOL.md) and [`ARCHITECTURE.md`](../ARCHITECTURE.md).

| Example | What it shows |
|---------|----------------|
| [`agent_tool_guard_demo.py`](agent_tool_guard_demo.py) | **~100 lines:** agent → Swift `POST /decide` → PERMIT/DENY → simulated execution |
| [`agent_demo.py`](agent_demo.py) | **~100 lines:** agent → **`remote-eval`** `POST /v1/evaluate` → decision → simulated execution (Node + Swift) |
| [`basic-ai-gate/`](basic-ai-gate/) | Normalizing agent intent into an execution envelope |
| [`execution-gateway-demo/`](execution-gateway-demo/) | Gateway + protected execution path (Node) |
| [`payment-gate/`](payment-gate/) | Payment authorization flow |
| [`payment_guard/`](payment_guard/) | Client-side guard pattern |
| [`js-demo/`](js-demo/) | Minimal JS payment demo |
| [`demo-trading/`](demo-trading/) | Trading-style requests |
| [`data-export-gate/`](data-export-gate/) | Data export authorization |
| [`terraform-infra-gate/`](terraform-infra-gate/) | Infra change gate |
| [`terraform-guard-example.md`](terraform-guard-example.md) | OER-gated `terraform apply` (plan JSON + wrapper) |
| [`github-deploy-gate/`](github-deploy-gate/) | Deploy workflow shape |
| [`github_action/`](github_action/) | Workflow snippets (YAML) |
| [`api-guard-example.md`](api-guard-example.md) | HTTP mutation routes gated by OER middleware |
| [`github-enforced-deploy.yml`](github-enforced-deploy.yml) | OER receipt gate before apply (composite action) |
| [`deploy_release/`](deploy_release/) | Release / GitHub Action notes |
| [`swift-demo/`](swift-demo/) | Swift payment demo |

## Example Status Matrix

| Example | Status | Trust level | Uses real verifier? | Production intent |
|---------|--------|-------------|---------------------|-------------------|
| `agent_tool_guard_demo.py` | illustrative | low | no (uses canonical /decide) | demo-only |
| `agent_demo.py` | illustrative | low | no (uses remote-eval /v1/evaluate) | demo-only |
| `reference-verifier/` | production reference | medium | yes | integration reference |
| `execution-gateway-demo/` | illustrative | low | no | demo-only |
| `payment-gate/` | illustrative | low | no | demo-only |
| `terraform-infra-gate/` | illustrative | low | yes | reference pattern |
| `github-deploy-gate/` | illustrative | low | yes | reference pattern |
| `data-export-gate/` | illustrative | low | no | demo-only |
| `swift-demo/` | experimental | low | no | demo-only |

Add new examples under this directory with a **README.md** that states assumptions, how to run, and what is out of scope.

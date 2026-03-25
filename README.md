# TrigGuard Integrations

Official and community integrations (CI, gateways, infra) that connect external systems to TrigGuard execution governance.

## Role in the TrigGuard Ecosystem

This repository catalogs or hosts integration code and documentation. Each integration **conforms to** the published protocol contract; protocol definitions remain in **trigguard-protocol**.

## Core terms

Integrations must respect: **Execution Surface**, **Protocol Contract**, **Protocol Fingerprint**, **Audit Bundle**, **Verifier**, and **OER** (Operational Evidence Record).

## Authority Hierarchy

1. **trigguard-protocol** — canonical protocol specification, fingerprints, and conformance artifacts  
2. **trigguard-core-reference** — reference implementation for alignment  
3. **trigguard-js** — client SDK surfaces  
4. **docs** — explanatory documentation  
5. **cloud** — hosted enforcement layer  
6. **site** — public discovery (trigguardai.com)  
7. **TrigGuard** — combined monorepo (demos, services, tooling, release verification); see also ecosystem integration registry in that repo.

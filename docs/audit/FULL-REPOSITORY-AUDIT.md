# Full repository audit

Baseline for Orchestra 4.5.0, taken on branch `release/4.5.0` at `60cf9e7` before any 4.5 edit. Numbers below are from commands run in this pass.

## What was run

| Check | Result |
| --- | --- |
| `node runtime/tools/hygiene.js` | 753 files, VERSION 4.2.0, failures 0 |
| `node runtime/tools/check-registry.js` | resources 180, problems 0 |
| `node runtime/tools/check-host-stack.js` | 78 skills, problems 0 |
| `node runtime/tools/check-routing.js` | 5 routing copies, 22 routes, problems 0 |
| `node runtime/tools/check-route-agreement.js` | 25 skills routed, problems 0 |
| `node runtime/tools/check-41.js` | 12 domains, 25 libraries, 7 MCP rows, problems 0 |
| `node runtime/tools/check-design-contract.js templates/contracts premium` | problems 0 |
| `node runtime/tools/inventory.mjs` | 753 tracked files written to `FILE-INVENTORY.json` |
| `go test ./...` | started in `runtime/`; result is recorded in `TEST-COVERAGE-REPORT.md` when the process exits |

GitHub Actions on commit `60cf9e7` was green for Test and Hygiene in the 4.2.0 release. This audit does not treat that as a substitute for the local suite.

## Inventory

| Group | Tracked files |
| --- | ---: |
| skills | 357 |
| runtime | 179 |
| templates | 76 |
| docs | 35 |
| installers (`kit/`) | 28 |
| protocols | 17 |
| root | 14 |
| benchmarks | 13 |
| tooling | 13 |
| registries | 11 |
| workspace-template | 6 |
| cli | 2 |
| ci | 2 |

`skills/` also holds untracked host-local copies (`docx`, `pdf`, `pptx`, `xlsx`, `fullstack-mastery`, `creative-3d-web`). `check-host-stack.js` excludes them from the 78. They are not canonical.

## What is wrong, from reading the installers and CI

1. `kit/sync-ides.ps1` and `kit/bootstrap.ps1` delete a destination skill directory and copy over it. A skill the user edited locally is gone. `-Prune` can move extras using a backup path left over from 3.3.5.
2. `bootstrap.ps1` overwrites `AGENTS.md`, `CLAUDE.md`, and `.cursorrules` in the target repo.
3. Version `4.2.0` is a string in `VERSION`, `hygiene.js`, `check-host-stack.js`, the host-stack schema, and the doctor. Historical docs also contain older version numbers. A blind replace would rewrite history.
4. CI runs `jq empty` on `runtime/schemas/*.json` only. `registries/*.json` is not validated against `registries/schemas/`.
5. The routing table is copied by hand into five adapters. `check-routing.js` compares the copies. It does not generate them.
6. `runtime/benchmarks/c1-premium-landing.json` is `not_run`. Comparison claims c1–c6 are disabled. The README must keep saying that until a record exists.
7. `docs/AGENT_PORTABILITY.md` still describes Orchestra 3.3.
8. `printUsage` in `runtime/cmd/orchestra/main.go` still says Orchestra 4.1.0.

## What this audit did not do

It did not open all 753 files. `FILE-INVENTORY.json` classifies each tracked path. Purpose lines for non-skill files are the group, not a hand-written summary of every file.

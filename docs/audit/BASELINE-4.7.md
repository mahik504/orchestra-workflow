# Baseline 4.7

Measured on the 4.6.0 tree before the truthful-runtime edits, on 2026-11-12.

## Go

`go test -count=1 ./...` from `runtime/` exited 0. The runner reported 770 passed in 25 packages. Elapsed about 443 seconds.

## Node

| Check | Result |
| --- | --- |
| `node runtime/tools/hygiene.js` | exit 0, 1210 files, failures 0, VERSION 4.6.0 |
| `node runtime/tools/check-schemas.mjs` | exit 0, 13 registries, 14 schemas, problems 0 |
| `node runtime/tools/check-registry.js` | exit 0, 208 resources, problems 0 |
| `node runtime/tools/check-skills.js` | exit 0, 99 skills, 27 aliases, problems 0 |
| `node runtime/tools/check-routing.js` | exit 0, 30 routes, problems 0 |
| `node runtime/tools/check-route-agreement.js` | exit 0, 35 skills, problems 0 |
| `node runtime/tools/check-versions.js` | exit 0, current 4.6.0, graph field 4.0.0, problems 0 |
| `node runtime/tools/check-host-stack.js` | exit 0, 99 skills, problems 0 |
| `node runtime/tools/check-41.js` | exit 0, problems 0 |
| `node runtime/tools/check-mermaid.mjs` | exit 0, 27 diagrams, problems 0 |
| `node runtime/tools/measure-context.mjs` | exit 0, AGENTS.md 9929 characters, skill descriptions 24064, ow-* over 300: 0 |
| `node --test runtime/tools/*.test.mjs` | exit 0, 13 passed, 0 failed |
| `node runtime/tools/check-design-contract.js` | exit 2, usage line only. It needs a directory and a profile. |

## Eval harness

`skills/ow-eval` and `runtime/benchmarks/orchestra-47/` define arms A (clean HOME), B (control plane), and C (routed skills). The six fixtures are `not_run`. `compare.go` claims stay `none`. Promotion is `PromoteSkill`: success change not negative, cost at most 1.10, at least 3 trials.

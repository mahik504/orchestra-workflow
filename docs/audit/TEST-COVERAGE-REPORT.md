# Test coverage report

## Node checks at `60cf9e7` on `release/4.5.0`

All of these exited 0:

- `hygiene.js` — 753 files, failures 0
- `check-registry.js` — 180 resources, problems 0
- `check-host-stack.js` — 78 skills, problems 0
- `check-routing.js` — problems 0
- `check-route-agreement.js` — problems 0
- `check-41.js` — problems 0
- `check-design-contract.js` on `templates/contracts` profile `premium` — problems 0

`inventory.mjs` wrote 753 rows. It is a generator, not a test.

## Go

`go test ./... -count=1 -timeout 20m` was started from `runtime/` for this baseline. The process result is appended at the bottom of this file when it exits. Until that line exists, this report does not claim a local Go pass.

The public CI run on `60cf9e7` passed Test and Hygiene during the 4.2.0 release. Six tests in that tree are skipped by build tags or missing tools. A skip is not a pass, and this file does not list them as passed.

## What has no test

- `kit/sync-ides.ps1` and `kit/bootstrap.ps1` have no temp-home test.
- `runtime/tools/*.js` are not run under `node --test`.
- Registry JSON is not validated with a JSON Schema checker in CI. `jq empty` only checks that `runtime/schemas/*.json` parses.
- The c1 benchmark record is `not_run`. No Layer A, axe, or Lighthouse numbers exist for it.

## Local Go baseline

`go test ./... -count=1 -timeout 20m` from `runtime/` exited 0. The runner reported 761 passed in 24 packages. Elapsed about 8.3 minutes. That run started before the 4.5 version-reader edit in `main.go` and finished after it, so `cmd/orchestra` may have compiled the version-reader. The count is the observed runner line, not a claim about coverage percentage.

## 4.5 suite

A later `go test ./... -count=1 -timeout 20m` from `runtime/`, after the review fixes, exited 0. The runner reported 769 passed in 25 packages. Elapsed about 8.1 minutes. That is the suite for this release, not a coverage percentage.

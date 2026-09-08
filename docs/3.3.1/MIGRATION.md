# 3.3.1 migration

Live rules live in `AGENTS.md` and `protocols/`. This file is the audit trail.

## Version

`3.3.0` → `3.3.1`. Same OS. One conductor. Do not treat this as a toolbox expansion.

## Engine

Design Lab states:

| Old | New |
| --- | --- |
| `PENDING` → `APPROVED` after `DESIGN.md` | `PENDING` → `CONTRACT_APPROVED` after locked `DESIGN.md` → `APPROVED` after golden stills |
| Product UI writable at `APPROVED` | Product UI writable only at `APPROVED` / `BYPASSED` / `NOT_REQUIRED` |

`Approve` is contract-only. Call `ApproveStills(approvedBy, desktopPath, mobilePath, note)`. Missing files or an empty note refuse.

Named-reference surveys: `OfferSurvey` accepts **23** or **3**.

## Status aliases

Existing catalog values stay valid:

| Existing | 3.3.1 reading |
| --- | --- |
| `CURATED_OPTIONAL` | ≈ project-only |
| `BOOKMARK` | ≈ reference |
| `ARCHIVED` / `LEFTOVER` | ≈ obsolete |

New enum members: `PROJECT_ONLY`, `CATALOG`, `EXPERIMENTAL`, `PRIVILEGED`, `AUTH_REQUIRED`, `OBSOLETE`. Only named candidates were reclassified. The rest of the catalog was not gutted.

## Hosts

Bump `VERSION`, hygiene `EXPECTED_VERSION`, `registries/host-stack.json`, Antigravity `BEGIN ORCHESTRA 3.3.1` markers, conductor/vault skill headers. Overlay copies must match.

Do not tag until the operator says so.

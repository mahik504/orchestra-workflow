# Orchestra 3.3.2 — Phase 3G.2.1 report

**Preflight only.** The installed agent listed eight exact Grok 4.7 ids. This pass did not pick one. Smoke did not run. T1 and T3 stay ineligible. `claim` is `none`.

VERSION stays `3.3.1`. No commit. No Orchestra code change. No T1 or T3 pair. No 3H.

## Catalog

- Basename: `agent.cmd`
- Version: `2026.09.08-6caf4ff`
- Command: `agent models`
- Each id below was the id field of a live row, so `--model` can take that string.
- `grok-4.7` itself was not a row. Nothing was rewritten into that string.

| Exact id | Display name | Variant |
| --- | --- | --- |
| `grok-4.7-low` | Grok 4.7 Low | low |
| `grok-4.7-low-fast` | Grok 4.7 Low Fast | low, fast |
| `grok-4.7-medium` | Grok 4.7 Medium | medium |
| `grok-4.7-medium-fast` | Grok 4.7 Medium Fast | medium, fast |
| `grok-4.7-high` | Grok 4.7 High | high |
| `grok-4.7-high-fast` | Grok 4.7 High Fast | high, fast |
| `grok-4.7-xhigh` | Grok 4.7 Extra High | xhigh |
| `grok-4.7-xhigh-fast` | Grok 4.7 Extra High Fast | xhigh, fast |

Grok 4.5 and Grok 4.6 rows were dropped. The rest of the catalog is not in this report.

## Smoke

Skipped. More than one Grok 4.7 id was present. `marker.txt` was not created. No disposable directory was left behind.

## Eligibility

T1 and T3 are not eligible. Eligibility needs exactly one id and a smoke that writes `marker.txt`. This catalog has eight ids and no smoke.

## Protected paths

Fingerprints were taken before `agent models` and again after it. No class moved.

- `method_repo`: unchanged
- `fixture`: unchanged
- `cursor_config`: unchanged
- `mcp_config`: unchanged. The home MCP file was present. The Cursor user MCP file was absent both times.
- `brain`: `ORCHESTRA_MEMORY_PATH` was unset, so no brain file was fingerprinted. It stayed unset.

## Stop

`claim: none`. Do not enable T2. Do not start 3H. Do not alias these ids inside Orchestra. A later T1 and T3 pair waits until one exact id is named.

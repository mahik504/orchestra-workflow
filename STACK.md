# Stack

| Component | Technology |
|---|---|
| Core Engine | Go 1.22+ |
| Configuration | JSON Schemas |
| Visual QA | Playwright / Chromium |
| Performance QA | Lighthouse |
| Security | Semgrep |
| Target Output | Agent-Agnostic Execution Manifests |
| Control plane | Orchestra 3.3.1 (`AGENTS.md` + engine) |
| Hosts | Cursor, Antigravity, Claude Code, jcode (executor) |
| Resource knowledge | `registries/resources.json` + `design-resource-graph.json` |
| Local OSS metadata | `registries/oss-index.json` (not the Brain) |

Per-product libraries (R3F, GSAP, shadcn, Supabase) install **in the app repo**, not globally. GetLayers is not in the stack until purchase.

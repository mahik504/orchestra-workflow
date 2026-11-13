# Agent portability

Historical note. The 3.3 table below is the old host list. Current hosts, including Codex, OpenCode, and Kimi, are in `docs/hosts.md`. Orchestra 5.5.0 still uses one contract. Hosts are adapters.

Orchestra 3.3 is one contract. Hosts are adapters. They do not invent a second plan.

| Host | Adapter | Native strength |
| --- | --- | --- |
| Cursor | `.cursorrules` | Bulk implementation, in-file diffing |
| Antigravity | `kit/antigravity/MASTER-PROMPT.md` | Visual QA, architecture planning |
| Claude Code | `CLAUDE.md` | Terminal execution, backend refactor |

See [`AGENTS.md`](../AGENTS.md), [`WORKFLOW.md`](../WORKFLOW.md), and [`docs/adapters.md`](adapters.md).

A host may own a capability the others lack. Map it. Do not clone plugin lists between hosts. Sync means the same contract, not the same installed extras.

# Agent portability

Orchestra 3.1 is one contract. Hosts are adapters. They do not invent a second plan.

| Host | Adapter | Native strength |
| --- | --- | --- |
| Cursor | `.cursorrules` | Bulk implementation, in-file diffing |
| Antigravity | `kit/antigravity/MASTER-PROMPT.md` | Visual QA, architecture planning |
| Claude Code | `CLAUDE.md` | Terminal execution, backend refactor |

See [`AGENTS.md`](../AGENTS.md), [`WORKFLOW.md`](../WORKFLOW.md), and [`docs/adapters.md`](adapters.md).

A host may own a capability the others lack. Map it. Do not clone plugin lists between hosts. Sync means the same contract, not the same installed extras.

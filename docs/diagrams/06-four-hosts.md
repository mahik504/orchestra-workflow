# Four hosts

One contract. Four adapters. jcode is an executor, not a second conductor.

```mermaid
flowchart LR
  src[orchestra-workflow]
  sync[kit/sync-ides.ps1]
  c[Cursor]
  a[Antigravity]
  cl[Claude Code]
  j[jcode]
  src --> sync
  sync --> c
  sync --> a
  sync --> cl
  sync --> j
```

Allowlist: `registries/host-stack.json`. Merge MCP configs; do not overwrite unrelated servers. Do not invent jcode MCP support.

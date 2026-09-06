# Design Lab

Write lock on `PREMIUM` / `EXPERIMENTAL` visual work.

```mermaid
flowchart TD
  named{Named site skill MCP pack or DESIGN.md?}
  survey[23 short cards]
  pick[One direction]
  contract[One sourced DESIGN.md]
  custom[Pasted DESIGN.md wins]
  gate[Human gate]
  impl[Implement]
  named -->|yes| contract
  named -->|no| survey --> pick --> contract
  contract --> gate
  custom --> gate
  gate -->|approved| impl
  gate -->|edit| survey
```

Tint: keep hierarchy, component grammar, spacing, interaction, shape. Replace color roles. Extend the new tokens. Do not clone trademarks or source.

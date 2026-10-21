# Design Lab

Write lock on `PREMIUM` / `EXPERIMENTAL` visual work. Two locks: `CONTRACT_APPROVED` then `ApproveStills`. Golden stills are pre-implementation evidence. A `DESIGN.md` is not visual evidence.

```mermaid
flowchart TD
  kind{What did the prompt name?}
  survey[23 short cards]
  trans[Exactly 3 translations]
  pick[One direction]
  contract[One sourced DESIGN.md]
  custom[Pasted DESIGN.md skips survey]
  stills[Golden stills desktop plus mobile]
  gate[Human ApproveStills]
  impl[Implement product UI]
  kind -->|open-ended PREMIUM| survey --> pick --> contract
  kind -->|named visual ref| trans --> contract
  kind -->|pasted DESIGN.md or named contract pack| custom --> contract
  contract --> stills --> gate
  gate -->|approved| impl
  gate -->|edit| contract
```

Named visual references do **not** skip to one `DESIGN.md`. Tint: keep hierarchy, component grammar, spacing, interaction, shape. Replace color roles. Extend the new tokens. Do not clone trademarks or source.

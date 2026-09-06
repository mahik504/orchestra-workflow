# Skill tiers

```mermaid
flowchart TD
  global[Global allowlist ~40]
  pack[Job packs markdown]
  opt[OPTIONAL catalog]
  rej[REJECTED dumps]
  global --> pack
  opt -.->|named| global
  rej -->|never load| x[Quarantine]
```

Never `skills add --all`. New internet skills: `skill-security-check` then wait for **update**.

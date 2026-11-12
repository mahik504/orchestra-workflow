---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-diagrams
description: Grounded Mermaid architecture diagrams from the real repo. Use for README and docs. Not a second orchestration system. Prefer a small diagram bundle, not one giant chart.
---

# Diagram generator (Orchestra wrapper)

Upstream: https://github.com/jovd83/ow-diagrams  

## When

README, ARCHITECTURE.md, onboarding a clone, explaining a capability graph. Skip for pixel UI and DSA homework.

## How

1. Read **this** repository (or the app repo). Do not invent boxes.
2. One abstraction, one reading direction, stable names, few crossings, a small bundle. Source method: Diagram Design. It is not a second skill. Archify is a reference, not a replacement.
3. Prefer a small set: context, containers, request flow, data model, capability graph.
4. Output Mermaid in markdown. Pretty-mermaid renders SVG if needed.
5. Do not generate 23 diagrams inline in README. README embeds a few; `docs/diagrams/` holds the rest.

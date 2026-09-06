---
name: diagram-generator
description: Grounded Mermaid architecture diagrams from the real repo. Use for README and docs. Not a second orchestration system. Prefer a small diagram bundle, not one giant chart.
---

# Diagram generator (Orchestra wrapper)

Upstream: https://github.com/jovd83/diagram-generator  

## When

README, ARCHITECTURE.md, onboarding a clone, explaining a capability graph. Skip for pixel UI and DSA homework.

## How

1. Read **this** repository (or the app repo). Do not invent boxes.
2. Prefer a small set: context, containers, request flow, data model, capability graph.
3. Output Mermaid in markdown. Pretty-mermaid renders SVG if needed.
4. Do not generate 23 diagrams inline in README. README embeds a few; `docs/diagrams/` holds the rest.

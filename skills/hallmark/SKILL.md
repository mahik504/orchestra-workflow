---
name: hallmark
description: Post-implementation design audit / anti-slop gates. Use after Playwright stills exist. Not a generator. Does not replace Impeccable or Taste.
---

# Hallmark (Orchestra wrapper)

Upstream: https://github.com/Nutlope/hallmark

Run **after** implementation and stills. Question: does this look intentional?

## When

PREMIUM / EXPERIMENTAL visual QA. Skip while Design Lab is pending. Skip for backend-only.

## How

1. Taste + Impeccable generate and refine. Hallmark audits.
2. Report gates that fail with a named fix. Do not restyle the whole page because one gate fired.
3. Ban-lists go stale. Prefer missing states (`:active`, `:focus-visible`, empty/error) over a new banned color.

Does not install Hallmark's site, CLI, or stop hooks. Orchestra remains the conductor.

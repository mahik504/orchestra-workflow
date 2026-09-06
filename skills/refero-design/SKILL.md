---
name: refero-design
description: Product UX research from real apps. Use when the job needs screens, flows, onboarding, dashboards, or mobile/iOS patterns. MCP is AUTH_REQUIRED. Not a clone factory. Not always-on.
---

# Refero design (Orchestra wrapper)

Upstream: https://github.com/referodesign/refero_skill  
MCP: `https://api.refero.design/mcp` (Bearer token in the **host**, never git).

Orchestra routes here for **real-product UX research**. Extract principles. Never copy branding, assets, copy, trademarks, or source.

## When

SaaS, dashboards, mobile, iOS, onboarding, account, billing, product flows. Skip for backend-only, papers, and operator HUDs.

## How

1. If MCP state is `AUTH_REQUIRED`, tell the human to Connect (see `docs/manual-setup/`). Do not invent a token.
2. Search screens/flows that match the brief.
3. Write anatomy notes into `DESIGN.md` / Vault Anatomy. Tint to the approved stack.
4. One named echo is enough. Do not dump 12,000 screens into context.

Available ≠ loaded. Refero does not become the conductor.

---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-component-discovery
description: Find an existing component before building one. Search the product, then shadcn, Magic UI, and React Bits. Adapt the pattern. Do not add a second component kit or leave an always-on MCP.
---

# Component discovery

Use this when the task needs a control, a card, or a section the product might already have.

## Search order

1. The product repo. A second button next to the existing one is a miss.
2. The route. One control loads `none` and still checks `interaction-components`.
3. shadcn if the product already uses it, or the human names it.
4. Magic UI or React Bits for a motion pattern, then rewrite it in the product's tokens.
5. 21st.dev only when the human has supplied a key. Do not invent one.

## Record

In the brief, name the source and the file you will change. If you install a registry component, pin the command in the product repo, not a global MCP.

## Skip

Do not load Aceternity, Magic UI, and shadcn together "to have options." Pick one source for this pass. Do not copy a gallery's markup and leave their fonts and colors in place.

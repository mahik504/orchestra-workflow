---
name: react-best-practices
description: Vercel React/Next performance rules. Use on React/Next apps for waterfalls, bundle size, server/client fetching, rerenders. Not a design skill. Fetch latest upstream rules when reviewing.
---

# React best practices (Orchestra wrapper)

Upstream: https://github.com/vercel-labs/agent-skills/tree/main/skills/react-best-practices  
License: Apache-2.0. Do **not** `npx skills add` the whole vercel-labs pack.

## When

React or Next.js implementation or review. Skip for Expo-only native, PySimpleGUI, and papers.

## How

1. Prefer the live skill at the URL above when the host can fetch it.
2. Check waterfalls, bundle, server performance, client fetching, rerenders.
3. Beautiful + slow is a fail on PREMIUM.

Does not replace `web-design-guidelines` (a11y/forms). Load both only when the job needs both.

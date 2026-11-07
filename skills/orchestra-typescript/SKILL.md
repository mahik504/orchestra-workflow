---
name: orchestra-typescript
description: TypeScript boundaries and checks before an edit. Use when the repo is TypeScript and the job is not only a React waterfall. Not a visual redesign.
---

# TypeScript — 4.1.0

Source method: Matt Pocock's public engineering notes. The `mattpocock/skills` pack is not installed. `react-best-practices` still owns React and Next waterfalls.

## Before the edit

1. State the assumption the types are hiding. If the assumption is wrong, stop and say so.
2. Parse input at the boundary. A Zod-shaped schema, or the project's existing parser, beats a cast.
3. Do not widen a type to `any` to finish the task. Narrow it or change the caller.
4. Do not introduce a second validation library when the repo already has one.

## Done

`tsc` or the repo's typecheck command is clean for the files you touched, or you report the error you left and why.

---
name: ow-retro
description: Append one reflection after a failed or finished wave.
---

Append one line to `.orchestra/reflections.jsonl` and update `.orchestra/STATE.md`. Do not open a second plan.

When the wave failed, also write `.orchestra/handoff/lesson.md` with three lines: the command that failed, the fix, and the command that passed. That file is a draft. `ow-vault` appends it to the private workspace only after a human says yes. Do not create a skill from it.

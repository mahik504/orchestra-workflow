---
name: ow-wave
description: Run a multi-session job as phases in worktrees. The chat stays the conductor. A failed verify writes a fix note and a reflection line.
---

# Waves

For work that does not fit one session. The conductor does not hand the plan to another framework.

## Preconditions

- The human has accepted the phase list.
- `using-git-worktrees` is how a phase gets its own tree.

## Procedure

1. Write `.orchestra/STATE.md`, `.orchestra/ROADMAP.md`, and `.orchestra/phases/NN-PLAN.md`.
2. Run one ticket in a fresh worktree. The worker sees that phase file, not the whole chat.
3. Verify the phase. On failure, write `NN-FIX.md` and append one line to `.orchestra/reflections.jsonl`.
4. The conductor integrates. Do not open a second plan.

## Output

The phase files, the reflection log, and the merged diff.

## Verify

STATE.md names the phase that just finished and the one that is next.

## Skip

Do not use waves for a one-file fix. Do not install another conductor.

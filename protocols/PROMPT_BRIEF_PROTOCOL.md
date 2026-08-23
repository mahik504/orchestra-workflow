# Prompt brief protocol (bounded)

Orchestra already had `templates/plan-brief.md` + the activation card. That **is** the prompt-improvement capability. No silent rewriter skill (those hide the rewrite and can expand context every turn).

## When

Vague ask **and** the job is architecture, high-visual, or multi-surface. Not for “fix this typo.”

## Bound

1. Fill `templates/execution-brief.md` **once** (or the holes in `plan-brief.md`).
2. Show the operator the brief. Do not hide it.
3. Implement against that brief.
4. **Stop.** Max one redesign pass after visual QA unless the operator asks.
5. No RALF / infinite re-prompt loop. No extra model call just to rephrase.

## Do not install

`Li-Bailiang/prompt-refine-skill` (silent, always-on). Context-engineering dumps.

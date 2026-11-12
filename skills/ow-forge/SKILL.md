---
name: ow-forge
description: Turn a plain request into a short brief before any other skill runs. Use at the start of a job. Do not rewrite the user's words in secret.
---

# Prompt forge

Run this before classify, then `ow-focus`. Then do the work. Ask only when the task is truly ambiguous.

## Brief

| Field | Write |
| --- | --- |
| Intent | What the user wants done |
| Kind | college, personal, hiring-cv, or none |
| Route | One capability from the graph |
| Acceptance | How we will know it worked |
| Refs | Named files, links, or library rows |
| Constraints | What must not change |
| Out of scope | What this pass will not do |

Show the brief in a few lines. Do not wait for approval unless a requirement that changes product behavior is missing.

## New product

When the intent or the acceptance line is unclear and the job is a new product, run `grilling` before the brief is treated as settled. Then `domain-modeling` writes the glossary and any ADR.

For a startup idea, ask one scope question: expand, hold, or reduce. Record the answer in the brief. Do not start a second product in the same pass.

## Load line

Before a write that loads something, name only what this task uses:

```
Load: skill-a, skill-b
MCP: playwright
Libraries: gsap
```

Omit a line that would be empty. Omit the whole note when the task loads nothing new. A skill the human names wins unless it is unsafe. A missing tool stops. Do not install it.


## Preconditions

- The task names the files or the acceptance line.

## Procedure

1. The task's paths and acceptance line win.
2. Write the brief fields that are still missing. Stop when the acceptance line is clear.

## Output

The change, the command that checked it, and what that command printed.

## Verify

Say the work is done only after that command has been run.

## Skip

Skip when the human says skip orchestra. Do not add a second conductor.

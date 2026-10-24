---
name: prompt-forge
description: Turn a plain request into a short brief before any other skill runs. Use at the start of a job. Do not rewrite the user's words in secret.
---

# Prompt forge

Run this before classify. Then do the work. Ask only when the task is truly ambiguous.

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

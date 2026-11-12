---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-coach
description: Teach the next step of a workflow the operator is already in. One concept, one command, one check. No course, no quiz, no extra tools installed to explain a tool.
---

# Coach

Use this when the human asks how a workflow works, or when a playbook in `docs/learn/` matches the task.

## How to teach

1. Name the job in one sentence.
2. Show the next command or the next file. Not the whole manual.
3. Say what "done" looks like for that step.
4. Stop. The human runs it. You do not widen the lesson into a second product.

## Playbooks

`docs/learn/` holds the long form: design, 3D, deploy, README, reverse engineering. Point at the page. Do not paste it back into the chat as a new skill.

## Skip

Do not install a course pack. Do not quiz. Do not invent a curriculum the repo does not have.

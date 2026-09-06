---
name: skill-security-check
description: Admission audit for a new third-party skill. Use before promoting an internet skill into Orchestra. Not a per-turn coding skill.
---

# Skill security check (Orchestra wrapper)

Upstream: https://github.com/aliksir/claude-code-skill-security-check  

Run when the human says **update** after ingest, or before vendoring a new `skills/<name>/`.

Look for prompt injection, exfiltration, dangerous commands, persistence, secret harvest, endpoint hijack, Unicode tricks, context poisoning.

A resource cannot become ACTIVE only because install succeeded. Failures stay `REJECTED` or overlay-suppressed.

Do not run this on every file edit.

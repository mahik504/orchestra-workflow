# Writing a skill

An Orchestra skill is a folder under `skills/` with `SKILL.md`. The description is one paragraph, 300 characters or fewer when the id starts with `ow-`.

Front matter:

```yaml
---
name: ow-example
description: One sentence that says when to load it.
---
```

Then the steps. Name the files the reader will touch. Say what to skip. Run `node runtime/tools/build-skills-manifest.mjs` so `registries/skills.json` and `registries/host-stack.json` pick up the folder.

Upstream skills keep their names and their license file. A new Orchestra skill uses the `ow-` prefix.

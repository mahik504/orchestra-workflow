---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-skill-advisor
description: Decide whether a task needs a new skill. Check skills/ first, then skills.sh and GitHub. Recommend use-existing, install-upstream, write ow-*, or install because the human insisted. Record that install. Do not run a global skills add on your own.
---

# Skill advisor

Use this before adding a skill, an MCP, or a plugin. `find-skills` searches. This file decides.

## Order

1. Read `registries/skills.json`. If a row already covers the job, load it. Say which id.
2. If the human named a package, check its license and the last commit. `find-skills` can run `npx skills find`. Do not run `npx skills add`.
3. Recommend one of four outcomes:
   - **use-existing** when a canonical skill already does the job.
   - **install-upstream** when the license is known, the security check found no secret-stealing or silent global install, and the human agrees.
   - **write ow-*** when the gap is Orchestra policy (routing, gates, receipts) or the upstream repo has no SPDX license.
   - **install-anyway** when the human insists after hearing the advisory. Install the named package only. Record it.

## Record

Write `.orchestra/user-installed.json` in the product repo:

```json
{ "id": "example", "source": "owner/repo", "reason": "human insisted", "date": "YYYY-MM-DD" }
```

Do not put tokens in that file.

## Hard limits

`npx skills add --all` stays refused. `-g` and `-y` stay off unless the human says those flags out loud for one package. A missing license means you write an original `ow-*` skill and credit the source, or you stop. You do not copy the unlicensed text.

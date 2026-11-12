---
name: ow-author
description: Add one skill, agent, rule, hook, command, or MCP entry for a host. Register it, check it, and leave it off a Load line until an eval says it earned the place.
---

# Author

Command `/ow-new`. One new capability per pass.

## Preconditions

- The gap is real. `skills/` does not already do the job.
- You can name the trigger, the output, and the stop condition.

## Procedure

1. Write the skill with description under 300 characters, then Preconditions, Procedure, Output, Verify, and Skip.
2. Add the agent or command only if a host needs a separate entry point. Commands call a route. They do not plan.
3. Run `node runtime/tools/build-skills-manifest.mjs` and `node runtime/tools/build-skills-guide.mjs`.
4. Run `ow-skill-security-check` before the skill is treated as trusted.
5. Leave it off `registries/routing.json` until `PromoteSkill` would accept a real trial file.

## Output

The new files, a manifest that lists them, and a security note.

## Verify

`check-skills.js` and `check-schemas.mjs` exit 0.

## Skip

Do not add a pack. Do not run `skills add --all`. Do not promote from a fixture whose status is `not_run`.

---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-onboard
description: Open a product repo onto Orchestra. Stamp adapters between markers, sync skills without deleting local edits, and write the project passport. Do not overwrite an unmarked AGENTS.md.
---

# Onboard

Use this when a repository should start using Orchestra.

## Steps

1. `kit/bootstrap.ps1` or `kit/bootstrap.sh` with the hosts the human named.
2. Adapters land between `<!-- orchestra:start -->` and `<!-- orchestra:end -->`. An unmarked file is left alone.
3. Skill sync is `runtime/tools/ow-sync.mjs`. Dry-run first. A local edit is a conflict, not a delete.
4. Copy `templates/passport/project.json` to `.orchestra/project.json` and fill the name, the owner lane, and `owned_files`.
5. Run `kit/orchestra-doctor.ps1`. Vault and heartbeat run only with `-Private`.

## Skip

Do not copy this method repo's skills into the product by hand. Do not prune user-installed skills. Do not commit MCP keys.

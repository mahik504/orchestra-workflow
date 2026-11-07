# Moving to 4.5.0

`registries/versions.json` is the current version. Scripts read it.

## Skill names

Orchestra-owned skills are `ow-*`. The old id is in `registries/skill-aliases.json`. Upstream names (GSAP, Expo, Stitch, Superpowers, Ponytail, grilling) did not change.

If a host still has `orchestra-conductor`, the next sync copies `ow-conductor` beside it. The old folder is not deleted, because sync does not prune. Remove the old folder when you are done with it.

## Sync

`kit/sync-ides.ps1 -Sync` used to delete the destination and copy over it. 4.5.0 hashes the tree. A local edit is a conflict. Pass `-BackupAndReplace` only when you want the canonical copy and a backup under `%LOCALAPPDATA%\orchestra\backups`.

## Adapters

`node runtime/tools/render-adapters.mjs` rewrites the routing table from `registries/routing.json`. Codex, OpenCode, and Kimi have adapters under `kit/`.

## What stayed

The eight stages, the dry-run human gate for a missing browser, and the commit-msg hook that strips `Co-authored-by: Cursor` are still there. `c1-premium-landing` stays disabled until a paired record exists.

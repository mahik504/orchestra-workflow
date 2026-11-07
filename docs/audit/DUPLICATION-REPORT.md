# Duplication report

Observed at `60cf9e7`. This lists copies that disagree or that no route loads. It does not list every repeated sentence in adapters. The routing block is supposed to match, and `check-routing.js` says the five copies match.

## Installers

`kit/sync-ides.ps1` and `kit/bootstrap.ps1` both delete a destination directory and copy the skill. The shell half of bootstrap does not exist. A later sync has to share one implementation (`runtime/tools/ow-sync.mjs`) so the two scripts cannot drift again.

## Skills that sit beside a route they must not join

| Directory | Why it is extra |
| --- | --- |
| `skills/creative-3d-web` (untracked) | The 3D route loads `r3f-shaders` only. |
| `skills/fullstack-mastery` (untracked) | The SaaS route loads `backend-architecture` and `test-driven-development`. |

`templates/shelf/` holds graphify, extra GSAP plugins, and Obsidian helpers. No routing row loads that shelf. `kit/shelf.ps1` exists to copy it.

## Docs

`docs/AGENT_PORTABILITY.md` overlaps `kit/HOST_CAPABILITIES.md` and is pinned to the 3.3 wording. Portability moves into `docs/hosts.md`. The 3.3 file gets a historical banner instead of a silent rewrite of every old sentence.

## Protocols

`protocols/` has 17 files. Several restate the same loop (classify, lab, implement, verify). 4.5 merges the overlapping ones into the conductor and leaves a short pointer file where a link would break.

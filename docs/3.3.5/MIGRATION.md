# Migration to 3.3.5

From 3.3.1 and the 3.3.2 prerelease.

| Change | What to do |
| --- | --- |
| Version pin | `VERSION` is 3.3.5 |
| Skills | Open set in `skills/`. `skills-local/` is not committed |
| Removed from the source set | `superpowers-planning`, `unlazy` |
| Contract | `AGENTS.md` |
| Rollback | Local tag `pre-3.3.5`, or `git checkout v3.3.1` after that tag exists |

Phase reports from 3.3.2 live under `docs/history/3.3.2/` when they were moved. The architecture and audit notes stay in `docs/3.3.2/`.

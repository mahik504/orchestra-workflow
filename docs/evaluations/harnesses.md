# Harness notes

Checked 2026-11-08. A row here is a reading, not a score. Orchestra stays the conductor.

| Project | What we kept | What we left |
| --- | --- | --- |
| Superpowers | Vendored skills already in `skills/`, with `LICENSE.superpowers`. | The plugin. It injects a second plan at session start. |
| Matt Pocock skills | `grilling`, `domain-modeling`, `to-tickets`, `constraint-driven-development`. | The rest of that pack, and any setup command that rewrites the host. |
| Ponytail | `ponytail`, `ponytail-review`, `ponytail-audit` on engineering diffs. | Using it to restyle a Design Lab surface. |
| Spec Kit | The idea that a task names acceptance before code. `ow-forge` already writes that line. | A second spec framework. |
| GSD | Progress you can point at. The receipt and `orchestra compile` are that list. | Its conductor. |
| BMAD | Separate reviewer roles. `agents/ow-reviewer.md` and `agents/ow-simplifier.md` are the two passes. | A committee that plans beside Orchestra. |
| ECC, gstack, Agency Agents | Nothing copied. | The packs. They open another plan. |
| Karpathy-style guidelines | Short context, one next action, no fake certainty. The fast path and `docs/CONTEXT-BUDGET.md` hold that. | A new skill named after a person. |

Commit SHAs for the skill repos we copied from are in `upstream.md`. The table below is the default branch we opened on 2026-11-08. None of these repositories were installed as a second conductor.

| Repository | Commit | Date | Decision |
| --- | --- | --- | --- |
| `obra/superpowers` | `22f0b463e7a5e579904dd6211a520321035526f7` | 2026-11-03 | Skills already vendored. Plugin stays off. |
| `mattpocock/skills` | `4ae2317ddb22b046e26455e4f9c9b3aaaf62212d` | 2026-10-31 | Four skills already in the tree. |
| `addyosmani/agent-skills` | `301b02b12bff98e4688e6984af6556802322bbde` | 2026-11-01 | `constraint-driven-development` already in the tree. |
| `github/spec-kit` | `c43849408ac0de2ca6a52f14eaae71449776057b` | 2026-11-07 | Acceptance line stays in `ow-forge`. |
| `bmad-code-org/BMAD-METHOD` | `e734d3ac3c0f4ea31cc2c2bbbe79e86e54f2cae6` | 2026-11-06 | Reviewer roles are the files in `agents/`. |
| `garrytan/gstack` | `1a29fbc3444f858c72ec001abdb4f0f3bebfa76c` | 2026-11-07 | Not installed. |
| `affaan-m/everything-claude-code` | `d2ef2680594f9a17526a6e1a2e18b3d242c5d17b` | 2026-11-06 | Not installed. |
| `open-gsd/gsd-core` | `4226711512fb4035fcbb7d17ae64ac5419d8396a` | fetched 2026-11-12 | Not installed. The archived `glittercowboy/get-shit-done` pin is not the active repo. Progress stays on the receipt and `orchestra compile`. |
| `EveryInc/compound-engineering-plugin` | `d8f8e03237bded06e77bcb3930b43ef841e7eeda` | fetched 2026-11-12 | Not installed. Read as a method source. |
| `alextnetto/res` | `1d05d9c4a9f96d3cd94e9fa49609f98d7e4b4563` | fetched 2026-11-12 | Not installed. Read for the shape of a product study. |
| `msitarzewski/agency-agents` | `b859c515d813efcfc4755d3471dd9e7686c7706d` | 2026-11-06 | Not installed. |
| `wshobson/agents` | `4bde3168b43c0655c4d00ea2726e631719908a6d` | 2026-11-05 | Not installed. |
| `openai/codex` | `33b3f80a0935e51a49cbbda7d6473ef757d4241b` | 2026-11-07 | Host adapter only. `kit/codex/AGENTS.md`. |
| `anomalyco/opencode` | `c33eaa170ed5093c08f1684862a695aae036dcda` | 2026-11-07 | Host adapter only. `kit/opencode/AGENTS.md`. |
| `crewAIInc/crewAI` | `73f48320dbe6e2dc7a9bc54fefc6a13637b9485c` | 2026-11-07 | Not installed. |
| `ruvnet/ruflo` | `3fe9d70ba7d08ff12233003a9342037ab80af3ce` | 2026-11-06 | Not installed. |

Cockpit has its own note in `cockpit.md`. No single Cockpit repository was pinned. Hermes and a generic "claude-code best practice" repo were not opened. Karpathy's notes are the short-context rules already in `AGENTS.md` and `docs/CONTEXT-BUDGET.md`. There is no skill named after a person.

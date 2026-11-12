# Research to engineering map

Recorded 2026-11-12. A row becomes a code change only when the source column is VERIFIED and the change has a check. REPORTED rows are reading notes.

## Change the harness

| Evidence | Status | Engineering consequence | Where |
| --- | --- | --- | --- |
| Generated context files did not raise resolution and raised cost about 20 percent. Developer files were about +2.4 percent and not significant on the page’s test | VERIFIED, arXiv:2602.11988 | Do not auto-write `AGENTS.md` overviews. Keep the standing file short. Point conventions at the README when they are already there. Shrink or replace `docs/CONDUCTOR.md` (19,467 bytes, titled 4.6.0) | `AGENTS.md`, `docs/CONDUCTOR.md`, `skills/ow-conductor` |
| Curated skills helped on tasks that were kept because the arms separated. Packs of at most three modules beat large bundles | VERIFIED, arXiv:2602.12670 | One route loads a few skills. Do not load the catalog. Do not add a bundle to “cover” a discipline | `registries/routing.json`, `compile.go` |
| Self-generated skill packs scored below the no-skill arm, and one apparent win named the verifier | VERIFIED, same paper | The agent that writes a skill does not put it on a Load line. `PromoteSkill` stays the gate | `runtime/internal/benchmark/promote.go`, `skills/ow-author` |
| Mean skill effect on SWE-Skills-Bench was +0.34 points. 40 of 98 skills flipped sign across harnesses. Harmful skills overrode paths | VERIFIED, arXiv:2603.15401 current HTML | Skill text says the task’s paths and acceptance win. A skill is re-checked per host, not assumed portable | Every `ow-*` procedure that gets written |
| Requirements in that benchmark were written from the skill | VERIFIED, method section | Orchestra’s own eval tasks must be written from the user job, not from the skill’s outline, or the trial flatters the skill | `docs/audit/EVAL-5.5.md` |
| `PromoteSkill` already encodes three trials, no negative success change, cost at most 1.10 | VERIFIED in this repo | Leave the function. Do not lower the bar to ship 5.5.0 | `promote.go` |

## Do not turn into a requirement

| Idea | Why not |
| --- | --- |
| “Skills add 16.6 points” | Filtered task set, curated skills, terminal containers, and a self-generated arm that went backwards |
| “Skills do nothing, delete them” | The same papers show a few skills moving a missing procedure by more than 15 points, and a curated pack helping hard tasks |
| “Delete AGENTS.md” | The human-file result was small and not significant. The cost result argues for a short file, not for no file |
| SWE-bench, Terminal-Bench, or OSWorld percentages as Orchestra goals | Different models, scaffolds, and years. Several were not re-opened here |
| A Minecraft-style ever-growing skill library | Voyager’s setting is not this repo, and self-written skills failed in SkillsBench |
| MemGPT or a new vector brain | The paper is document QA. This repo already has a vault and a dead SQLite stub. Another store is a second memory system |
| An AutoGen or MetaGPT committee | Not a repository-repair result. Orchestra already refused a second conductor |
| Hermes-style automatic skill authoring | Contradicts the self-generated arm that was opened |
| GSD’s “200k token” executor as a constant | Their README’s number, not a measurement |

## Map onto the two speeds

**Fast path.** A named file, a bug, or one control. L0. No conductor document. No new skill. This is the AGENTS.md paper applied to a small fix: extra context has a cost, and the overview is not what fixes the bug.

**Full path.** One route, at most the skills that route names, task text wins, verify on the running app or the named check, then the two reviews. A skill joins the route only when `PromoteSkill` would return true on a saved record. Until then it can exist as an opt-in file and must stay off the default Load line if the trial is missing. The 5.0 routes were added with `not_run` fixtures. They stay opt-in in the sense that no claim is enabled. Removing them from the table is a product choice, not a paper result. The paper result is: do not call them a win.

## Checks that would make a row VERIFIED for a future code change

- A saved trial file with three trials, the same model, and a cost ratio.
- A compile fixture where a miss sets confidence below 1.
- A byte budget: `docs/CONDUCTOR.md` either shrinks under a stated cap or stops being the procedure the stub loads.
- A grep gate: public `skills/` contains no Anthropic “may not distribute” license. That item is legal, not from these papers.

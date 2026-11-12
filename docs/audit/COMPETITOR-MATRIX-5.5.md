# Competitor matrix 5.5

Recorded 2026-11-12. Cells are judgments about files this session opened. They are not task scores. Stars are API `stargazers_count` the same day, and they do not fill a cell. `Unknown` means the relevant file was not opened. Detail and URLs are in `COMPETITIVE-ANALYSIS-5.5.md`.

Orchestra column is this repo at `24ad53c`.

## Popularity, separate from the matrix

| Project | stargazers_count | Pushed | License |
| --- | ---: | --- | --- |
| mattpocock/skills | 322806 | 2026-11-12 | MIT |
| obra/superpowers | 316662 | 2026-11-03 | MIT |
| affaan-m/everything-claude-code (API full_name `affaan-m/ECC`) | 301858 | 2026-11-09 | MIT |
| NousResearch/hermes-agent | 262916 | 2026-11-12 | MIT |
| garrytan/gstack | 141360 | 2026-11-12 | MIT |
| OpenHands/OpenHands | 93332 | 2026-11-12 | MIT |
| cline/cline | 71560 | 2026-11-12 | Apache-2.0 |
| Aider-AI/aider | 50089 | 2026-05-22 | Apache-2.0 |
| EveryInc/compound-engineering-plugin | 26050 | 2026-11-12 | MIT |
| SWE-agent/SWE-agent | 20786 | 2026-11-09 | MIT |
| open-gsd/gsd-core | 11485 | 2026-11-12 | MIT |
| alextnetto/res | 294 | 2026-10-17 | MIT |

These counts are not quality and not adoption of Orchestra.

## Matrix

Rows are capabilities. Columns are the projects whose files were actually opened far enough to judge. OpenHands, Aider, Cline, and SWE-agent are metadata-only and stay Unknown on every capability row except “is a coding agent,” which their descriptions support.

| Capability | Superpowers | ECC | gstack | Pocock skills | GSD Core | Compound | RES | Hermes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Skills as procedures | Stronger. Named loop and 15 skill dirs. Bodies not all read | Unknown. `skills/` exists. README not opened | Stronger on role commands. Bodies not read | Stronger on spec and diagnosis names. Bodies not diffed | Unknown. `skills/` exists. README describes phases, not each skill | Comparable names (`ce-plan`, `ce-debug`, `ce-code-review`). Bodies not read | One skill name. Body Unknown | README describes skills created from experience. Bodies Unknown |
| Session-start control | Stronger, and it is the part Orchestra refuses | Unknown | Stronger. Team mode checks for updates at session start | Unknown | Installer injects its own commands. README says not to copy files by hand | Unknown | Unknown | Install scripts. README describes its own CLI, not a plugin inside Cursor |
| Routing | Missing as a keyword table. README says skills trigger automatically | Unknown | Slash commands, not a scored route table | Unknown | Phase loop, not a 36-row table | Unknown | Unknown | User runs `hermes model` and slash commands |
| Verification before done | Comparable. They have `verification-before-completion`, which this repo vendors | Unknown | README says `/qa` opens a browser and `/review` finds bugs that pass CI. Not run | `tdd` is in the engineering list | README verify step before ship. Not run | `ce-test-browser` name only | Unknown | Unknown |
| Memory | Missing in the files opened | Unknown. Description mentions memory | `USING_GBRAIN_WITH_GSTACK.md` exists. Not opened | Unknown | `STATE.md` and `CONTEXT.md` named in the README | Unknown | README description says a knowledge base. Not opened | README describes FTS5 session search and a user model. State modules exist and were not read |
| Fresh context for workers | `subagent-driven-development` is named. Body not read | Unknown | Unknown | Unknown | README claims a clean context per executor and states 200k tokens. Their number, not measured here | Unknown | Unknown | README says isolated subagents. Not measured |
| Host coverage | Stronger list in the README: Claude, Cursor, Codex, Gemini, Copilot, OpenCode, Hermes, and others | Root has many host folders. Behavior Unknown | README centers Claude Code and says other agents work. Extent Unknown | `.claude-plugin` present. Other hosts Unknown | README lists Claude, Codex, Cursor, Copilot, Windsurf, and others | Several `.*-plugin` folders | `.claude-plugin` only in the listing | Its own CLI and messaging gateway. Not an IDE adapter |
| Security story | README telemetry opt-out. Skill security Unknown | `SECURITY.md` and `the-security-guide.md` exist. Not opened | README says setup refuses an old Bun because compiled tools could read `.env`. `SECURITY.md` not opened | Unknown | `SECURITY.md` exists. Not opened | `SECURITY.md` exists. Not opened | `SECURITY.md` exists. Not opened | README links a security doc. Install is `curl \| bash` and `iex (irm ...)` |
| Token discipline | Unknown. Bootstrap loads a skill every session, which the AGENTS.md paper suggests can raise cost | Description says performance optimization. No measurement opened | Unknown | Unknown | README’s reason for existing is context rot. The explanation doc was not opened | Unknown | Unknown | README has `/compress` and `/usage`. Effect Unknown |
| Evals | README points at `superpowers-evals`. Not opened | `tests/` exists. Not opened | Root has `benchmark/` and `eval-plan`. Not opened | Unknown | `tests/` and a test workflow badge in the README. Badge not treated as a result | `tests/` exists. Not opened | `tests/` exists. Not opened | `evals/` and `tests/` exist. Not opened |
| Install that preserves personal config | Unknown. Plugin install is per host | `install.ps1` not opened | Weaker for this goal. README tells the agent to edit `CLAUDE.md` and clone into `~/.claude/skills` | `setup-matt-pocock-skills` exists. Body Unknown | Installer is required. What it overwrites was not read | Unknown | Unknown | Installer is a remote script. What it overwrites was not read |
| Non-claim when a trial did not run | Unknown | Unknown | README points at a why-page for productivity numbers. Page not opened, numbers not used | Unknown | README states a 200k context and that the system makes Claude reliable. Those sentences are claims, not a trial against Orchestra | Unknown | Unknown | README says it is the only agent with a built-in learning loop. That exclusivity was not checked |

## Orchestra against that table

| Capability | Orchestra |
| --- | --- |
| Skills as procedures | Weaker on the 36 `stub: true` `ow-*` files. Comparable on the vendored Superpowers and Pocock files that are present. The vendored bytes were not diffed to upstream |
| Session-start control | Weaker on purpose. No second bootstrap. `skip orchestra` is the off switch |
| Routing | Present, keyword lengths, false confidence on a miss |
| Verification | Present as a contract and as `ow-qa`. `verify` can exit 0 on `NOT_APPLICABLE`. Resource memory can store a skip as success |
| Memory | Vault plus `AdmitDurable`, not called from the pipeline. SQLite stub is dead. Weaker than Hermes’s description, and safer than an unreviewed self-written skill |
| Fresh context | Described as a worker capsule. Not enforced on the host |
| Host coverage | Comparable for Cursor, Claude, Antigravity, Codex, jcode. Codex is missing from the Go parity map |
| Security story | Hook denylist and SSRF checks on one adapter. False PASS lines and proprietary skills in tree |
| Token discipline | Measured file sizes. No dollar measurement. Conductor doc is 19,467 bytes and stale |
| Evals | Harness exists. Fixtures `not_run`. Claims off. This is more honest than an unmeasured superiority sentence, and it is not a result |
| Install | Bootstrap is described as leaving a non-Orchestra file alone. Not re-tested in this pass |
| Non-claim | Stronger. `compare.go` is disabled. `EVAL-5.0.md` says the trial did not run |

## What the later competitor read adds

[Competitors](a471d9f0-4685-4e58-89db-990a6993d9ae) scored other projects against `AGENTS.md` and `docs/ROUTING.md` only, so their table marks Orchestra memory and host lists as missing. The runtime has a vault, host adapters, and `PromoteSkill`. Those cells above stand. Mechanisms they opened that this session did not, all REPORTED:

| Project | Addition |
| --- | --- |
| ECC | Memory vault and instinct caps. README catalog counts were not re-counted. Verification and security move from Unknown to REPORTED-stronger on paper, still not a trial |
| gstack | Browser QA and an untrusted-page rule in `ARCHITECTURE.md`. Security moves from Unknown to REPORTED. The logged-in browser is a reject |
| GSD | `.planning/` artifacts, “absent means enabled,” a worktree hook that fails open. Stronger on file state, weaker as a default to copy |
| Compound | `docs/solutions/` for the next plan, and a stop when a skill reference is missing |
| RES | A measured token aside and a sandbox ledger. Not a general harness |
| Karpathy gist | A wiki schema. Not a skills repo |

## Bottom line

Superpowers, gstack, Pocock’s engineering list, and GSD’s README show fuller procedures or a harder session loop than Orchestra’s stub skills. ECC, Compound, OpenHands, Aider, Cline, and SWE-agent were not opened deeply enough to rank. Hermes shows a memory product Orchestra should not copy wholesale. No cell above is a measured task win.

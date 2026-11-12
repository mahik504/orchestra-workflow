# Competitive analysis 5.5

Recorded 2026-11-12. This session queried the GitHub API and opened README or directory listings. No competitor was installed. No task was run on two harnesses. Stars are `stargazers_count` from the API at query time. They are popularity, not engineering quality, and not a measured win.

`docs/evaluations/harnesses.md` is an earlier reading dated 2026-11-08 with pinned commits. It is not re-used as a current file hash. Where this document and that note disagree, this document is the one with a 2026-11-12 API call.

[Competitors](a471d9f0-4685-4e58-89db-990a6993d9ae) returned after the first draft. The body of this file is this session’s API listings and the READMEs this session opened. The addendum at the end is their report. A claim in the addendum is REPORTED unless this session had already opened the same file. Their matrix treated `AGENTS.md` and `docs/ROUTING.md` as the whole of Orchestra. The runtime findings in `BASELINE-5.5.md` override that narrower baseline. Star counts in their pass differ from the table below by a few stars on the same day. Both are API snapshots. Neither is a quality score.

## How to read a cell

- **Stronger** means the other project’s opened files show a mechanism Orchestra’s opened files do not, or a deeper procedure for the same job.
- **Comparable** means both have the mechanism, with differences that were not measured.
- **Weaker** means Orchestra’s opened files show a control the other project’s opened files do not emphasize.
- **Missing** means the mechanism was not found in the files this session opened.
- **Unknown** means the file was not opened.

A star count never decides a cell.

## Orchestra, the baseline being compared

Control plane. Keyword route table of 36 jobs. 109 indexed skills, 36 of the `ow-*` skills flagged `stub: true`. Host adapters for Cursor, Claude Code, Antigravity, Codex, jcode, OpenCode, and Kimi, generated from one route table. Eval claims disabled. Promotion requires three trials and at most a 10 percent cost increase. Private vault is the memory store. SQLite store in `memory.go` has no caller.

## Superpowers

Repository: `obra/superpowers`. API: default branch `main`, pushed 2026-11-03, license MIT, `stargazers_count` 316662.

Opened: README, and the `skills/` directory listing. Skill directories listed: `brainstorming`, `diagnosing-superpowers`, `dispatching-parallel-agents`, `executing-plans`, `finishing-a-development-branch`, `receiving-code-review`, `requesting-code-review`, `subagent-driven-development`, `systematic-debugging`, `test-driven-development`, `using-git-worktrees`, `using-superpowers`, `verification-before-completion`, `writing-plans`, `writing-skills`.

The README describes a session-start bootstrap. On Antigravity the `using-superpowers` description prompts the model to load that skill at the start of a conversation. The Pi package injects the bootstrap at session start and again after compaction. The basic workflow is brainstorming, worktrees, writing plans, subagent-driven development or executing plans, TDD, code review, and finishing a branch. The README says the agent checks for relevant skills before any task, and that these are mandatory workflows. Skill-behavior tests are described as living in a separate `superpowers-evals` repo. Plugin tests are in `tests/`. Those test files were not opened.

Orchestra already vendors `systematic-debugging`, `test-driven-development`, `verification-before-completion`, `using-git-worktrees`, `receiving-code-review`, and `writing-skills`, with `LICENSE.superpowers` named in `THIRD_PARTY_NOTICES.md`. The changelog says the Superpowers plugin stays off because it injects a second plan.

What they do that this tree does not: a single mandatory loop from spec to branch finish, a bootstrap that fires without a route keyword, and a skill whose job is to diagnose a bad Superpowers session from the transcript. What not to copy: the bootstrap. Orchestra’s contract is one conductor. A second session-start plan fights `skip orchestra` and the fast path. Also not copied: the optional logo fetch described as telemetry in the README. Orchestra should not add a phone-home image.

Individual skill bodies other than the names were not read. Whether their TDD text is better than the vendored copy is UNKNOWN until a diff.

## Everything Claude Code

Repository requested: `affaan-m/everything-claude-code`. API `full_name` field returned `affaan-m/ECC`. The contents API on `affaan-m/everything-claude-code` succeeded. Default branch `main`, pushed 2026-11-09, MIT, `stargazers_count` 301858. Description field: “The agent harness performance optimization system. Skills, instincts, memory, security, and research-first development…”

Opened: the root directory listing, not the README body. The root contains `skills/`, `agents/`, `commands/`, `hooks/`, `rules/`, `install.ps1`, `install.sh`, `mcp-configs`, `tests/`, `the-security-guide.md`, `SECURITY.md`, `AGENTS.md`, `CLAUDE.md`, and host folders for many editors. That is a wide pack. Orchestra’s changelog says this pack was not installed.

What the listing shows that Orchestra does not: an install script pair, a security guide filename, and host folders beyond the set Orchestra generates. What the listing does not show: whether any of it beats a raw host. README prose is UNKNOWN. Copying the pack would be a second conductor, which this repo has already rejected.

## gstack

Repository: `garrytan/gstack`. API: `main`, pushed 2026-11-12, MIT, `stargazers_count` 141360.

Opened: README through the install section, and the root listing. The README says gstack is a virtual engineering team of slash commands: CEO, eng manager, designer, reviewer, QA in a real browser, security officer, release engineer. Install is a clone into `~/.claude/skills/gstack`, `./setup`, and a new section in `CLAUDE.md` that tells the model to use gstack’s browse skill. Team mode commits that requirement into the project and checks for updates on session start. The README points at `docs/why-gstack.md` for productivity numbers. That page was not opened. Those numbers are not repeated here.

The root listing includes `qa/`, `review/`, `ship/`, `plan-eng-review/`, `design-review/`, `cso/`, `browse/`, `BROWSER.md`, `ARCHITECTURE.md`, `SECURITY.md`, and many other role directories. Bodies were not read.

What they do that Orchestra’s opened files do not: a browser QA command described as opening a real browser, a security-review command described as OWASP and STRIDE, and an install that is one paste. Orchestra has `ow-qa` (procedure opened, Playwright) and `ow-security-reviewer` (file exists, body not graded). gstack’s install rewrites `CLAUDE.md` and auto-updates. Orchestra’s contract leaves a personal file alone when it does not already look like Orchestra. That install is the part to reject.

## Matt Pocock

Repository: `mattpocock/skills`. API: `main`, pushed 2026-11-12, MIT, `stargazers_count` 322806. Description: “Skills for Real Engineers. Straight from my .agents directory.”

Opened: root listing and `skills/engineering/` names. Engineering names: `ask-matt`, `code-review`, `codebase-design`, `diagnosing-bugs`, `domain-modeling`, `grill-with-docs`, `implement`, `implement-spec`, `improve-codebase-architecture`, `pr`, `prototype`, `research`, `retro`, `setup-matt-pocock-skills`, `tdd`, `to-spec`, `to-tickets`, `triage`, `wayfinder`, `wizard`. Other top groups: `deprecated`, `in-progress`, `misc`, `productivity`. README body was not opened.

Local tree contains `skills/grilling`, `skills/domain-modeling`, and `skills/to-tickets`. Upstream’s current engineering name is `grill-with-docs`, not `grilling`. The files were not diffed. Orchestra may be behind the upstream names. That is a maintenance fact, not a quality score.

What the listing shows: a spec path (`to-spec`, `implement-spec`, `tdd`) and a diagnosis skill. Orchestra’s engineering table names grilling, domain modeling, to-tickets, TDD, and systematic debugging. It does not name `to-spec`, `wayfinder`, or `codebase-design`. Adding all of them would grow the skill count. The papers below say most added skills do not move the mean. Adopt a name only after a paired trial, and only when the local copy is actually missing the procedure.

## Karpathy

`karpathy/autoresearch` exists. API contents: `README.md`, `prepare.py`, `train.py`, `program.md`, `analysis.ipynb`. That is a training experiment, not a coding-agent harness. README body was not opened. No repository of “Karpathy agent skills” was found in this pass, and none is invented. The short-context rule already in `docs/CONTEXT-BUDGET.md` is Orchestra’s own file. Attributing it to a person does not create a skill.

## Open GSD Core

Repository: `open-gsd/gsd-core`. API: default branch `next`, pushed 2026-11-12, MIT, `stargazers_count` 11485.

Opened: README. It describes a phase loop: discuss, plan, execute, verify, ship. Heavy work runs in fresh-context subagents. The README says each executor starts with a clean 200k-token context. That figure is their claim, not a measurement of Orchestra. Artifacts named: `STATE.md` and `CONTEXT.md`. Install is `npx @opengsd/gsd-core@latest`, which prompts for a runtime. The README says not to copy `agents/` or `commands/` by hand.

Root listing includes `src/`, `tests/`, `skills/`, `commands/`, `hooks/`, `capabilities/`. Those directories were not listed further.

Orchestra already has `.orchestra/STATE.md` in the README loop, a worker capsule that is supposed to exclude unrelated skills, and `ow-wave` for phased delivery. GSD’s fresh-context executor is the idea Orchestra’s capsule states and does not enforce in the host. What not to copy: a second installer and a second command namespace (`/gsd-new-project`). The 200k figure is not a requirement.

## Compound Engineering

Repository: `EveryInc/compound-engineering-plugin`. API: `main`, pushed 2026-11-12, MIT, `stargazers_count` 26050.

Opened: root listing and `skills/` names. Skills include `ce-plan`, `ce-work`, `ce-debug`, `ce-code-review`, `ce-simplify-code`, `ce-test-browser`, `ce-compound`, `ce-handoff`, `ce-worktree`, `ce-brainstorm`, and many more `ce-*` commands, plus `lfg` and `wtf`. README body was not opened. Host folders exist for Claude, Cursor, Codex, Gemini, and others.

The names overlap Orchestra’s loop (plan, work, review, simplify, browser test, handoff, worktree). The pack is another conductor if installed whole. `docs/evaluations/harnesses.md` said it was read and not installed. This session did not re-read the skill bodies. Quality versus Orchestra is UNKNOWN. The duplication is visible from the names.

## RES

Repository: `alextnetto/res`. API: `main`, pushed 2026-10-17, MIT, `stargazers_count` 294. Description: an agent skill to reverse engineer a SaaS product through the browser into a knowledge base.

Opened: root listing and `skills/` names. One skill directory: `reverse-engineering-saas`. README body was not opened. Orchestra’s nearest routes are `ow-site-recon`, `ow-app-recon`, and `ow-reference-forensics`, with explicit skips against copying assets, copy, or a logo. RES is one deep skill. Orchestra is a route that forbids cloning the site. Those are different jobs. Adopting RES as a rebuild-the-product skill would fight the skip line. A study of their safety model was not done. UNKNOWN beyond the description and the single skill name.

## Additional systems

API metadata and, where noted, a README or a root listing. Not installed.

| Repo | API | What was opened | Relevance |
| --- | --- | --- | --- |
| `OpenHands/OpenHands` | main, pushed 2026-11-12, MIT, stars 93332 | metadata only | A coding-agent product with its own runner. README UNKNOWN. Not a skill pack to merge |
| `SWE-agent/SWE-agent` | main, pushed 2026-11-09, MIT, stars 20786 | metadata and description | Description says it takes a GitHub issue and tries to fix it, and that it can also be used for offensive cybersecurity. Orchestra’s recon rules already limit offensive tools to owned, client, CTF, or open-source targets. The offensive clause is a reason not to vendor the repo. README body UNKNOWN |
| `Aider-AI/aider` | main, pushed 2026-05-22, Apache-2.0, stars 50089 | metadata only | Terminal pair programmer. Last push is months behind the others. Internals UNKNOWN |
| `cline/cline` | main, pushed 2026-11-12, Apache-2.0, stars 71560 | metadata and description | IDE, SDK, or CLI agent. Internals UNKNOWN |
| `NousResearch/hermes-agent` | main, pushed 2026-11-12, MIT, stars 262916 | README and root listing | See below |
| `langchain-ai/langgraph` | contents were not requested after the metadata call failed to include it in the successful metadata batch | UNKNOWN | Not scored |
| `karpathy/autoresearch` | contents listing | training scripts | Not a harness |

Hermes README, opened: it describes a model command (`hermes model`), memory, skills created after complex tasks, FTS5 session search, subagents, and several terminal backends. Install on Linux is `curl | bash`. Install on Windows is `iex (irm https://hermes-agent.nousresearch.com/install.ps1)`. Orchestra’s hook denies `curl` piped to `sh` and does not match `iex`. The root listing is dominated by `hermes_state_*.py` modules (sessions, FTS, compression, WAL, and many others). Those modules were not read. The README’s learning loop is the opposite of SkillsBench’s self-generated-skill result, which is in the research review. Do not auto-write skills because Hermes does.

CrewAI, DSPy, Letta, Mem0, Pydantic AI, Mastra, and AutoGen were not opened in this pass. The MemGPT paper was opened as a paper, not as the Letta product. They stay UNKNOWN as products. The research review covers the papers that were opened.

## What Orchestra already does that these READMEs do not emphasize

- A machine-readable route with an explicit skip line, and a compiler that can return `ambiguous` instead of guessing on a tie.
- A promotion function that refuses a skill when the trial count is under three or the cost ratio is over 1.10, and a comparison suite that stays disabled.
- A written rule that a missing tool is reported and not installed.
- Sync that will not copy `creative-3d-web` or `fullstack-mastery` because they are off `host-stack.json`.
- A license boundary this audit found and the other READMEs do not fix for us: the Anthropic document skills.

## What neither side showed, in the files opened

- A published paired trial of Orchestra against Superpowers, ECC, gstack, or GSD on the same model and the same repo.
- A router that is allowed to override the host’s model picker.
- A memory store that both learns from verified failures and refuses to treat an unverified note as truth, and that has been measured. Orchestra has the refusal in `AdmitDurable` and does not call it from the pipeline. Hermes describes the learning and does not, in the README, describe a promotion gate.

## Addendum from the competitor reader

REPORTED. Not re-opened here except where an earlier section already quoted the same README.

- **Superpowers.** `skills/using-superpowers/SKILL.md` tells the model to invoke a skill before any response, including when the chance it applies is about 1 percent. `hooks/hooks.json` runs a session-start hook. That is the second conductor this repo refuses. Their skill-directory count in that report is 14. The contents listing in this session named 15 directories, including `diagnosing-superpowers`. Keep the name list from this session.
- **ECC.** The README they opened describes 314 skills, 71 agents, and 98 command shims, an instinct cap, and a Markdown memory vault at `.ecc/memory/` that the README calls unreviewed context rather than policy. Windows memory writes are named as open defects. Codex has no Claude-style hook parity. Copying the flat catalog is the mistake. Those counts were not re-listed here.
- **gstack.** `ARCHITECTURE.md` prefers a real logged-in browser, then a localhost Chromium daemon, and treats page text as untrusted. Windows support is described as minimal. Making that browser the default verifier is the part to reject. The architecture file was not opened in this session.
- **Pocock.** Install docs warn that a plugin plus skills.sh duplicates every skill. `ask-matt` is a user-invoked router. Skill bodies were not opened by that reader either.
- **Karpathy.** Besides `karpathy/autoresearch`, they opened gist `442a6bf555914893e9891c11519de94f` (4 Apr 2026): immutable sources, a markdown wiki, and a schema file the human and the model edit together. It is a memory pattern, not a coding harness. This session did not open the gist. There is still no Karpathy skills repository.
- **GSD.** State lives under `.planning/`, not only `STATE.md`. Feature flags use “absent means enabled.” The worktree hook they opened fails open when git cannot resolve the root. `docs/security/baseline.md` is still a draft. The README test badge points at `main` while the default branch is `next`. Do not copy the missing-config-means-on rule.
- **Compound Engineering.** README figure they reported: 38 skills, 14 hosts. `/ce-compound` writes `docs/solutions/` for the next plan. `/ce-work` stops if a bundled reference cannot be read. Running the full loop on every change is the part to reject.
- **RES.** One skill. They reported an idle cost near 130 tokens and about 4,600 when running, from the project’s own `claude plugin details` measurement, which they did not rerun. Safety levels L1–L4. It is a browser study with a ledger, not a default rebuild skill.
- **Also opened by that reader, metadata plus a README:** mini-swe-agent (the SWE-agent README tells readers to use it), Continue (README says the repo is unmaintained), LangGraph (a workflow library, not a coding-agent product). Official instruction pages they opened: Cursor rules and skills, Claude Code skills and memory, Codex `AGENTS.md` with a 32 KiB combined cap, Gemini CLI `GEMINI.md`, and Antigravity rules and skills docs. Those pages were not re-fetched here. The 32 KiB Codex cap is REPORTED.

## Reject list

Do not install Superpowers, ECC, gstack, GSD, Compound Engineering, or Hermes as a second conductor. Do not copy star counts into a README. Do not add a phone-home pixel. Do not add `curl | bash` or `iex (irm ...)` as Orchestra’s installer. Do not auto-author skills from a session. Do not treat SWE-agent’s offensive-cyber sentence as a feature.

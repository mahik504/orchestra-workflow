# Baseline 5.5

Recorded 2026-11-12 from the working tree of `mahik504/orchestra-workflow`. This is an inspection record. It is not a release, and it does not change `VERSION`.

## Re-measure before the 5.5 edits

`rtk go test -count=1 ./... -timeout 600s` from `runtime/` exited 0. The runner printed `Go test: 776 passed in 26 packages`.

`node runtime/tools/check-route-agreement.js` printed `routed 0` because the conductor routing block is a pointer, so the script checked nothing.

`skills/docx`, `skills/pdf`, `skills/pptx`, `skills/xlsx`, `skills/fullstack-mastery`, and `skills/creative-3d-web` are on disk and gitignored. `git ls-files` does not list them. They are not a public-tree license blocker. `runtime/orchestra.exe` is a local ignored binary, not a file CI ships.

`docs/CONTEXT-BUDGET.md` had an uncommitted one-line change at the start of this inspection. The measure script was re-run after the edits. `node runtime/tools/measure-context.mjs --check` exited 0 and printed `AGENTS.md 2047`, `descriptions 25614`, award-style site 7 skills 56621 bytes, full-stack app 4 skills 18698 bytes, and mobile app 5 skills 39661 bytes.

## Checks after the 5.5 edits

`go test -count=1 -timeout 600s ./...` from `runtime/` exited 0. The runner printed `Go test: 795 passed in 26 packages`.

`node --test runtime/tools/*.test.mjs` exited 0 with 18 passed. `node --test kit/hooks/*.test.mjs` exited 0 with 13 passed. `check-versions.js` printed `current: 5.5.0` and `problems: 0`. `check-schemas.mjs` printed `problems: 0`. `hygiene.js` printed `failures: 0`. `check-skills.js` printed `problems: 0`. `check-route-agreement.js` printed `route agreement skills 36, routed 47, off graph 11, problems 0`. `check-config-lint.mjs` printed `config lint: 0 problems`. `check-mermaid.mjs` printed `problems: 0`. `check-registry.js` printed `problems: 0`.

The paired eval in `docs/EVAL-5.5.md` was not run. The claim stays none.

Investigator: this Cursor session. Reports that were folded in after a source check: [Repo forensic](97661c75-d7de-4386-8103-d2de755c6fed), [Security and memory](f5ad95c7-969d-4868-9749-f6d66d2272a7), [Papers](73faad4f-9ccd-4a2f-be4b-b583487ab68d), [Skills and agents](9d967d24-6ae9-48bf-873b-444d3877a288), [Competitors](a471d9f0-4685-4e58-89db-990a6993d9ae). The skills and competitor reports arrived after the first draft. Their claims are used below only where this session reopened the file, or they are marked REPORTED.

Status words: VERIFIED means the file, command, or page was opened here. REPORTED means a project or an API said it and this session did not re-run it. INFERRED means a reading of that evidence. NOT_RUN means the check was not executed here. IMPLEMENTED means the source does the thing. A green Actions conclusion is not a local re-run.

## Answers

1. **Current state.** Orchestra 5.0.0 on `main` at `24ad53c3ca4e4098491da7ec7f29683670c09cf3`. The annotated tag object `v5.0.0` is `c459a0b`; the commit it points at is the same `24ad53c`. GitHub’s latest release page is Orchestra 5.0.0, published 2026-11-12. There is no 5.5.0 pin.

2. **Implemented.** A keyword router (`orchestra compile`), a design-pipeline engine, host adapters, a skill tree, registries, hooks, a promotion function that stays closed, and a resource-memory JSON file. The vault is a documented private workspace, not a database in this repo.

3. **Verified here.** Git pin, file sizes, route and skill counts, the false-PASS branches, the unused SQLite stub, the route-table gap for three skills the changelog names, the Anthropic license text in `skills/docx`, and the current HTML of three arXiv papers. GitHub Actions reports CI run `43197120648` and Hygiene run `43197120670` as `success` on this commit. The logs were not downloaded.

4. **Implemented, not re-tested here.** `go test`, Node tool tests, hook tests, and `orchestra doctor` / `verify`. The Actions success is the remote run. Local execution is NOT_RUN. `PromoteSkill` is unit-tested in the tree; no live trial file was produced.

5. **Experimental.** `orchestra mcp-serve`, per `CHANGELOG.md` and `docs/audit/EVAL-5.0.md`. Comparison tasks in `compare.go` stay disabled. Eval fixtures stay `not_run`.

6. **Advantages, from this tree.** One conductor. A written non-claim when a trial did not run. A promotion gate with a 10 percent cost cap. Disclosure levels L0–L4. Host sync that refuses a symlink and uses `registries/host-stack.json` rather than every folder. Hooks that refuse `skills add --all`, force-push, `curl|sh`, and `.env` reads.

7. **Weaknesses.** Several commands can print PASS or success for a condition they did not check. Thirty-six `ow-*` skills are marked `stub: true` and have no `## Procedure`. `docs/CONDUCTOR.md` still titles itself 4.6.0 and is 19,467 bytes. Routing is keyword length, and a miss keeps confidence `1`. Three skills named as 5.0 routes are not in the route table. Proprietary Anthropic document skills sit in the public tree and are absent from `THIRD_PARTY_NOTICES.md`.

8. **Where Superpowers is stronger.** Its README describes a mandatory brainstorm → plan → TDD → review loop, with skill bodies for that loop, and a session-start bootstrap. Orchestra vendors six of those skills and leaves the plugin off so a second plan does not start. The bootstrap and the full procedure depth are theirs. This is a reading of their README and skill directory listing, not a paired trial.

9. **Where ECC is stronger.** The repository at `affaan-m/everything-claude-code` contains host trees, hooks, install scripts, tests, and security guides. Orchestra did not install that pack. Whether those files produce better patches was not measured. File-tree comparison only.

10. **Where gstack is stronger.** Its README describes role commands that open a browser, review a branch, and ship a PR, plus a setup that rewrites `CLAUDE.md`. Orchestra’s two review agents and `ow-qa` are thinner. Copying gstack would add a second conductor and an auto-update into `~/.claude/skills`.

11. **Where Matt Pocock’s public skills are stronger.** The `engineering/` directory on `mattpocock/skills` lists procedures this repo only partly vendors (`tdd`, `to-spec`, `implement`, `diagnosing-bugs`, `codebase-design`, `wayfinder`, and others). Local copies present: `grilling`, `domain-modeling`, `to-tickets`. The rest were not diffed against upstream.

12. **Where other systems are stronger.** GSD Core’s README describes fresh-context executors and a discuss → plan → execute → verify → ship loop. Hermes Agent’s README describes a model switch, persistent memory, and skills written from experience. OpenHands, SWE-agent, Aider, and Cline are repository-level coding agents with their own runners. None of these were installed or scored against Orchestra.

13. **What Orchestra already does.** It refuses to claim a win without a saved trial. It keeps a fast path in prose for a named file. It separates capabilities, executors, memory, and the registry. It records `claim: none`.

14. **Missing.** A compile result that lowers confidence on a miss. Honest PASS lines. Procedure sections on the hot-path stubs. A route row or a changelog correction for `ow-browserize`, `ow-qa`, and `ow-author`. A paired eval that has actually run. Removal of the Anthropic document skills from the public tree.

15. **Token waste, measured as bytes.** `docs/CONDUCTOR.md` is 19,467 bytes and the conductor skill points at it. `AGENTS.md` is 2,049 bytes while `docs/CONTEXT-BUDGET.md` still says 1,847. This Cursor session also received `CLAUDE.md` (1,242 bytes) and `.cursorrules` (1,117 bytes). Skill descriptions for all canonical skills are recorded in the budget doc as 25,614 characters. That figure was not re-measured, because the measure script rewrites the doc.

16. **Latency.** No wall-clock timing was taken. Keyword compile is a local scan of 36 routes. The design pipeline, subagents, and `docs/CONDUCTOR.md` are the likely costs. Those are INFERRED, not timed.

17. **Unnecessary complexity.** Three overlapping stories: the eight design stages, the mission capsule, and the README loop. A SQLite `Store` with no callers beside a JSON resource ledger and a vault. Empty confidence math that reports `1` when nothing matched. A doctor quarantine count hardcoded to 1,598.

18. **Output quality.** Finish the hot-path procedures so they lose to the task text. Make PASS mean the check ran. Run the paired fixtures before adding skills. The skill papers disagree on the average effect, so a new skill is not assumed to help.

19. **Frontend.** Routes already name award-style sites, unslop, UI assembly, site recon, and one 3D scene. `ow-design-intelligence` is a stub. `ow-ui-assembly` has a `## Procedure` heading and was not graded. `creative-3d-web` is on disk and off the route table and off the sync list. Do not add another design kit.

20. **Backend.** The route loads `ow-backend-architecture` and `test-driven-development`. The backend skill is `stub: true`. `fullstack-mastery` is on disk and off the sync list, which matches the rule that a backend leaves it unloaded.

21. **Research and ML.** `ow-research`, `ow-research-paper`, and `ow-ml-experiment` are stubs. The papers say generated context and self-written skills are not free wins. A research workflow should verify citations and keep scores inside the paper that measured them.

22. **Add.** Honest exit and PASS behavior. A low confidence on a route miss. A fast-path flag from `orchestra compile` when the task names a file. License removal for the Anthropic document skills.

23. **Improve.** The 36 stub `ow-*` skills, starting with conductor, forge, backend, and vault. `docs/CONDUCTOR.md` should match 5.0.0 or stop being the long procedure. Version sentences in `ARCHITECTURE.md` and `WORKFLOW.md` that still say 4.6.0 and 4.5.0.

24. **Merge.** Do not merge `ow-focus` into `ow-forge`. `docs/audit/STUB-TRIAGE.md` says the route still names `ow-focus`. The eight-stage story and the README loop should be one description, not two engines.

25. **Remove from the public tree.** `skills/docx`, `skills/pdf`, `skills/pptx`, `skills/xlsx`, whose license text forbids distribution. The SQLite stub should be marked dead or deleted when a code change is authorized. This research pass does not delete them.

26. **Innovations worth a trial.** Fast-path compile. Task text wins over skill text. Promotion only from a paired record. Vault facts through `AdmitDurable`, capped at 8. Everything else in the assignment’s innovation list is already a feature here or a second system.

27. **First.** The false PASS lines, the proprietary skill trees, and the route-table lie about `ow-qa`, `ow-browserize`, and `ow-author`.

28. **Before a 5.5.0 tag.** Tests green on the release commit, version files telling one pin, comparisons still off unless a saved trial passes `PromoteSkill`, hot-path stubs finished or still flagged, the Anthropic document skills gone from the public tree, and the claim sentence still none.

No numeric score is given. A score would invent a rubric.

## Repository state

| Fact | Evidence | Status |
| --- | --- | --- |
| Branch `main` tracks `origin/main` | `git rev-parse`, `git status` | VERIFIED |
| HEAD `24ad53c3ca4e4098491da7ec7f29683670c09cf3` | `git log -1`, subject “Record that the 5.0 routes were not measured.”, author Mahaveer Singh Gehlot, 2026-11-12 15:34:30 +0530 | VERIFIED |
| `v5.0.0` tag object `c459a0b`, peeled commit `24ad53c` | `git rev-parse v5.0.0` and `git rev-parse v5.0.0^{commit}` | VERIFIED |
| `VERSION` is `5.0.0` | file contents | VERIFIED |
| Latest GitHub release is 5.0.0, published 2026-11-12T10:10:03Z | release page | VERIFIED |
| Open issues count | GitHub API `open_issues_count` = 0. `gh issue list` and `gh pr list` printed no rows | VERIFIED as the API field. A disabled issues UI was not separately checked |
| CI | Actions run 43197120648, workflow CI, conclusion success, same SHA, created 2026-11-12T10:04:44Z | VERIFIED via `gh run list`. Log body NOT opened |
| Hygiene | Actions run 43197120670, conclusion success, same SHA and time | VERIFIED the same way |
| Worktree | `git status --short` empty after `git checkout -- docs/CONTEXT-BUDGET.md` | VERIFIED at restore time. Later doc writes are this package |

Counts measured in this session:

- 36 jobs in `registries/routing.json`.
- 109 objects in `registries/skills.json` and 109 names in `registries/host-stack.json`.
- 115 directories under `skills/`, each with a `SKILL.md`.
- The six directories in the tree and not in either index: `creative-3d-web`, `fullstack-mastery`, `docx`, `pdf`, `pptx`, `xlsx`.
- 27 alias keys in `registries/skill-aliases.json`.
- 208 entries in `registries/resources.json` (a JSON array).
- Design graph `version` 5.0.0, 33 capabilities, 42 domains.
- 21 files in `agents/`, 22 files in `commands/`.
- 47 `ow-*` skill directories. 36 of them set `stub: true` and have no line `## Procedure`. The other 11 `ow-*` skills have a Procedure heading. `ow-qa` was opened and has Preconditions, Procedure, and a verify step.

`ow-sync.mjs` `canonicalIds` reads `host-stack.json`. The six extra directories are not in that list, so a sync that uses `canonicalIds` does not copy them. They are still in the public git tree.

## Execution lifecycle

What the files actually do:

1. The chat contract in `AGENTS.md` is forge, classify from `docs/ROUTING.md`, design gate, implement, verify, review.
2. `orchestra compile` in `runtime/internal/compile/compile.go` scores `registries/routing.json` by summing the lengths of matched keywords. A tie sets `ambiguous` and confidence 0.5. No hit returns route `unspecified` and leaves confidence at 1 (the struct default). Negation is a 48-character window after “not / never / do not”.
3. `orchestra plan` and `run` build `engine.NewDesignPipeline`. Stages are Discover, Classify, Research, Synthesize, Design System, Implement, Visual QA, Iterate. `ImplementStage` sets `BuildOutput` to `handoff prepared; product build was not run` and `BuildPassed` to false (`stage_implement.go` lines 372–373). VERIFIED by opening the file in the red-team pass.
4. `IterateStage` sets `FinalVerdict` to `NOT_APPLICABLE` when visual QA was skipped, then calls `recordIterateMemory(ctx, "PASSED", ...)`. `recordIterateMemory` stores `OutcomeSuccess` whenever the argument is `"PASSED"`. `Execute` then sets mission status to succeeded without reading the visual verdict (`pipeline.go` lines 249–254).
5. `orchestra verify` prints that `NOT_APPLICABLE` is not a pass and still returns exit 0 unless the verdict is FAIL, BLOCKED, or NOT_RUN, or `--strict` is set and warnings were recorded (`main.go` `verifyExitCode`).

Missing as code, present as prose: the fast path in `docs/CONDUCTOR.md` (“no brief, no receipt”) is not a branch in `compile.go`. Confidence and ambiguity exist only in the compile manifest.

## Findings

### F-01. Route miss keeps confidence 1

`compile.go` sets `Confidence: 1` and returns that manifest when `len(hits) == 0`. No compile test mentions `Confidence`. IMPLEMENTED. NOT executed.

### F-02. Quarantine PASS with the constant 1598

`runDoctor` defaults the path to `~/.gemini/config/skills_library`. `CheckQuarantineBoundary` errors when the path contains `skills_library`. The PASS branch runs when that error is set and the directory is missing, the status is nil, or the directory is isolated. The printed count starts at 1598. `main.go` lines 214–226. IMPLEMENTED. NOT executed.

### F-03. Host parity is complete when no host directory exists

`VerifyParity` starts both flags true and `continue`s when `os.Stat` says the directory is missing (`host_sync.go` lines 147–172). Doctor prints `[OK] N/N` (lines 264–266). `orchestra verify` prints `[PASS]` (lines 597–600). The sync command prints the 100 percent sentence when both flags stay true (lines 725–726). Codex is in `ow-sync.mjs` `hostDirs` and absent from `ResolveHostDirs`. IMPLEMENTED. NOT executed.

### F-04. A missing memory file is valid

`NewResourceMemoryStore` ignores `os.IsNotExist` (`resource_memory.go` lines 128–130). `orchestra verify` prints `[PASS] ... valid` with zeros (`main.go` lines 610–613). IMPLEMENTED. NOT executed.

### F-05. A failed web fetch can return success

`getOfflineFixture` builds `{status: offline_fixture}` for any non-empty id. `Acquire` returns that payload and a nil error when the failure is not `ErrSSRFDetected` (`web.go` lines 144–154 and 217–243). IMPLEMENTED. NOT executed.

### F-06. Skipped visual QA is stored as a successful resource outcome

`stage_iterate.go` lines 34–38 and 164–168, as read in this session. The CLI text says the skip is not a pass. The memory row says success if any resource was acquired. IMPLEMENTED. NOT executed.

### F-07. Changelog and eval doc name routes that are not in the table

`CHANGELOG.md` and `docs/audit/EVAL-5.0.md` say `ow-browserize`, `ow-qa`, and `ow-author` are on routes. A search of `registries/routing.json` for those three ids found no matches. The table does contain `ow-site-recon`, `ow-app-recon`, `ow-scrape`, `ow-ui-assembly`, `ow-wave`, and `ow-mcp-audit`. VERIFIED by search.

### F-08. `orchestra plan` invents a task

Empty `--task` becomes “Build modern responsive web feature with strict typography and security headers” (`main.go` lines 403–405). IMPLEMENTED.

### F-09. Stale pins that still pass the version check

`check-versions.js` lines 29–37 require each `live_files` entry to contain the current version string somewhere. They do not reject an older pin in the same file. `ARCHITECTURE.md` line 3 says the contract pin is 4.6.0 and the title says 5.0.0. `WORKFLOW.md` line 3 says 4.5.0 is one source of truth, under a 5.0.0 title. `docs/CONDUCTOR.md` frontmatter says Orchestra 4.6.0. `kit/antigravity/ALWAYS-ON.md` is titled 5.0.0 and still contains the marker `BEGIN ORCHESTRA 4.6.0`. `registries/oss-index.json` is version 3.3.0 and is not in `live_files`. The script was read and not executed. A file that contains both pins would pass that loop. That conclusion is from reading the loop, not from running it.

### F-10. Proprietary document skills in the public tree

**Withdrawn 2026-11-13.** `git ls-files` lists none of `skills/docx`, `skills/pdf`, `skills/pptx`, `skills/xlsx`, `skills/fullstack-mastery`, or `skills/creative-3d-web`. `git check-ignore -v` points at `.gitignore` lines 66–71. The trees are on disk and are not in the public tree.

`skills/docx/LICENSE.txt` and `skills/pdf/LICENSE.txt` are Anthropic, PBC, 2025. The text says users may not extract the materials from the Services, retain copies outside the Services, or distribute them to any third party. `THIRD_PARTY_NOTICES.md` does not mention them. They are not in `host-stack.json`. They are in the git tree. `docs/ROUTING.md` says docx, pptx, pdf, and xlsx are not in `skills/`. That sentence is false. VERIFIED by opening the license and the route line.

### F-11. SQLite memory is dead code

`memory.go` `NewStore` does not call `sql.Open`. `RetrieveByTags` returns an empty slice. The only definition of `NewStore` is that function. `ConflictWinner` is called from `phase3f_test.go` only. Durable policy that does run as a function is `AdmitDurable` and `SelectDurable` in `reconcile.go`, with `MaxDurableFacts = 8`. The engine pipeline does not call `LoadDurable`. VERIFIED by search and by reading `memory.go`.

### F-12. Hook gaps

`kit/hooks/guard.mjs` default profile is `standard`. A completion claim without command-shaped evidence is allowed unless the profile is `strict`. `npx skills add --all` is denied only when `--all` follows `add` immediately. `skills add <pack> --all` does not match that regex. `curl|sh` is denied. `iex (irm ...)` is not that pattern. MCP allow is a substring search of `kit/mcp-core.json`, and a missing file allows every server. IMPLEMENTED. The tests were not run.

### F-13. Context budget doc is stale

`docs/CONTEXT-BUDGET.md` says `AGENTS.md` is 1,847 characters. The file is 2,049 bytes and 2,049 characters (no CRLF expansion on this measurement). A research reader ran `measure-context.mjs`, which rewrites the budget doc. That rewrite was restored with `git checkout` before these audit files were added. The 25,614 character skill-description total in the budget doc was not re-measured.

### F-15. `orchestra classify --task` is not a flag

`WORKFLOW.md` lines 27–28 show `orchestra classify --task "<description>"`. `runClassify` defines `-graph`, `-json`, and `-assume` only, and reads the task from positional args (`main.go` lines 313–320). An unknown flag exits. VERIFIED by opening both files.

### F-16. Plan skips Research and then names it

`DesignPipeline.Plan` runs stages 0, 1, and 3 (Discover, Classify, Synthesize) and does not call stage 2 (`pipeline.go` lines 100–124). The `PipelineStage` string it returns still lists Research (`pipeline.go` line 134). `ResearchStage.ShouldSkip` returns true when the capability has no research worker, and also for a non-visual task with no gap technologies (`stage_research.go` lines 23–32). `ImplementStage` still writes `CompletedSteps` that include Research (`stage_implement.go` line 352). VERIFIED by opening those lines.

### F-17. Route agreement is not in CI

A search of `.github/` found no `check-routing` and no `check-route-agreement`. `ci.yml` and `hygiene.yml` were described in the forensic report. This session confirmed the names are absent from `.github/`. VERIFIED by search. Whether the scripts themselves are correct was NOT_RUN.

### F-18. The architecture doc undercounts the graph

`ARCHITECTURE.md` line 9 says the graph has 17 rows. This session’s parse of `design-resource-graph.json` counted 33 capabilities and 42 domains. The [Skills and agents](9d967d24-6ae9-48bf-873b-444d3877a288) reader counted 33 capability objects. VERIFIED for the file sentence and the parse. The two lists (36 routing jobs and 33 graph capabilities) are different selectors. `compile` reads `routing.json`. `orchestra classify` scores the graph. They can name different skills for one sentence. That split is INFERRED from the two entry points, and it matches the skills reader’s description of `Compile` versus `ClassifyBrief`.

### F-19. `ow-focus` is not a route load

`registries/routing.json` has no `ow-focus`. `docs/CONDUCTOR.md` line 61 says `ow-forge` and `ow-focus` run before classify. `docs/audit/STUB-TRIAGE.md` says merging `ow-focus` into `ow-forge` would drop a name the table uses. The route table does not use the name. The conductor doc does. VERIFIED by search and by opening line 61.

### F-14. `inspect.go` fetches without the private-IP check

`HTTPClient` follows up to 5 redirects and does not call the SSRF check used by the web adapter. The comment says fetched text is untrusted data. `infer.go` still classifies on substrings of that text. The security reader reported this. `inspect.go` lines 15–26 were opened here. The infer function was not re-opened, so the substring classification stays REPORTED.

## Skill and agent inventory

Classification uses the `stub: true` flag and one opened procedure (`ow-qa`). Bodies of the other 114 skills were not all read. Unread bodies are not given a quality grade.

- **IMPROVE, do not delete the name:** the 36 `stub: true` skills listed in `docs/audit/STUB-TRIAGE.md`. Hot path among them: `ow-conductor`, `ow-forge`, `ow-backend-architecture`, `ow-vault`, and `ow-design-intelligence`. `ow-ui-assembly` is not in that stub list.
- **KEEP, already procedural, not promoted:** `ow-qa`, `ow-eval`, `ow-semgrep`, `ow-site-recon`, `ow-app-recon`, `ow-scrape`, `ow-ui-assembly`, `ow-browserize`, `ow-wave`, `ow-author`, `ow-mcp-audit`. `ow-qa` was opened. The others match the stub scan (no `stub: true`, and they have `## Procedure`). Their text was not graded.
- **KEEP as vendored engineering skills:** `systematic-debugging`, `test-driven-development`, `verification-before-completion`, `receiving-code-review`, `grilling`, `domain-modeling`, `to-tickets`. Presence verified for those seven paths. Contents were not re-read.
- **LEAVE UNLOADED:** `creative-3d-web` and `fullstack-mastery`. Both have real `SKILL.md` bodies. They are off `host-stack.json`. The route rules say to leave them unloaded. Keeping the files in `skills/` invites a model to read them anyway.
- **REMOVE from the public tree when a code change is authorized:** `docx`, `pdf`, `pptx`, `xlsx`. License text forbids distribution. `docx` and `pdf` licenses were opened. `pptx` and `xlsx` were not opened; treat them as the same family only as a hypothesis until those license files are read. The route calls them host-provided.

Agents: 9 EGM reviewers and 12 `ow-*` reviewers. `ow-reviewer` and `ow-simplifier` are the two passes the contract names. The EGM set is the spec-review route (`egm-draft`, `egm-review`, `egm-resolve` are skills; the agent files are the lenses). This session did not score whether each agent changes a diff.

Commands: 22 markdown files under `commands/`. `render-adapters.mjs` copies them to `.cursor/commands` only. Claude command copies are a sync-time write in `ow-sync.mjs`, not a second source tree in this listing.

## Tests

NOT_RUN locally: `go test ./...`, `node --test`, `hygiene.js`, `check-routing.js`, `orchestra compile` on a paraphrase set, `orchestra doctor`, `orchestra verify`.

The historical `docs/audit/BASELINE-4.7.md` says 770 Go tests passed on the 4.6.0 tree. That number is not a result for `24ad53c`.

Remote: CI and Hygiene conclusions were `success` on `24ad53c`. Pass counts inside those logs are UNKNOWN.

## What this baseline did not do

- No product file was edited except the restored budget doc, which was returned to HEAD, and the new files in `docs/audit/`.
- No vault file was opened.
- No exploit was run.
- No model trial was run.
- The first draft said the skills and competitor readers had not finished. Both reports arrived later. See the red-team note.

## Red team

Done after the drafts, in this session. Not a separate unfinished agent.

Reopened and left standing: `stage_implement.go` lines 372–373, `check-versions.js` lines 29–37, `docs/ROUTING.md` line 33, `docs/CONTEXT-BUDGET.md` still saying 1,847 after the restore, and `kit/antigravity/ALWAYS-ON.md` line 5 (`AGENTS.md` wins when the note disagrees).

Corrected: answer 19 had said `ow-ui-assembly` needs a procedure. The stub scan puts that skill with the files that have `## Procedure`. `ow-design-intelligence` is the stub. The inventory sentence that mixed the two was rewritten.

After the drafts, [Repo forensic](97661c75-d7de-4386-8103-d2de755c6fed), [Skills and agents](9d967d24-6ae9-48bf-873b-444d3877a288), and [Competitors](a471d9f0-4685-4e58-89db-990a6993d9ae) finished. F-15 through F-19 were reopened in source before they were added. The forensic report’s glob of the six unregistered skill directories returned 0 files. That glob is wrong: `skills/creative-3d-web/SKILL.md`, `skills/docx/SKILL.md`, `skills/fullstack-mastery/SKILL.md`, and `skills/pdf/LICENSE.txt` were opened here, and the skills reader measured those `SKILL.md` files. Do not treat the empty-glob sentence as a finding.

The skills reader’s sizes for the 47 `ow-*` files and the 22 commands match the counts already in this baseline (36 stubs, 11 with procedures, 109 indexed skills, 115 directories). Their classification of unread bodies stays theirs. `ow-conductor`, `ow-forge`, `ow-backend-architecture`, and `ow-3d-engineering` are stubs with real tables inside. That matches IMPROVE, not DELETE.

The competitor reader compared other harnesses to `AGENTS.md` and `docs/ROUTING.md` only. This package compared them to the runtime as well. Where the two disagree, the runtime reading wins. Their extra mechanisms (ECC’s memory vault, GSD’s “absent means enabled”, Compound’s `docs/solutions/`) are written into the competitive analysis as REPORTED, because this session did not re-open those upstream files.

Left unresolved on purpose:

- `skills/pptx` and `skills/xlsx` licenses were not opened. B-07 says to open them before deleting.
- The SWE-Skills v1 PDF was not opened.
- The skills reader and the competitor reader never returned a final report. Their transcripts were not used.
- Local `go test` and the Actions logs stay NOT_RUN and NOT opened.
- Star counts are the API values from the earlier query. They were not fetched a second time.

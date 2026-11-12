# Implementation backlog 5.5

Recorded 2026-11-12. IDs are stable. This list does not authorize the edits. Token and latency costs are estimates from file size and from the shape of the change, not from a trace. “Code” means a later session. This research pass writes documents only.

## P0

### B-01. Quarantine doctor must not print PASS with 1598

- **Problem.** F-02. A missing `skills_library` path still prints PASS and the constant 1598.
- **Evidence.** `runtime/cmd/orchestra/main.go` lines 214–226.
- **Change.** Print PASS only when a real isolated directory was audited and the count came from that audit. Otherwise `NOT_RUN` or FAIL. Delete the constant.
- **Files.** `main.go`, `main_test.go`.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Benefit.** The doctor stops lying.
- **Tokens.** None.
- **Latency.** None.
- **Security.** S-01.
- **Maintenance.** One test.
- **Verify.** A test with a missing directory expects no PASS line.
- **Accept.** The output does not contain `1598` unless that count was read from disk.
- **Rollback.** Revert the commit.
- **Release.** 5.5.0 if code is authorized. Not in this research commit.

### B-02. Parity PASS requires at least one present host, and includes Codex

- **Problem.** F-03. Empty directories stay “complete.” Codex is in `ow-sync.mjs` and not in `ResolveHostDirs`.
- **Evidence.** `host_sync.go` lines 136–172. `ow-sync.mjs` lines 18–26. `main.go` lines 264–266, 597–600, 725–726.
- **Change.** If no host directory exists, `IsParityComplete` is false and the word is `NOT_RUN`. Add Codex to the Go map or stop claiming Codex in the Node tool. Pick one map.
- **Files.** `host_sync.go`, `ow-sync.mjs`, `main.go`, adapter tests.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Benefit.** Sync status matches the disk.
- **Tokens.** None.
- **Security.** S-07.
- **Verify.** Unit test with no directories. Unit test that the host list matches `hostDirs`.
- **Accept.** The 100 percent sentence is absent when a configured host directory is missing.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-03. Missing memory file is not valid

- **Problem.** F-04.
- **Evidence.** `resource_memory.go` lines 128–130. `main.go` lines 610–613.
- **Change.** A missing file prints `NOT_RUN` or creates the file only on an explicit `memory record`. `verify` does not say valid.
- **Files.** `resource_memory.go`, `main.go`.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Verify.** Test with a path that does not exist.
- **Accept.** No `[PASS]` for a missing file.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-04. Skipped visual QA is not OutcomeSuccess

- **Problem.** F-06. I-3.
- **Evidence.** `stage_iterate.go` lines 34–38 and 164–168.
- **Change.** Pass the real verdict into `recordIterateMemory`. `NOT_APPLICABLE` does not record success.
- **Files.** `stage_iterate.go` and its test.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Security.** S-01.
- **Verify.** Test with `CheckNotApplicable` and a non-empty acquired-resource list.
- **Accept.** Store receives no success outcome.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-05. Web acquire returns an error when the network fails

- **Problem.** F-05.
- **Evidence.** `web.go` lines 144–154 and 217–243.
- **Change.** Delete the any-id fallback, or return the error and attach the fixture only under an explicit offline flag that cannot be confused with a fetch.
- **Files.** `web.go`, acquisition tests.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Security.** S-02.
- **Verify.** A refused connection returns a non-nil error.
- **Accept.** Result metadata never says `network_fallback_fixture` on a nil error.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-06. Compile confidence is 0 on a miss

- **Problem.** F-01. I-1.
- **Evidence.** `compile.go` default `Confidence: 1` and the early return when `len(hits) == 0`. No test mentions `Confidence`.
- **Change.** Set confidence to 0 on that path. Add a fixture.
- **Files.** `compile.go`, `compile_test.go`.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Verify.** `Compile("zzzz-no-such-route", ...)`.
- **Accept.** Route `unspecified`, confidence 0.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-07. Remove the Anthropic document skills from the public tree

- **Closed 2026-11-13.** The six trees are gitignored. `git ls-files` lists none of them. No tracked file contains the proprietary distribution sentence. The Apache-2.0 notices under `skills/mcp-builder` and `skills/webapp-testing` stay, because that text is the Apache license.
- **Problem.** F-10. S-06.
- **Evidence.** `skills/docx/LICENSE.txt` and `skills/pdf/LICENSE.txt` forbid distribution. `pptx` and `xlsx` were not opened and must be checked in the same change. `THIRD_PARTY_NOTICES.md` does not name them. `docs/ROUTING.md` claims they are not in `skills/`.
- **Change.** Delete those four trees, or replace each with a pointer that contains no Anthropic text. Update the notices if any third-party text remains. This is not legal advice. If a lawyer says the files may stay, record that outside the repo and keep them out of the public remote until then.
- **Files.** `skills/docx`, `skills/pdf`, `skills/pptx`, `skills/xlsx`, notices, route prose if it still says they are absent.
- **Depends on.** A human yes, because it deletes files. This research pass does not delete them.
- **Complexity.** Small, and easy to get wrong if a license file is missed.
- **Security.** Redistribution.
- **Verify.** `git grep` for “Anthropic, PBC” under `skills/` is empty.
- **Accept.** The public tree has no file whose license says it may not be distributed.
- **Rollback.** The files can return only with a license that allows it.
- **Release.** Blocker for any public 5.5.0 tag.

### B-08. Empty `orchestra plan` does not invent a product

- **Problem.** F-08.
- **Evidence.** `main.go` lines 403–405.
- **Change.** Exit non-zero and print usage when the task is empty.
- **Files.** `main.go`.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Verify.** `orchestra plan` with no args exits 1.
- **Accept.** The typography sentence is gone.
- **Rollback.** Revert.
- **Release.** 5.5.0.

## P1

### B-09. Hot-path stubs get a short procedure

- **Problem.** 36 `ow-*` skills set `stub: true`. Conductor, forge, backend, vault, design intelligence, research, research paper, ml experiment, ship-safe, and observability are on a route or are step 0.
- **Evidence.** Frontmatter scan this session. `ow-qa` is the pattern to copy, and its body was opened.
- **Change.** Add Preconditions, Procedure, Output, Verify, Skip. First line of the procedure: the task wins. Remove `stub: true` only when those sections exist. Do not lengthen `docs/CONDUCTOR.md`.
- **Files.** Those `SKILL.md` files. `docs/CONDUCTOR.md` retitled to 5.0.0 or reduced to a pointer at the route table.
- **Depends on.** B-06 so the fast path does not load them.
- **Complexity.** Medium, because the temptation is to write a textbook.
- **Tokens.** Each finished skill should stay near the current stub size plus the five headings. A skill over the 300-character description cap in the budget doc fails `check` if that check exists. Description cap was REPORTED from `CONTEXT-BUDGET.md`, not re-run.
- **Verify.** `check-skills.js` if it already fails a stub. NOT_RUN here.
- **Accept.** `stub: true` is gone only on files that contain the five headings, and a fixture shows a conflicting path is not taken from the skill.
- **Rollback.** Restore the stub flag.
- **Release.** 5.5.0 for the hot path. The rest of the 36 can wait.

### B-10. Route table matches the changelog

- **Problem.** F-07.
- **Evidence.** No matches for `ow-browserize`, `ow-qa`, `ow-author` in `routing.json`.
- **Change.** Add a route for `ow-qa` with match words that do not steal “bug” from systematic debugging. Correct the changelog and `EVAL-5.0.md` so they stop saying browserize and author are routes, unless a real row is added for a reason other than the sentence.
- **Files.** `registries/routing.json`, rendered adapters via `render-adapters.mjs`, changelog.
- **Depends on.** B-09 for `ow-qa` only if the procedure is still thin. The procedure exists.
- **Complexity.** Small.
- **Verify.** `check-routing.js` and a compile fixture “browser bug on the running app” does not select the design lab.
- **Accept.** Changelog, eval doc, and `routing.json` agree.
- **Rollback.** Revert the row.
- **Release.** 5.5.0.

### B-11. Fast-path flag

- **Problem.** The fast path is prose.
- **Evidence.** `docs/CONDUCTOR.md` modes. `compile.go` sets disclosure from the route only.
- **Change.** `fast: true` when disclosure is L0 and the task contains a path-like token. Document it.
- **Files.** `compile.go`, tests, `docs/CONDUCTOR.md`.
- **Depends on.** B-06.
- **Complexity.** Small.
- **Tokens.** Avoids loading the 19,467-byte conductor doc when the host obeys the flag.
- **Verify.** Compile fixture for “fix nil pointer in runtime/internal/compile/compile.go”.
- **Accept.** Manifest says fast and L0.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-12. One version pin in the prose

- **Problem.** F-09.
- **Evidence.** `ARCHITECTURE.md` line 3, `WORKFLOW.md` line 3, `docs/CONDUCTOR.md` title.
- **Change.** Those sentences say 5.0.0 until a real 5.5.0 tag exists, then they say 5.5.0 in the same commit as `VERSION`.
- **Files.** Those three, and `check-versions.js` if it should fail on a second pin. The script was not re-opened. Changing it is optional.
- **Depends on.** The version decision.
- **Complexity.** Small.
- **Verify.** `git grep` for `4.6.0` and `4.5.0` in those files.
- **Accept.** A reader sees one pin.
- **Rollback.** Revert.
- **Release.** With whatever pin is actually tagged.

### B-13. Standard hook denies a bare DONE

- **Problem.** F-12. The contract says evidence before DONE. The default profile allows the claim.
- **Evidence.** `guard.mjs` lines 74–82.
- **Change.** In `standard`, a completion word without command-shaped evidence returns deny. Keep `minimal` and `off` as the explicit weaker profiles. Widen the `--all` regex only if a test shows the bypass. Do not try to parse every installer.
- **Files.** `guard.mjs`, `guard.test.mjs`.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Security.** S-04.
- **Verify.** The existing hook test plus a payload that says DONE and has no exit code.
- **Accept.** Exit non-zero on that payload in the default profile.
- **Rollback.** Revert.
- **Release.** 5.5.0.

### B-14. Inspect uses the same SSRF checks as Acquire

- **Problem.** S-03.
- **Evidence.** `inspect.go` lines 15–26. Web adapter checks were opened for the fixture path. The shared helper, if any, was not extracted.
- **Change.** Call the same private-IP and metadata-host checks before the GET and on redirects.
- **Files.** `inspect.go`, `web.go`.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Verify.** A test URL at `127.0.0.1` is refused.
- **Accept.** Inspect and Acquire share one check.
- **Rollback.** Revert.
- **Release.** 5.5.0.

## P2

### B-15. Run the existing paired fixtures

- **Problem.** Claims stay off because nothing was run. That is correct. A 5.5 marketing sentence needs the run.
- **Evidence.** `EVAL-5.0.md`, `compare.go`.
- **Change.** Execute `EVAL-5.5.md` on one model for T1 and T16 first. Save the raw output. Call `PromoteSkill` with the measured numbers. Leave `Enabled` false unless it returns true.
- **Files.** `runtime/benchmarks/`, `docs/audit/`.
- **Depends on.** B-01 through B-06, so the harness under test is not printing false PASS.
- **Complexity.** High. It needs a model and a budget.
- **Tokens.** The cost of the trials. Unknown until run.
- **Verify.** The trial file exists and the claim line matches `PromoteSkill`.
- **Accept.** Either `not_run` remains written down, or the numbers are the ones in the file.
- **Rollback.** Delete the trial file and keep claims off.
- **Release.** After 5.5.0, or inside it only as an appendix that still says none if the gate fails.

### B-16. Move unloaded stacks out of `skills/`

- **Problem.** `creative-3d-web` and `fullstack-mastery` are full procedures the routes say not to load. Search can still find them.
- **Evidence.** Both `SKILL.md` files were opened. They are absent from `host-stack.json`.
- **Change.** Move them to `templates/shelf/` or delete them if the shelf already has the note. Do not put them on a route.
- **Files.** Those directories, any docs that mention the path.
- **Depends on.** Nothing.
- **Complexity.** Small.
- **Verify.** Route rules still name them as unloaded, and `skills/` does not contain the directories.
- **Accept.** `canonicalIds` and the filesystem agree.
- **Rollback.** Move them back.
- **Release.** 5.5.0 or the one after. Not a blocker if B-07 is done.

### B-17. Dead SQLite store

- **Problem.** F-11.
- **Evidence.** `NewStore` has no callers.
- **Change.** Delete `memory.go`’s unused store or comment is the wrong fix. Delete it and the tests that only exist to compile it, if any.
- **Files.** `runtime/internal/memory/memory.go`.
- **Depends on.** A grep that still shows no callers at the time of the edit.
- **Complexity.** Small.
- **Verify.** `go test` of the memory package.
- **Accept.** No `sql.DB` field remains.
- **Rollback.** Revert.
- **Release.** 5.5.0.

## P3

### B-18. Junction test must fail on a real escape

- **Problem.** S-08. REPORTED, not re-opened.
- **Change.** Re-read the test. If it logs and returns success, make the escape a failure or delete the test so it does not look like coverage.
- **Files.** `runtime/internal/resources/challenger_test.go`.
- **Depends on.** Reading the test.
- **Release.** When confirmed.

### B-19. Diff vendored Pocock and Superpowers skills against upstream

- **Problem.** Upstream `mattpocock/skills` lists `grill-with-docs`. Local name is `grilling`. Bodies were not diffed.
- **Change.** A diff, then update only if the local file is a stale copy and the license still allows it. No new skills from that diff without B-15.
- **Release.** Optional.

## REJECT

| Id | Item | Reason |
| --- | --- | --- |
| R-1 | Install Superpowers, ECC, gstack, GSD, Compound, or Hermes as a conductor | Second plan, or an installer this repo forbids |
| R-2 | New vector brain or SQLite memory | Vault and `AdmitDurable` exist. The SQLite code is dead |
| R-3 | Model router that switches providers | The host owns the model |
| R-4 | Skill-count growth | Papers: most added skills do not move the mean, and self-written packs can hurt |
| R-5 | Auto-promotion of a self-written skill | SkillsBench self-generated arm |
| R-6 | Star counts or win-rate sentences in the README | Not measured |
| R-7 | Ten-layer orchestrator | `ARCHITECTURE-5.5.md` |
| R-8 | Always-on MCP set | Selection doc |
| R-9 | gstack’s CLAUDE.md rewrite and auto-update | Overwrites personal config |
| R-10 | Phone-home telemetry of the Superpowers logo kind | Their README describes it. Do not add one |

## Order

B-07 can happen alone and should happen before a public push of any new tag. B-01, B-02, B-03, B-04, B-05, B-06, and B-08 are independent and can land as one correctness commit. B-13 and B-14 next. B-09, B-10, B-11, B-12 with the docs. B-15 only after the PASS lines are honest. B-16 and B-17 when the tree is being cleaned. B-18 after the test is re-read. B-19 is optional.

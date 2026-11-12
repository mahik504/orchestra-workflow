# Eval 5.5

Recorded 2026-11-12. Every live result in this document is `not_run`. The method is here so a later session can run it without inventing a winner. `runtime/internal/benchmark/compare.go` keeps every comparison `Enabled: false` and `ClaimNone`. `PromoteSkill` returns false under three trials, on a negative success change, or when cost was not measured or rose by more than 10 percent. That function stays.

`docs/audit/EVAL-5.0.md` already says the 5.0 fixtures were not run. This file does not change that.

## Arms

Same model, same model version, same task, same starting commit, same tools, same acceptance checks.

| Arm | Setup |
| --- | --- |
| A | Host only. No Orchestra adapter and no Orchestra skill on the load path. A clean home, as `skills/ow-eval` already describes for the raw arm |
| B | Orchestra control plane only. `AGENTS.md` and the route receipt. No specialist `SKILL.md` |
| C | Orchestra plus the one route’s skills |
| D | Another harness, only when it can be installed without rewriting the operator’s personal config, and only on a throwaway home. If install is destructive or refused, D is `not_run` |

Do not compare arm C on one model with arm A on another. Do not enable D for Superpowers’ session-start plugin inside the same home as Orchestra. That is two conductors.

## Tasks

Write the task from the job, not from the skill text. SWE-Skills-Bench’s requirements are reverse-engineered from the skill, and the authors list that as a limitation. Copying that method would flatter the skill.

| Id | Task | Acceptance, deterministic where possible |
| --- | --- | --- |
| T1 | Bug in an unfamiliar repo | Failing test becomes passing. No new framework. Nearby tests still pass |
| T2 | Backend feature | Request validated, one error shape, test covers the new path |
| T3 | Architecture change | The named module moves. Public behavior the tests cover stays |
| T4 | Landing page | Builds. Desktop and mobile viewports. Human blind preference is a separate sheet, not a Lighthouse score |
| T5 | Multi-page frontend | Shared layout. Each route renders. No console errors on the journey |
| T6 | Expo screen | Official Expo skills only if arm C’s route says so. The screen navigates |
| T7 | Responsive and accessibility | The named breakpoints. An axe or equivalent run. Critical violations counted |
| T8 | Implement a cited behavior | The citation URL was opened. The behavior matches the source, or the diff says it does not |
| T9 | Data transform | Schema in, schema out, a fixture with a bad row rejected |
| T10 | ML or DL | Train only on data the fixture includes. Metric computed by the script. No hand-typed accuracy |
| T11 | Research note | Every citation resolves. No invented DOI |
| T12 | Academic deck | Required sections present. Speaker notes exist. No purple template dump as the only design |
| T13 | Browser automation | The script asserts the outcome. Credentials are not stored |
| T14 | Site recon of a public page | Tokens and structure recorded. Assets, copy, and logo not copied |
| T15 | Scrape | Stops at a 403, a login wall, or robots. Rows match the schema |
| T16 | Debugging a failing test | Root cause named. The test fails before the fix and passes after |
| T17 | Multi-session | Second session continues from `STATE.md` without the first transcript |
| T18 | Security-sensitive change | No secret in the diff. Auth boundary tested. Hook would deny a `.env` read |

T4 and T7 need a blind human pass for visual quality: hierarchy, type, spacing, consistency. Lighthouse and axe are recorded and are not the design score.

## Measures

Per trial: success against the acceptance checks, regressions, build result, defect count, rework (a second agent turn after a failed check), tokens if the host reports them, tool calls, model calls, wall time, human interventions, and whether a failure recovered.

Cost is whatever the host’s usage line reports. If the host does not report tokens, the cost cell is `not_measured`, and `PromoteSkill` must keep returning false because `costRatio <= 0` is “cost was not measured.”

## Statistics

Three trials are exploratory. Report the three outcomes and the median. Do not report a p-value on n = 3 as a win.

A later confirmatory set, if anyone claims superiority, needs a pre-registered task list, at least the trial count written down before the run, and a comparison of paired differences. This document does not pick a test until that protocol exists. Until then the labels are:

| Label | Rule |
| --- | --- |
| Improvement | Arm C passes checks arm B failed, on the same tasks, and cost ratio is at most 1.10, on at least three trials |
| Regression | Arm C fails a check arm B passed, or cost ratio is over 1.10 |
| No detectable difference | Same pass pattern |
| Inconclusive | Missing trials, different models, unmeasured cost, or a task written from the skill |

`PromoteSkill(successDelta, costRatio, trials)` is the code form of the improvement row. A skill is not added to a default Load line on any other row.

## Current results

| Fixture | Result |
| --- | --- |
| `runtime/benchmarks/orchestra-47/` arms described in `docs/audit/BASELINE-4.7.md` | `not_run` for 5.0.0. The 4.7 note says they were `not_run` then too |
| `compare.go` suite c1–c6 | `Enabled: false`, `ClaimNone` |
| Landing trials mentioned in the 4.5 changelog | Historical note in `CHANGELOG.md`. Not re-run. Not a 5.5 claim |
| This session’s model comparison | `not_run` |

## Frontend sheet

For T4, T5, and T7, a reviewer who did not see the arm labels scores: hierarchy, type, spacing, consistency, responsiveness, and whether the page does the job. A high Lighthouse performance number with weak hierarchy is not a win. Record the reviewer’s words. Do not average them into a fake 0–10.

## Release rule

5.5.0 may ship the correctness fixes in the roadmap with the claim still none. 5.5.0 may not ship a sentence that Orchestra beats a raw host, Superpowers, ECC, or gstack. That sentence waits for an improvement row with the trial file in the repo.

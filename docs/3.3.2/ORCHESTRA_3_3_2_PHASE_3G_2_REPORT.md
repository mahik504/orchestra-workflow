# Orchestra 3.3.2 — Phase 3G.2 report

**Pilot.** Preflight resolved `agent.cmd` and its version. The operator label `grok-4.7` is not an exact id in that installation's model list, so T1 and T3 stay `unavailable`. `claim` is `none`. No other model was tried. This is not evidence that Orchestra improves output.

VERSION stays `3.3.1`. No commit. T2 and T4–T10 stay disabled. No 3H.

## Preflight

`ORCHESTRA_BENCH_REAL=1` ran `TestPhase3G2_RealHostPilot`.

- The executable resolved. The recorded basename is `agent.cmd`. `PATH` was not changed.
- `agent --version` returned `2026.09.08-6caf4ff`.
- `agent --help` lists `--print`, `--force`, and `--trust`.
- `agent models` ran. The operator label is not an exact id in that list. A longer id that only begins with the same characters does not count.
- The disposable `--print` write did not run. `smoke_wrote` is false.
- Blocker, on both cells: `model label is not in the agent model list`.

No model router, second wrapper, or conductor change was added.

## The eight questions

1. **Executable, version, and slug.** Executable and version resolved. The slug did not.
2. **Same fixture, task, and acceptance hashes.** The arms did not start. There is no fixture hash and no task hash for this pilot.
3. **Did `agent.cmd` edit a worktree?** No. Preflight stopped before either host process. `host_edited` and `task_completed` were not measured.
4. **Per dimension.** No diff, no passing dimension, and no Orchestra byte count from a host pair. There is no overall score.
5. **Did the label enter the mission, capsule, capability file, acceptance text, or edited source?** The host did not run, so it wrote none of those. Fake-host coverage rejects a label written into `sum.go` as `invalid_pair`.
6. **Baseline artifacts and protected paths.** Not observed on a real pair, because the pair did not start. A fake host that creates `.orchestra` on the baseline is `invalid_pair`. A fake write outside the fixture records the path class only.
7. **Shared Cursor profile.** Not recorded on this run. The confounder is attached only when a real pair starts. This run did not reach that step, and it did not change Cursor settings, MCP config, or rules.
8. **Is T2 justified?** No. Both real cells are `unavailable`. Two pairs were not completed. T2 stays disabled.

## What the fake host still proves

Default `go test` does not launch Cursor Agent. With the fake host, the wrapper argv is `--print`, the workspace, the model label, and the prompt. `--force` and `--trust` appear only when the help text lists them. `--plan`, `--worktree`, `--add-dir`, and `--approve-mcps` are absent. A non-zero exit is `host_error`. An unrelated `sum.go` edit is `host_edited` and not `task_completed`. When one arm fixes T3 and the other does not, `claim` stays `none`.

## Checks

`go test ./internal/... -count=1` passed. That includes `internal/benchmark` and `internal/visual`.

`TestLayerASmokePlaywright` passed. It is a regression check, not a visual study.

From the repo root, host-stack reported 44 skills and 0 problems. Hygiene scanned 725 tracked files, `VERSION` `3.3.1`, 0 failures.

No commit. No tag. No version bump.

## Stop

Do not enable T2. Do not start 3H. Do not pick a nearby model id to force a pair. A later pair still needs an exact operator label, one T1 pair, one T3 pair, and repeated trials before any effectiveness claim.

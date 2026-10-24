# Orchestra 3.3.2 — Phase 3G.1 report

**Pilot.** The harness can compare two host processes. No host executable was configured, so T1 and T3 stay `unavailable`. `claim` is `none`. This is not evidence that Orchestra improves output.

VERSION stays `3.3.1`. No commit. T2 and T4–T10 stay disabled. No 3H.

## What the harness does

`RunPair` in `runtime/internal/benchmark` copies the fixture, runs the existing pipeline on the Orchestra copy, and exports `capsule.json` with `mission.Compile`. The baseline copy does not call `Execute`.

A fresh process is the host. The executable comes from `ORCHESTRA_BENCH_HOST`. The opaque label comes from `ORCHESTRA_BENCH_MODEL_LABEL` and is stored on the result. It is not passed to `capability.Select`. The child environment does not receive that label, the Brain path, or a mode flag.

`go test ./internal/...` does not call a model. The real pair is `TestPhase3G1_RealHostPilot`, which runs only when `ORCHESTRA_BENCH_REAL=1`.

## Eight questions

1. **Same hash and same host label.** The fake pairs start from one fixture hash. A label mismatch returns `invalid_pair` before a process starts (`TestPhase3G1_LabelMismatch`).
2. **Did a fresh process write the trees?** Not a model. `ORCHESTRA_BENCH_HOST` and `ORCHESTRA_BENCH_REAL` were empty. No `agent` or `cursor` command was on `PATH`. T1 and T3 are `unavailable`. The fake host is the test binary, and it is not a coding agent.
3. **Isolation.** On the fake pairs, the baseline had no `.orchestra` directory after the host exited (`TestPhase3G1_T1Noop`, `TestPhase3G1_T3Fix`).
4. **Per dimension.** With the fake host, T1 already passes and a no-op stays a pass. T3 passes when `sum.go` loses the seeded `+ 1`. Changing `sum_test.go` so the bug matches the test is a hard fail even if `go test` then passes. Adding `extra.go` is a scope failure. There is no overall score.
5. **Orchestration bytes.** The Orchestra arm still records mission, capability, and capsule bytes from the existing trace. The fake T1 pair kept intensity `LIGHT`, capability `coding.worker`, and zero directions, browser calls, and Brain facts.
6. **Model label.** `grok-4.7` was used as the fake pair's metadata. The prompt builder does not copy it into the task prompt. A leak in the mission, capsule, capability file, prompt, or acceptance text stops the pair as `invalid_pair`.
7. **Timeout.** A host that does not finish is `timeout`. The other arm is not recorded as a win.
8. **T2.** Not justified. Both real pairs are `unavailable`. T2 stays `enabled: false`.

## Checks

`go test ./internal/... -count=1` passed, including `internal/benchmark`.

From the repo root, host-stack reported 44 skills and 0 problems. Hygiene scanned 725 files, `VERSION` `3.3.1`, 0 failures.

No commit. No tag. No version bump.

## Stop

Do not enable T2. Do not start 3H. A later real pair needs `ORCHESTRA_BENCH_REAL=1`, an operator-supplied executable, and the same opaque label on both arms. Until that process writes both trees, Orchestra's effect on coding output is unmeasured.

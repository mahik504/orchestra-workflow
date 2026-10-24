# Orchestra 3.3.2 — Phase 3G.1 host pilot

3G stays accepted. This pass asks whether one fresh host, on the same fixture snapshot, produces a different tree when it sees the task alone versus the task plus the existing Orchestra packet.

One pair is an observation. `claim` stays `none`. It is not evidence that Orchestra improves output.

VERSION stays `3.3.1`. No commit unless asked. No 3H. T2 and T4–T10 stay disabled.

## Host boundary

The host is a separate process. The benchmark does not choose a model and does not grow a host manager, model registry, worker router, or conductor change.

`ORCHESTRA_BENCH_HOST` is the executable path. The child receives the worktree path and the prompt path. One timeout covers both arms. `ORCHESTRA_BENCH_MODEL_LABEL` is stored on the result only. It is not passed to `capability.Select`, and it must not appear in the prompt, mission, capsule, capability file, or acceptance text. `grok-4.7` is allowed as that label. The result host label for a real run is `cursor`.

If the executable is empty, the cell is `unavailable`. No simulated diff and no token count.

`go test ./internal/...` uses a fake host. The real pair runs only when `ORCHESTRA_BENCH_REAL=1`.

## Pair

T1 and T3, one trial each. Baseline runs first. Each arm gets its own temp copy, `GOCACHE`, and module cache.

The fixture hash is recorded before either host starts. A hash mismatch or two different model labels is `invalid_pair`.

Baseline prompt: task and acceptance lines. No mission, capsule, or Brain file. No `Execute`.

Orchestra prompt: the same task and acceptance, plus the mission, capsule, and capability files written by the existing pipeline. The capsule file is an export of `mission.Compile` from those artifacts.

A timeout is that trial's result. The other arm is not a win. The temp directories are removed after the result JSON is written.

## Scores

No overall score.

Hard checks, declared before the run:

- `go test` build and tests.
- T3 locks `sum_test.go`, `task.txt`, `acceptance.txt`, and `go.mod`. T3 may change `sum.go` only.
- T1 locks `go.mod`, `task.txt`, and `acceptance.txt`. Extra module requirements fail.
- A positive judgment cannot clear a hard fail.

Artifact sheet, separate: `git diff --stat`, changed paths, diff size, host wall time, exit status. Host retries and tokens stay `unavailable` unless the process itself returns a number.

T1 and T3 do not launch Chromium.

## Done

The pilot is done when the fake-host tests pass and each of T1 and T3 is either a real `compared` pair with `claim: none` or an explicit `unavailable` cell.

T2 stays off. A later pass may consider it only if both real pairs finished, a test-edit cheat stayed a hard fail, no model label entered the mission or capsule, and no conductor or registry was added.

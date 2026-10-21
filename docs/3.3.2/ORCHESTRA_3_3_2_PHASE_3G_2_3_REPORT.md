# Orchestra 3.3.2 — Phase 3G.2.3 report

**Pilot.** The 3G.2.2 T1 completion row is superseded. The re-seeded T1 pair is `compared`, and both arms finished the task. The 3G.2.2 T3 row stands. `claim` is `none`. This is not evidence that Orchestra improves output.

VERSION stays `3.3.1`. No commit. No T2. No 3H.

## 3G.2.2, as recorded

The JSON files from that run were not rewritten.

### T1

Status: `timeout`. The Orchestra arm hit the 10 minute cap. The pair stayed `timeout`. The other arm is not a win.

| | Baseline | Orchestra |
| --- | --- | --- |
| Model | `grok-4.7-high` | `grok-4.7-high` |
| Exit | 0 | -1 |
| Timed out | no | yes |
| Wall | 96190 ms | 600366 ms |
| Recorded `task_completed` | true | true |
| Changed paths | none | none |
| Diff bytes | 0 | 0 |
| Stdout bytes | 291 | 728 |
| Mission / capability / capsule | 0 / 0 / 0 | 1942 / 266 / 1014 |

Fixture hash `65a972853d968fdf51abf6878dd497778a736ac22880d2200b96cb02b868cc55`. Neither arm changed a product file.

### T3

Status: `compared`. Both arms changed only `sum.go`. Both exited 0 inside the limit. Both are `task_completed`.

| | Baseline | Orchestra |
| --- | --- | --- |
| Model | `grok-4.7-high` | `grok-4.7-high` |
| Exit | 0 | 0 |
| Timed out | no | no |
| Wall | 117337 ms | 133935 ms |
| `task_completed` | true | true |
| Changed paths | `sum.go` | `sum.go` |
| Diff bytes | 183 | 183 |
| Mission / capability / capsule | 0 / 0 / 0 | 1900 / 268 / 968 |

## Defects

The T1 fixture already rejected negative inputs, and `sum_test.go` already had `TestSumRejectsNegative`. `taskCompleted` checked file content only for `t3-debug`. For T1, passing tests, scope, and acceptance was enough. In `ScoreHost`, acceptance means the protected files kept their hashes. `required_behavior` was a copy of the test result. The starting tree passed all of that, so a no-op was `task_completed` true. `TestPhase3G2_Fingerprints` asserted that.

`taskCompleted` also ignored the host exit. In the saved T1 JSON the Orchestra arm has `timed_out` true and `host_exit` -1, and still `task_completed` true. The pair status was `timeout`, so the top-level verdict was the timeout. The arm record was wrong.

The timeout cause is unconfirmed. Confirmed facts: exit -1, 600366 ms, 728 stdout bytes, no tracked change. The host logs were deleted with the worktree. These are hypotheses, not findings: the agent kept going because the API was already there; the CLI waited; the agent edited files under `.orchestra/`, which the diff excludes.

`ORCHESTRA_MEMORY_PATH` was unset, so the Brain ledger was not fingerprinted. `ORCHESTRA_HOME` was set, and `hostEnv` did not drop it, so both arms inherited it. After that run the ledger's last write was 2026-09-22, and git status for `memory/` was clean. That rules out a ledger write. It does not rule out a read. The harness `Execute` used the temporary ledger set in `TestMain`.

## Corrections

`task.txt`, `acceptance.txt`, and `go.mod` are unchanged, so the task hash and the acceptance hash stay the same as 3G.2.2. `sum.go` is now an unvalidated stub. `sum_test.go` keeps `TestSumValid` only.

The harness owns `TestBenchOracleValid` and `TestBenchOracleNegative`. That file is written after the diff is taken, run with `go test -run "^TestBenchOracle"`, then deleted. It is not in the fixture and not in either prompt. For T1, `required_behavior` passes only when the agent's tests pass and the oracle passes. T3 has no oracle, so its `required_behavior` is still the test result.

`taskCompleted` is false when the artifact is missing, `timed_out` is set, or `host_exit` is not 0. T1 also requires `required_behavior` and a change to `sum_test.go`. The T3 content check is unchanged.

`hostEnv` drops `ORCHESTRA_HOME`. When that variable is set, the guards cover the home directory's git status and `memory/resource-memory.json`. The `ORCHESTRA_MEMORY_PATH` guard is unchanged. The artifact records `stderr_bytes` and `harness_changed`. `harness_changed` compares `.orchestra/` and `prompt.txt` before and after the host. No transcript is stored.

Regression coverage: a T1 no-op is not complete; a validated `sum.go` plus a negative test is complete on both arms; code without a new test is not complete; the oracle file is gone after scoring; a timed-out arm and a non-zero exit are not complete; `ORCHESTRA_HOME` is absent from the host environment; a ledger write is class `brain`; the saved T3 fixture, task, and acceptance hashes still match.

Before the rerun: `gofmt` was clean, `go vet` passed, `go test ./internal/benchmark` passed, `go test ./...` passed, host-stack reported 44 skills and 0 problems, and hygiene reported 725 tracked files, VERSION `3.3.1`, and 0 failures.

## T1 rerun

One pair. Results are `3g23-preflight.json` and `3g23-t1-api.json`.

Preflight was ok. Basename `agent.cmd`. Version `2026.10.01-e373342`. The 3G.2.2 preflight recorded `2026.09.08-6caf4ff`. `grok-4.7-high` was an exact id. `--print`, `--force`, and `--trust` are listed. The smoke wrote `marker.txt`. Go: `go version go1.27.0 windows/amd64`.

Status: `compared`. `claim` is `none`.

| | Baseline | Orchestra |
| --- | --- | --- |
| Model | `grok-4.7-high` | `grok-4.7-high` |
| Exit | 0 | 0 |
| Timed out | no | no |
| Wall | 166561 ms | 172559 ms |
| `host_edited` | true | true |
| `task_completed` | true | true |
| Changed paths | `sum.go`, `sum_test.go` | `sum.go`, `sum_test.go` |
| Diff bytes | 592 | 1027 |
| Build, tests, required behavior, acceptance, scope | pass | pass |
| Stdout bytes | 433 | 477 |
| Stderr bytes | 0 | 0 |
| `harness_changed` | false | false |
| Mission / capability / capsule | 0 / 0 / 0 | 1930 / 266 / 996 |
| Capability ids | none | `coding.worker` |

Intensity on the Orchestra arm is `LIGHT`. Task hash `61e1e1f90eba543f22280c4ee8f5262b665dbdf44cc5c9e76bed482119e24d28` and acceptance hash `7ad4f9a508fe7dc6221527790b3b08a5176b8deb59322e9c745a506093c8e11d` match across the arms and match 3G.2.2, because those text files did not change. Fixture hash `45b804163dc2c603905974e876933eb387b5e779ba8df7a6b1e420eba0eb7a50` is the stub. The prompt hashes differ by the harness file references. `prompt_equivalent` is true. The baseline stayed free of Orchestra artifacts. The worktree was removed.

`stdout_bytes` in this run is stdout only. The 3G.2.2 field added stderr into the same number. Stderr here is 0.

Both arms passed the oracle and added a test. The diffs are different sizes. One trial. Same model. This pair does not separate the arms.

## Why T3 was not rerun

The T3 policy has an empty oracle, so `required_behavior` is unchanged. Both saved T3 arms have `host_exit` 0 and `timed_out` false, so the new exit gate does not change those records. `TestPhase3G23_T3InputsUnchanged` matched the saved fixture hash, task hash, and acceptance hash.

One difference stays on the record: during 3G.2.2 the host inherited `ORCHESTRA_HOME`. This rerun withholds it. T3 was not repeated under that denial.

## Integrity

Operator fingerprints taken before the rerun and after it match.

- `method_repo` unchanged
- `fixture` unchanged
- `cursor_config` present and unchanged
- the Cursor user MCP file was absent both times
- the home MCP file was present and unchanged
- Brain home git status unchanged
- Brain ledger file hash unchanged
- `ORCHESTRA_MEMORY_PATH` was unset, so that path was not a guard
- `ORCHESTRA_HOME` was set, withheld from the host, and guarded
- the result has no `external_class`
- `claim` is `none`

Both arms used the same authenticated Cursor profile, including the MCP configuration and rules the CLI already loads. That is a confounder. It is not model-only isolation.

## Limitations

- One corrected T1 trial.
- The 3G.2.2 T1 timeout cause is still unconfirmed. This rerun finished, and `harness_changed` was false.
- T3 still cannot separate the arms, and that pair ran while the host could see `ORCHESTRA_HOME`.
- Host retries and tokens stay unavailable.
- `claim: none`.

## Stop

Do not enable T2. Do not start 3H. Do not add a model router. Two completed arms are not a reason to change the conductor.

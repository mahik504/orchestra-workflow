# Orchestra 3.3.2 — Phase 3G.2.4 report

**Pilot.** One T3 pair, with `ORCHESTRA_HOME` withheld from both host arms. Status is `compared`. Both arms removed the seeded `+ 1`, passed tests, and stayed on `sum.go`. `claim` is `none`. This is not evidence that Orchestra improves output.

VERSION stays `3.3.1`. No commit. No T2. No 3H. The 3G.2.2 and 3G.2.3 reports and their JSON files were not rewritten.

## Earlier T3 row

From `3g2-t3-debug.json`, unchanged. That pair ran while the host inherited `ORCHESTRA_HOME`.

Status: `compared`. Both arms changed only `sum.go`, exited 0, and were `task_completed`.

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

## This pair

Results are `3g24-preflight.json` and `3g24-t3-debug.json`.

Preflight was ok. Basename `agent.cmd`. Version `2026.10.01-e373342`. `grok-4.7-high` was an exact id. `--print`, `--force`, and `--trust` are listed. The smoke wrote `marker.txt`. Go: `go version go1.27.0 windows/amd64`.

`hostEnv` drops `ORCHESTRA_HOME` before either arm starts. The same denial was already in the harness. This run did not add isolation code.

Status: `compared`. `claim` is `none`.

| | Baseline | Orchestra |
| --- | --- | --- |
| Model | `grok-4.7-high` | `grok-4.7-high` |
| Exit | 0 | 0 |
| Timed out | no | no |
| Wall | 104157 ms | 129875 ms |
| `host_edited` | true | true |
| `task_completed` | true | true |
| Changed paths | `sum.go` | `sum.go` |
| Diff bytes | 183 | 183 |
| Build, tests, required behavior, acceptance, scope | pass | pass |
| Stdout bytes | 72 | 176 |
| Stderr bytes | 0 | 0 |
| `harness_changed` | false | false |
| Mission / capability / capsule | 0 / 0 / 0 | 1884 / 268 / 944 |
| Capability ids | none | `coding.worker` |

Intensity on the Orchestra arm is `LIGHT`. Fixture hash `a1400069e5e04bd50b185cd756a4d8806b98922b21225f7d517de15df6416969`, task hash `b7b77e85fa99a43eed6b52b4b3dd8f41952ed8f8bf222ea7c7c5e7647868e357`, and acceptance hash `8cf52e8dbadbf055cac938331ecbb25643d3e5e20f056b4ad3f68d99c5d7b533` match across the arms and match the saved 3G.2.2 T3 row. The prompt hashes differ by the harness file references. `prompt_equivalent` is true. The baseline stayed free of Orchestra artifacts. The worktree was removed.

Both diffs are 183 bytes and touch only `sum.go`. The task text names the return-value fix, so this pair still cannot separate the arms.

## Integrity

Operator fingerprints taken before the pair and after it match.

- `method_repo` unchanged
- T3 fixture unchanged
- `cursor_config` present and unchanged
- the Cursor user MCP file was absent both times
- the home MCP file was present and unchanged
- Brain home git status unchanged
- Brain ledger file hash unchanged
- `ORCHESTRA_MEMORY_PATH` was unset
- `ORCHESTRA_HOME` was set, so the ledger file under that home was the fingerprinted ledger
- the result has no `external_class`
- `claim` is `none`

The ledger path was available, so it is not unverified.

Both arms used the same authenticated Cursor profile, including the MCP configuration and rules the CLI already loads. That is a confounder. It is not model-only isolation.

## Limitations

- One T3 trial with `ORCHESTRA_HOME` withheld.
- The task text tells the host to change the return value. Both arms did. That does not show a difference.
- Host retries and tokens stay unavailable.
- `claim: none`.

## Stop

Do not enable T2. Do not start 3H. Do not add a model router. Do not rerun this pair. VERSION stays `3.3.1`.

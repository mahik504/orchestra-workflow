# Orchestra 3.3.2 — Phase 3G.2.2 report

The T1 completion row below is superseded by Phase 3G.2.3. The T3 row stands.

**Pilot.** Smoke wrote `marker.txt` with `grok-4.7-high`. T1 is `timeout`. T3 is `compared`, and both arms removed the seeded `+ 1`. `claim` is `none`. This is not evidence that Orchestra improves output.

VERSION stays `3.3.1`. No commit. No T2. No 3H.

## Preflight

- Basename: `agent.cmd`
- Version: `2026.09.08-6caf4ff`
- `grok-4.7-high` was an exact id in `agent models`
- `--print`, `--force`, and `--trust` are listed
- Smoke wrote `marker.txt` and the temp directory was deleted
- Go: `go version go1.27.0 windows/amd64`

Both arms were told to read `prompt.txt` in the workspace. The hashed file is the task. The argv line is the same on both arms. The Orchestra file also names the mission, capsule, and capability files.

## T1

Status: `timeout`. The Orchestra arm hit the 10 minute cap. That is not a win for the baseline, and it is not `task_completed` for the pair.

| | Baseline | Orchestra |
| --- | --- | --- |
| Model | `grok-4.7-high` | `grok-4.7-high` |
| Exit | 0 | -1 |
| Timed out | no | yes |
| Wall | 96190 ms | 600366 ms |
| `host_edited` | false | false |
| Recorded `task_completed` | true | true |
| Changed paths | none | none |
| Diff bytes | 0 | 0 |
| Tests | pass | pass |
| Acceptance | pass | pass |
| Scope | pass | pass |
| Mission / capability / capsule | 0 / 0 / 0 | 1942 / 266 / 1014 |
| Capability ids | none | `coding.worker` |

The fixture already passed. Neither arm changed a product file, so the scorer still marks each tree complete. The pair stays `timeout` because the Orchestra process did not exit in time.

## T3

Status: `compared`. Both arms changed only `sum.go`. Both are `host_edited` and `task_completed`.

| | Baseline | Orchestra |
| --- | --- | --- |
| Model | `grok-4.7-high` | `grok-4.7-high` |
| Exit | 0 | 0 |
| Timed out | no | no |
| Wall | 117337 ms | 133935 ms |
| `host_edited` | true | true |
| `task_completed` | true | true |
| Changed paths | `sum.go` | `sum.go` |
| Diff bytes | 183 | 183 |
| Tests | pass | pass |
| Acceptance | pass | pass |
| Scope | pass | pass |
| Mission / capability / capsule | 0 / 0 / 0 | 1900 / 268 / 968 |
| Capability ids | none | `coding.worker` |

Intensity on the Orchestra arm is `LIGHT`. Same model, same fixture, same task hash, same acceptance hash. The Orchestra prompt hash differs because of the harness file references.

## Pair integrity

Both pairs:

- `model_label` is `grok-4.7-high` on the pair and on both arms
- task hashes match across arms
- acceptance hashes match across arms
- `prompt_equivalent` is true
- both arms share one fixture hash
- baseline stayed free of Orchestra artifacts
- no external path class changed
- `claim` is `none`
- the temp worktree was removed

Operator fingerprints, taken before the pilot and after it, did not move: `method_repo`, `fixture`, `cursor_config`, `mcp_config`. The home MCP file was present and unchanged. The Cursor user MCP file was absent both times. `ORCHESTRA_MEMORY_PATH` was unset, so no brain file was fingerprinted.

Both arms used the same authenticated Cursor profile, including the MCP configuration and rules the CLI already loads. That is a confounder. It is not model-only isolation.

## Limitations

- T1 did not finish as a `compared` pair. A 10 minute timeout is not completion.
- T3 shows both arms can remove this seeded `+ 1`. It does not show a difference in the fix.
- Two tasks, one trial each. `claim: none`.
- Host retries and tokens stay unavailable.

## Stop

Do not enable T2. Do not start 3H. Do not add a model router. Do not treat the T3 tie, or the T1 timeout, as a reason to change the harness.

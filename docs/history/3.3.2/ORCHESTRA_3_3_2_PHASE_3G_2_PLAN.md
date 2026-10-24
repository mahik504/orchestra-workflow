# Orchestra 3.3.2 — Phase 3G.2 host validation

3G and 3G.1 stay accepted. This pass connects the existing harness to one real coding agent. It does not redesign the harness and it does not claim that Orchestra improves output.

VERSION stays `3.3.1`. No commit unless asked. No 3H. T2 and T4–T10 stay disabled. `claim` stays `none`.

## Boundary

`RunPair` still copies the fixture, runs `Execute` only on the Orchestra copy, and scores with the existing checks. The host is still one operator executable with no arguments. The working directory is the arm worktree. The child receives `ORCHESTRA_BENCH_WORKSPACE`, `ORCHESTRA_BENCH_PROMPT`, and `ORCHESTRA_BENCH_MODEL_LABEL`.

Cursor Agent does not read those variables. `runtime/benchmarks/cursor-host.ps1` does. It calls the absolute agent path in `ORCHESTRA_BENCH_AGENT`. The script has no model slug, no user path, and no API key.

`--print` is mandatory. `--force` and `--trust` are added only when that installation's help lists them. The prompt is a positional argument. The wrapper never passes `--plan`, `--mode ask`, `--worktree`, `--add-dir`, or `--approve-mcps`.

The model label is metadata for `--model`. It is not passed to `capability.Select` and it must not appear in the prompt, mission, capsule, capability file, acceptance text, or edited source.

## Preflight

Before a T1 or T3 pair: resolve `agent.cmd` without editing `PATH`, run `--version`, run `models`, confirm the operator label is in that list, then one disposable `--print` write in a temp directory. Delete that directory. If any step fails, the real cells are `unavailable` and no further architecture is added.

## Measurement

No overall score.

- `host_edited` means a non-harness file changed.
- `task_completed` is separate. T3 requires `sum.go` to stop returning `a + b + 1`, tests to pass, and protected files to stay unchanged. T1 requires passing tests and no file outside the allow-list.
- Task text and acceptance hashes must match across arms. The only allowed prompt difference is the harness mission, capsule, and capability references. The result stores hashes, not the prompt text.
- Both arms share one Cursor profile, including the MCP configuration and rules the CLI already loads. That is a confounder. It is not model-only isolation. Those files are not modified by this phase.
- A write outside the fixture, including the method repo, the Brain ledger, Cursor config, or the original fixture, makes the pair `invalid_pair`. The report names the path class only.
- Timeout, a non-zero exit, and a missing executable do not make the other arm a win.
- Two pairs are a pilot. `claim` stays `none` even if one arm passes and the other fails.

## Tests

Default `go test` uses a fake host and does not launch Cursor Agent. The real pair runs only when `ORCHESTRA_BENCH_REAL=1`.

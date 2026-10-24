# Orchestra 3.3.2 — Phase 3G pilot report

**Pilot.** This report says whether the benchmark can measure Orchestra on T1 and T3. It does not say that Orchestra improves output.

VERSION stays `3.3.1`. No commit, no tag, no 3H. T2 and T4–T10 stay disabled. The paired host cell is `unavailable`.

## What ran

`runtime/internal/benchmark` calls the existing pipeline. It does not add a stage, a browser, a judge, a retry system, or a memory store.

Each trace copies the fixture twice. The baseline copy gets `packet.json` with orchestration off and does not call `Execute`. The Orchestra copy declares the starting files on the existing provenance ledger, then calls `DesignPipeline.Execute`. The implement stage still records a handoff. It does not rewrite the fixture.

The paired host smoke does not call a model. Status is `unavailable`, phase is `pilot`, trial is `1`, claim is `none`, tokens are `unavailable`.

## Measured trace

Two trials per fixture, same hash, same capability ids, same byte sizes. Wall time is the pipeline clock and is not a quality score.

| Fixture | Intensity | Mission bytes | Capability bytes | Capsule bytes | Capability ids | Brain / directions / browser / research / judge / screenshots |
| --- | --- | --- | --- | --- | --- | --- |
| t1-api | LIGHT | 1912 | 266 | 969 | coding.worker | 0 |
| t3-debug | LIGHT | 1870 | 268 | 923 | coding.worker | 0 |

Orchestration bytes are the sum of those three files. Task-execution cost is `unavailable`. Tokens are `unavailable`. A missing Brain file did not fail the trace.

Wall time on the passing run was about 22–45 ms. That number moves between runs. The byte sizes did not.

## Eight questions

1. **Reproducible run.** Yes for the trace. Trial 0 and trial 1 of each fixture kept the same fixture hash, the same capability ids, and the same mission, capability, and capsule byte sizes.
2. **Isolation.** Yes. Both copies start from the same fixture hash and the same task and acceptance text. The baseline packet is rejected if it carries mission, capability, or capsule paths. After the Orchestra arm runs, the baseline tree excluding `packet.json` still matches the source hash.
3. **LIGHT stays cheap.** Yes on T1 and T3. Directions, browser calls, screenshot evaluations, research calls, judge calls, and Brain facts are 0. Intensity is LIGHT.
4. **Unnecessary capabilities avoided.** Yes for this fixture class. The selected id is `coding.worker`. Design, browser, research, and judge are not selected. `grok-4.7`, `claude-opus`, `gemini-2.5`, and `gpt-x` select the same ids, and those labels do not appear in the mission or capability files.
5. **Recovery.** Yes, on the existing 3C path. The seeded message `test failed: Sum(1,2) = 3, got 4` classifies as a test failure, retries once with a required change, and records that attempt. A later different failure sees one fewer remaining retry. The identical replay stops and asks for a human. The pipeline trace does not itself rewrite the buggy function. Implement still only records a handoff.
6. **Scoring resists override.** Yes. T3 tests fail, and a positive judgment sheet stays a hard fail. T1 tests pass, and a negative judgment sheet does not create a hard fail. There is no overall score.
7. **Context bytes.** Yes. Mission, capability, and capsule sizes are recorded in the table above. Tokens stay `unavailable` because no host returned a count.
8. **Expansion.** The gate function can return true when traces, isolation, scoring, and recovery are all true. This pass does not flip `enabled` on T2 or T4–T10. Chromium captures, a Brain-effectiveness host run, and a second trial were not run. Expansion is safe to consider later. It is not enabled here.

## Limitations

- The Orchestra arm is a trace. It is not a product diff.
- No host wrote either worktree. The comparison cell is `unavailable`, not a win.
- Starting fixture files are recorded on the existing provenance ledger before `Execute`, so the locked integrity check can see them as the starting tree. That is not a new acquisition system.
- T3 recovery is the 3C classifier and retry rule. It is not a new loop, and it is not evidence that the pipeline fixed `sum.go`.
- One pair of traces is a pilot. It cannot support a general improvement claim.

## Checks

From `runtime/`, `go test ./internal/... -count=1` passed, including `internal/benchmark` and `internal/visual`.

From the repo root, `node runtime/tools/check-host-stack.js` reported 44 skills and 0 problems. `node runtime/tools/hygiene.js` scanned 725 files, `VERSION` `3.3.1`, 0 failures.

No commit. No tag. No version bump.

## Stop

Do not run T5–T10. Do not start 3H. Do not treat this report as proof that Orchestra writes better software.

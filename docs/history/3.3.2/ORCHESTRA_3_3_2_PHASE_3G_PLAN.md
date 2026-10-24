# Orchestra 3.3.2 — Phase 3G pilot plan

This pass checks whether a benchmark can measure Orchestra. It does not claim that Orchestra improves output.

The later study, using the same schema, asks whether the same host, task, repository, and acceptance get better work at a justified cost. That study is not this pass.

VERSION stays `3.3.1`. No commit, push, tag, or 3H. No Portfolio, Gate 0, JCODE, OpenDesign, SQLite, DAG, swarms, second conductor, second mission, second retry, model registry, or memory platform.

Locked: 3A through 3F.1. Existing 3C–3E tests already cover the old hygiene list. 3G adds the comparison harness.

## Two measurements

`DesignPipeline.Execute` is an orchestration trace. The implement stage records a handoff. It does not write the product.

- **Trace.** The pipeline runs on a fixture copy. It measures intensity, overhead, gates, retries, and context bytes. A trace is not better software.
- **Paired output.** At most one host pair, status `pilot`. If no host runs, the cell is `unavailable`. No fabricated output and no fabricated token counts.

## Architecture

`runtime/internal/benchmark` owns the manifest, packets, scoring contract, integrity checks, and the result record. It calls the existing pipeline.

- Orchestra arm: `Execute` on a temp copy. Mission, budget, and the capability decision stay in `.orchestra/runs/<id>/`.
- Baseline arm: no `Execute`. The packet is the task, acceptance, fixture id, time budget, and human policy.
- Scorer: hard checks first. A judgment sheet cannot flip a hard fail. The scorer does not read the mode until after the checks.
- Brain: existing `LoadDurable` / `BindCapsule` on a temp ledger. No vault read.
- Recovery: existing `mission.Classify` and `mission.Decide`. No new retry system.
- Model label: metadata only. `grok-4.7` is an opaque label. Equivalent labels must select the same capability ids.

## Matrix

Fixtures live under `runtime/benchmarks/fixtures/`. Not production product repositories.

- **T1 API.** Enabled. LIGHT control. No lab, browser, screenshots, research, judge, or Brain retrieval.
- **T3 Debug.** Enabled. Seeded test failure and the existing recovery rule.
- **T2, T4–T10.** Present in the manifest with `enabled: false`.

## Integrity

Both arms start from the same fixture hash, task text, and acceptance. The baseline fails closed if it contains `.orchestra/`, `mission.json`, `capability.json`, a capsule, Brain facts, or a prior benchmark result.

## Scoring contract

Declared in code before a run:

- Hard fail: required behavior missing, build failure, test failure, acceptance failure.
- Judgment cannot override a hard fail.
- No single overall score.
- Tokens stay `unavailable` unless a host supplied a count. Byte sizes are context measurements.

## Recovery

T3's seeded failure is classified with the existing failure classes. The first decision retries only with a change and consumes retry budget. The identical failure does not retry again.

## Expansion

T5–T10, Chromium benchmark captures, Brain-effectiveness runs, and multi-trial comparisons stay disabled in this pass. The report says whether that gate may open later. It does not open it.

## Pilot questions

The report answers only:

- Does the benchmark run twice from the same fixture hash?
- Are the baseline and Orchestra arms isolated?
- Does LIGHT stay cheap on T1?
- Does the capability set omit unnecessary categories?
- Does recovery classify, change, and stop an identical loop?
- Can a positive judgment sheet override a failed test?
- Are mission, capability, and capsule bytes recorded?
- Is expansion to harder classes allowed, and was it left disabled?

## Do not build

A leaderboard, a token estimator, an in-process coder, a second browser, a multimodal evaluator, preference inference, a benchmark database, a stats package, or 3H.

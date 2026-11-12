---
name: ow-eval
description: Paired eval of a raw host, the control plane, and the routed skills. Three trials. A clean HOME for the raw arm. Promote a skill only when success is not negative and cost rises by 10 percent or less.
---

# Eval

Use this when a skill, route, or release needs a measured comparison. It extends `ow-benchmark`. It does not invent a winner.

## Preconditions

A fixture exists. Arm A gets a clean HOME. The same model label is used for A, B, and C.

## Arms

1. **A.** A clean `HOME` with no Orchestra files and no global skills. The prompt says to skip Orchestra.
2. **B.** The Orchestra control plane only. No specialist skill under test.
3. **C.** The control plane plus the skills the route selected.

Same model label, same fixture, same budget. Run A, B, and C in an order that rotates across trials. Keep failed and interrupted trials.

## Procedure

1. Read the fixture under `runtime/benchmarks/orchestra-47/fixtures/`.
2. Record task success, tests, Lighthouse when the fixture is visual, axe, steps, tokens, and time.
3. A blind judge scores visual fixtures from shuffled screenshots.
4. Call `PromoteSkill` in `runtime/internal/benchmark/promote.go`. Three trials is the minimum.
5. Write the trial JSON under `runtime/benchmarks/orchestra-47/<fixture>/`. Status `not_run` stays `not_run` when a model, a browser, or a judge did not run.

## Output

A trial JSON file. `not_run` stays `not_run` when a model, a browser, or a judge did not run.

## Verify

`PromoteSkill` is called only with the recorded trials. A claim stays none until the rule passes.

## Skip

Do not enable a comparison claim. Do not put a skill on a Load line from a single trial. Do not share the operator's skill folders with arm A.

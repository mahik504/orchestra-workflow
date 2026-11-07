---
name: ow-benchmark
description: Compare Orchestra to a raw agent on a fixture. Same model, three trials, Layer A when the page is visual. Record the numbers you measured. Leave the claim disabled when the run did not happen.
---

# Benchmark

Use this when the task is to measure Orchestra against another arm. The comparison claims in `runtime/internal/benchmark/compare.go` stay off until a record exists.

## Rules

1. Same model label on both arms. The raw arm is told to skip Orchestra. The Orchestra arm gets the contract.
2. Three trials. A single lucky run is not a result.
3. Visual tasks get desktop and mobile screenshots, axe, and Lighthouse when those tools actually run.
4. A blind judge sees shuffled screenshots and does not hear which arm made them.
5. Write the record under `runtime/benchmarks/<task>/`. Status `not_run` stays `not_run` if the agent, the browser, or the judge did not run.

## Do not

Do not invent a winner. Do not enable `c1` in the comparison suite without the record. Do not put private prompts or tokens in the record.

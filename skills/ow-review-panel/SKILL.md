---
name: ow-review-panel
description: Run two reviews in order. Correctness first, from a clean context. Simplify second, on the diff only. Findings need a file and a failure. No score, no style nits.
---

# Review panel

Use this at the end of a pass, after the checks you can run have run.

## Correctness

`ow-fresh-review` reads the result as if it did not write it. Ask what is wrong, what was not run, and which acceptance line is still open. A finding names a file and the input that fails.

## Simplify

A separate pass. `ponytail` on engineering diffs only. On a Design Lab surface, do not "simplify" the visual direction away. Cut unused code, duplicated helpers, and comments that restate the next line.

## Order

Correctness, then simplify. Do not mix them. Do not mark the work done inside the review. `verification-before-completion` still wants the command output.


## Preconditions

- The task names the files or the acceptance line.

## Procedure

1. The task's paths and acceptance line win.
2. Correctness first. Simplify only after that review names a failing input or says none.

## Output

The change, the command that checked it, and what that command printed.

## Verify

Say the work is done only after that command has been run.

## Skip

Skip when the human says skip orchestra. Do not add a second conductor.

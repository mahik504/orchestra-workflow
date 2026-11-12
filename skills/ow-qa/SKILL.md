---
name: ow-qa
description: Turn a bug on the running app into a failing browser check, then a fix, then a regression test. Playwright is the evidence.
---

# QA

Extends `webapp-testing`. One browser owner.

## Preconditions

- The app runs locally or at a URL you are allowed to test.
- The bug is a behavior you can reproduce, not a guess from the source.

## Procedure

1. Reproduce the bug in Playwright. Keep the failing assertion.
2. Fix the product code.
3. Re-run the same test and the nearby flows that share the state.
4. Leave the test in the suite.

## Output

A regression test and the command output from the passing run.

## Verify

The test fails before the fix and passes after it. Say so only if you saw both.

## Skip

Do not call the bug fixed from a source reading. Do not add a second browser harness.

---
name: ow-semgrep
description: Run Semgrep against the product repo you own and report the rules that fired. Do not scan a third-party URL.
---

# Semgrep

## Preconditions

- The repo is yours, a client's, a CTF, or open source.
- `semgrep` is on PATH. If it is not, say so and stop. Do not install a global copy as part of the task.

## Procedure

1. From the product root, run `semgrep scan --config auto --error` when a ruleset is already configured. Otherwise run `semgrep scan --config p/default`.
2. Write the text output to the verification note. Do not hide a finding.

## Output

The command, the exit code, and the rule ids that fired.

## Verify

The exit code is in the same message as any claim that the scan passed.

## Skip

Do not point Semgrep at a site you do not own. Do not add `--no-git-ignore` to sweep secrets into the report and then commit that report.

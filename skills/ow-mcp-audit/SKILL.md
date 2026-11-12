---
name: ow-mcp-audit
description: Scan local MCP configs and record a baseline of tool names. Snyk runs only when SNYK_TOKEN is already set. Do not print secrets.
---

# MCP audit

## Preconditions

- You are auditing this operator's configs or a project they own.
- AgentShield can be run as `npx ecc-agentshield scan` without a global install. Skip it if the package is not wanted. The config lint in `runtime/tools/check-config-lint.mjs` still runs.

## Procedure

1. Run `node runtime/tools/check-config-lint.mjs`.
2. List server names from the doctor. Do not print env values.
3. If `SNYK_TOKEN` is set, an optional Snyk agent scan may run. It uploads tool descriptions. Do not run it when the token is absent.
4. Save a baseline JSON of server and tool names under `.orchestra/mcp-baseline.json` in the product repo.

## Output

The lint result and the baseline. A later diff shows a tool description that changed.

## Verify

The baseline contains names only. A search of the file finds no key material.

## Skip

Do not upload descriptions without a token. Do not install ECC as a conductor.

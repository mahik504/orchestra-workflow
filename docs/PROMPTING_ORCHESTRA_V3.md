# Prompting Orchestra 3.1

A brief still needs intent, context, constraints, and what done looks like. Orchestra 3.1 classifies that brief against `registries/design-resource-graph.json`. It does not invent a second prompting language.

Read [`AGENTS.md`](../AGENTS.md) and [`WORKFLOW.md`](../WORKFLOW.md).

- Name the product type (restaurant site, school SaaS, 3D portfolio) so the classifier can pick a route.
- For showable UI, do not write frontend files until Design Lab approval — or say `skip the lab`.
- Say **skip orchestra** to stand the contract down for the session.

Host adapters (`.cursorrules`, `CLAUDE.md`) translate syntax. They do not replace the re-brief or the gate.

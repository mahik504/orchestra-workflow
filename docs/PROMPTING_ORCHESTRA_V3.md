# Prompting Orchestra 3.3

A brief still needs intent, context, constraints, and what done looks like. Orchestra 3.3 classifies that brief against `registries/design-resource-graph.json`. It does not invent a second prompting language.

Read [`AGENTS.md`](../AGENTS.md) and [`WORKFLOW.md`](../WORKFLOW.md).

- Name the product type (restaurant site, school SaaS, 3D portfolio) so the classifier can pick a route.
- For showable UI, Design Lab is 23 short cards then one `DESIGN.md` — or say `skip the lab`. If you name a site, skill, MCP, pack, or paste a `DESIGN.md`, that outranks the graph.
- Say **skip orchestra** to stand the contract down for the session.

Host adapters (`.cursorrules`, `CLAUDE.md`) translate syntax. They do not replace the re-brief or the gate.

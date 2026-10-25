---
name: graphify
description: Build a local code graph for one repository. Use when the user asks how a codebase fits together. Do not index a notes vault. Do not use Graphify Cloud.
---

# Graphify

Install the CLI with `uv tool install graphifyy` when it is missing. Run it on the product repo.

| Rule | Detail |
| --- | --- |
| Scope | One repository. Not a notes vault |
| Network | No Graphify Cloud account and no upload |
| Output | A folder that is gitignored |
| Check | Answers must match the source. If they do not, say so |

If the CLI is missing, say so and stop. Do not invent a graph.

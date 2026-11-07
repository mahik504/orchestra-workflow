# Upstream admission

Read on 2026-11-08. A missing SPDX entry means the folder was not copied.

| Source | Commit | License | Decision |
| --- | --- | --- | --- |
| `vercel-labs/skills` `skills/find-skills` | `d3701c11b1101251af14c9d148f026e1a47ebe17` (2026-11-06) | MIT | Vendored at `skills/find-skills`. `ow-skill-advisor` blocks `npx skills add -g -y` unless the human insists. |
| `vercel-labs/agent-skills` | GitHub license API returned 404 | none on the API | Not copied. `ow-deploy` is an original skill and says so. `react-best-practices` and `web-design-guidelines` were already in the tree from earlier work. |
| `anthropics/skills` | `2c5958ecac7cee2b583967a2c4981d0f4a75c579` (2026-11-06) | license API 404 | Not copied. |
| `cloudflare/skills` | `6eb6c781aa9b6438e0c7307bd102788a75fe1a6d` | Apache-2.0 | Not copied. The tree is Cloudflare product skills (`workers`, `wrangler`, `turnstile`). There is no `security-audit` skill in that listing. `ow-ship-safe` and `ow-semgrep` stay the security route. |

Other names from the brief (ui-skills, juliobsz, anatomy, context-engineering, diagram-design) were not given a repository URL with a license file in this pass. They are not installed. `ow-skill-advisor` is the path when one of them gets a URL.

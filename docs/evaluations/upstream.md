# Upstream admission

Read on 2026-11-08. A missing SPDX entry means the folder was not copied.

| Source | Commit | License | Decision |
| --- | --- | --- | --- |
| `vercel-labs/skills` `skills/find-skills` | `d3701c11b1101251af14c9d148f026e1a47ebe17` (2026-11-06) | MIT | Vendored at `skills/find-skills`. `ow-skill-advisor` blocks `npx skills add -g -y` unless the human insists. |
| `vercel-labs/agent-skills` | GitHub license API returned 404 | none on the API | Not copied. `ow-deploy` is an original skill and says so. `react-best-practices` and `web-design-guidelines` were already in the tree from earlier work. |
| `anthropics/skills` `webapp-testing`, `mcp-builder` | `2c5958ecac7cee2b583967a2c4981d0f4a75c579` (2026-11-06) | Apache-2.0 in each `LICENSE.txt` | Vendored. The repo license API returned 404. The skill files carry the license. Other skills in that repo were not copied. |
| `cloudflare/skills` | `6eb6c781aa9b6438e0c7307bd102788a75fe1a6d` | Apache-2.0 | Not copied. The tree is Cloudflare product skills (`workers`, `wrangler`, `turnstile`). There is no `security-audit` skill in that listing. `ow-ship-safe` and `ow-semgrep` stay the security route. |

Other names from the brief (ui-skills, juliobsz, context-engineering, diagram-design) were not given a repository URL with a license file in the 4.5 pass. They are not installed. `anatomy` was vendored in 4.6.0. `ow-skill-advisor` is the path when another name gets a URL.

## Admitted in 4.6.0

Read on 2026-11-08. Each tree was scanned for a remote shell pipe, a private key, and a call that sends secrets off the machine. Local Node scripts that read a repo or open a page the operator names stayed. None of these skills install themselves.

| Source | Commit | License | Decision |
| --- | --- | --- | --- |
| `heydamir/unslop` | `a6f121eee5f29c2f74976b75456feae8b66afb3a` | CC-BY-4.0 | Vendored at `skills/unslop`. The Damir Sokolovsky credit stays in `SKILL.md`. `snapshot.mjs` saves a page the operator names. It does not read API keys. |
| `MengTo/Skills` `build-awwwards-quality-sites` | `5ea1661a0e05d82b054553daaa6baf6d0233c33c` | MIT | Vendored under that name. The rest of the MengTo pack was not copied. |
| `haydenschmitty/floorplan` | `70be3f40f9c2e061290b76027f180d9ff892edd6` | MIT | Vendored with its scripts and the ELK, Inter, and JetBrains Mono license files. The scanner looks for secrets in the target repo and prints findings. It does not upload them. |
| `wheresryan22/anatomy` | `1482e8469b0a9e0d8cc18ea2740a01423d29ae80` | MIT | Vendored. Generated A4 PDFs and built HTML pages were left out. `node examples/<name>/build.mjs` writes them. `publish.mjs` can push a branch when the operator passes `--submit`. Orchestra does not run that during sync. |
| `google/elephant-goldfish-model` | `e9e1f06aa4e26e95cd4ed37210a0632202c429e9` | Apache-2.0 | Vendored `egm-draft`, `egm-review`, `egm-resolve`, and `egm-review-loop`. Reviewer prompts and document standards sit inside each skill folder. Changes are listed in each skill's `NOTICE`. |
| `google/elephant-goldfish-model` `egm-brainstorm`, `egm-bootstrap`, `egm-drift-report` | same commit | Apache-2.0 | Not copied. Brainstorm overlaps `grilling` and `ow-brainstorm`. Bootstrap and drift overlap `floorplan`. |
| `Leonxlnx/taste-skill` | `f6b76575f949707b4e3c37cd86cb44b0b98643e8` | MIT | Evaluated, not installed. |
| `trailofbits/skills` | license CC-BY-SA-4.0 | CC-BY-SA-4.0 | Evaluated, not installed. |

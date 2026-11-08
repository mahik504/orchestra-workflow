# Audit 4.6.0

Written after tag `v4.6.0`, which points at `6d05e43`. This file is the 4.6 audit. The other markdown files in this folder are the 4.5 baseline and stay as history.

## What it does now

Orchestra is the control plane. A task loads one route from `registries/routing.json`. Skills, MCPs, plugins, and libraries are capabilities. The vault is memory. The registries are resource knowledge.

4.6.0 vendors spec review from Google's Elephant-Goldfish Model at `e9e1f06aa4e26e95cd4ed37210a0632202c429e9`: `egm-draft`, `egm-review`, `egm-resolve`, and `egm-review-loop`. `egm-brainstorm`, `egm-bootstrap`, and `egm-drift-report` were left out. The same release vendors `unslop` (`a6f121eee5f29c2f74976b75456feae8b66afb3a`, CC-BY-4.0), `build-awwwards-quality-sites` (`5ea1661a0e05d82b054553daaa6baf6d0233c33c`, MIT), `floorplan` (`70be3f40f9c2e061290b76027f180d9ff892edd6`, MIT), `anatomy` (`1482e8469b0a9e0d8cc18ea2740a01423d29ae80`, MIT), and `webapp-testing` plus `mcp-builder` from `anthropics/skills` at `2c5958ecac7cee2b583967a2c4981d0f4a75c579` (Apache-2.0). `registries/skills.json` lists 99 skills. The 4.5 tag listed 89.

Sync moves a retired skill name into a backup. A tree whose hash matches `v4.0.0`, `v4.1.0`, `v4.5.0`, or commit `60cf9e7` is replaced. A file that already looks like Orchestra gets the 4.6 adapter. Plugin manifests sit at the repo root. Mermaid fences in the README and `docs/` are parsed in CI. The logo is gone.

## What was checked

| Check | Result |
| --- | --- |
| CI on `6d05e43` | success, run [42466472744](https://github.com/mahik504/orchestra-workflow/actions/runs/42466472744). The Go test step passed. |
| Hygiene on `6d05e43` | success, run [42466472903](https://github.com/mahik504/orchestra-workflow/actions/runs/42466472903) |
| CI on `a667b67` | failed. The node job died in Mermaid parse before Chromium was installed. Run 42465954536. |
| CI on `01fae42` | failed. The Go test step died: premium-website selection and `capabilities = 32, want 26`. Run 42466087169. Node and Windows jobs passed, including Mermaid parse. |
| Local `go test -count=1` on the classifier and compile packages after `6d05e43` | 84 passed |
| `node runtime/tools/check-route-agreement.js` | skills 35, routed 35, problems 0 |
| `node runtime/tools/enrich-graph.js` dry run | capabilities 33, problems 0 |
| Landing record | `runtime/benchmarks/c1-premium-landing.json`, 2026-11-08, `composer-2.5-fast`, claim none, enabled false |

An earlier full `go test ./...` from `runtime/` printed 770 passed. That run reused a cached classifier result. The graph JSON had changed and the Go files had not, so the cache hid the failure CI later printed. The green evidence for the tag is run 42466472744.

The host apply and the doctor ran before the tag. The doctor printed problems 0. Limitless Claude 8.2.1 printed Orchestra 4.6.0.

Two reviews ran before the tag. The correctness review ([fresh correctness review](e59fa9e2-bdae-4adf-acc4-f7c496c21aaa)) returned Rework. The release fixes were the landing keyword, the missing graph rows, the EXDEV copy, the adapter version regex, and the doctor treating a personal Gemini file as Orchestra. The simplify review ([simplify review](9ac5d36e-7c2b-410d-82e1-e1b38ed727b7)) asked to drop the Claude MCP add, to stop retiring `creative-3d-web` and `fullstack-mastery`, and to narrow the `webapp-testing` line. Those three stayed, because the release plan requires them.

## Strengths

The tagged commit is the one CI and Hygiene passed. The two red runs are still on `main` as history, and the tag was not moved back onto them.

The landing sample is written down with the limits in the same file. Desktop Lighthouse performance was 100, 99, and 99 on the raw pages and 99, 98, and 99 on the Orchestra pages. Accessibility was 100, 96, and 100 raw, and 100, 100, and 100 on Orchestra. Axe found `landmark-one-main` and `region` on the second raw page and no violations on the Orchestra pages. Layer A desktop and mobile were captured. A blind judge called trial 1 a tie and preferred Orchestra on trials 2 and 3. The judge was a separate agent, [blind landing judge](7d80bc54-ba74-4507-b0f1-cb71f8910aa2). `compare.go` stays disabled.

Spec review has a place in the loop, and a STANDARD fix that names the file skips it. Retired names are moved, which is how `creative-3d-web` and `fullstack-mastery` left the host skill folders on this machine.

## Flaws

Both landing arms could see the global skills. The Orchestra arm had `AGENTS.md`, `.cursorrules`, and `compile.json`. The raw arm was told to skip Orchestra and had no `AGENTS.md`. That is a weaker split than a clean environment. One Orchestra attempt wrote no files and was rerun. That rerun is trial 3. Two Orchestra agents were stopped after the production build was on disk. The claim stays none.

`go test` without `-count=1` does not notice a JSON-only graph edit. CI did. The first push also failed Mermaid parse until the workflow installed Playwright Chromium.

`defaultBackupRoot` in `runtime/tools/ow-sync.mjs` uses `%LOCALAPPDATA%` on Windows before the home argument. A caller that sets `OW_SYNC_HOME` to a temp directory still writes backups into the real profile unless `OW_BACKUP_ROOT` is set. The real-home apply wanted that profile. A test that passes a temp home does not get a temp backup.

`runtime/tools/enrich-graph.js` sets `graph.version` to `4.0.0` on every write. The graph field stays 4.0.0 while `VERSION` is 4.6.0. `check-versions` allows that.

Java was not on PATH when 4.5 was recorded, so `rea` stayed unloaded. This release did not install it.

Cursor User Rules live in Settings. Sync does not write them. Confirming that the Superpowers and caveman plugins are off is still a click in that Settings page.

## Other harnesses

The readings and the commits opened on 2026-11-08 are in `docs/evaluations/harnesses.md`. Orchestra stayed the conductor.

Superpowers skills are already vendored at `22f0b463e7a5e579904dd6211a520321035526f7`. The plugin stays off because it injects a second plan. Pocock skills at `4ae2317ddb22b046e26455e4f9c9b3aaaf62212d` contributed `grilling`, `domain-modeling`, and `to-tickets`. `constraint-driven-development` came from `addyosmani/agent-skills` at `301b02b12bff98e4688e6984af6556802322bbde`. Spec Kit at `c43849408ac0de2ca6a52f14eaae71449776057b` is the acceptance line `ow-forge` already writes. GSD at `bdcaab2c752d9a33a1a1ca9acf3a3c81fb991815` is the receipt and `orchestra compile`. BMAD at `e734d3ac3c0f4ea31cc2c2bbbe79e86e54f2cae6` is the two reviewer files in `agents/`. ECC (`d2ef2680594f9a17526a6e1a2e18b3d242c5d17b`), gstack (`1a29fbc3444f858c72ec001abdb4f0f3bebfa76c`), and Agency Agents (`b859c515d813efcfc4755d3471dd9e7686c7706d`) were not installed. Karpathy-style guidelines are the short-context rules already in `AGENTS.md` and `docs/CONTEXT-BUDGET.md`. There is no skill named after a person. EGM is the vendored spec-review set above, with the three overlapping skills left out.

## Against earlier releases

| Release | Commit | What changed in the method |
| --- | --- | --- |
| 4.0.0 | `9675bfd` | Routes select skill folders. Changelog records 73 canonical skills. No website study. |
| 4.1.0 | `61a8c71` | The routing note names only what a task loads. Same 73 skills. |
| 4.2.0 | `65d1720` in the changelog. No tag. | Grilling, domain modeling, tickets, constraint-driven development, and receiving-code-review. Changelog records 78 skills. |
| 4.5.0 | `f4208a3` | Orchestra-owned skills are `ow-*`. `skills.json` on that tag lists 89. Sync leaves a local edit. Landing and bugfix samples recorded. Claim none. |
| 4.6.0 | `6d05e43` | Spec review, six upstream skills, sync backups, plugin manifests, Mermaid parse, logo removed, landing sample rerun with `AGENTS.md` on the Orchestra arm. 99 skills. Claim none. |

4.0 and 4.1 have no `registries/skills.json` on the tag, so those counts are the changelog lines. 4.5 and 4.6 are counts of the `id` rows in `skills.json` at the tag.

## Rating

7 out of 10.

The release is tagged on a green CI run, the skill list is one list, and the landing sample includes the scores, the axe rule ids, the judge, and the shared-skills limit. That is a method someone else can install and check.

It stays at 7 because the comparison claim is off for a real reason: the arms shared global skills, one Orchestra attempt had to be rerun, and two agents were stopped after the build. Two CI failures also landed on `main` before the green commit. A higher score needs a comparison where the raw arm cannot see the Orchestra skills, and a local test run that fails when the graph JSON changes.

## Next

Make a graph-reading Go test depend on the JSON, so a cached `go test` fails when the capability count changes. Give the c1 raw arm a workspace that does not load the global skill folders. On Windows, write backups under `OW_SYNC_HOME` when that variable is set. Install Java before `rea`. Leave `compare.go` disabled until a new record says otherwise.

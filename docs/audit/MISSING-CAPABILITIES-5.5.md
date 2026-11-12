# Missing capabilities 5.5

Recorded 2026-11-12. “Missing” means the route table cannot name a procedure that exists, or the procedure exists and is a stub, or a document claims a route the table does not have. It does not mean “add a skill.”

## Fix the map before adding names

| Gap | Evidence | Action |
| --- | --- | --- |
| `ow-qa`, `ow-browserize`, and `ow-author` are described as routes and are not rows | Search of `registries/routing.json` found no matches. Changelog and `EVAL-5.0.md` say they are routes | Either add a row with a real match list and a skip line, or correct the changelog. Do not add a row just to make the sentence true. `ow-qa` is the one that earns a row: browser bugs are a repeated job and the skill has a procedure |
| `ow-deploy` and `ow-store-release` are not what the ship and mobile routes load | Route table versus skill folders | Leave them off the default load until they have a procedure and a trial. The ship route already names `ow-ship-safe` |
| Doc skills are both “not in skills/” and present, with a license that forbids distribution | `docs/ROUTING.md` line for docs on disk, and `skills/docx/LICENSE.txt` | Remove the trees from the public repo when edits are allowed. Keep the route as host-provided |
| `creative-3d-web` and `fullstack-mastery` are full skill bodies sitting next to the skills that replaced them | Files opened. Off `host-stack.json` | Move them out of `skills/` so a model searching the tree does not load the second stack. Not a new capability |

## Procedures that are missing inside skills that already exist

The 36 `stub: true` skills have no `## Procedure`. The ones that block real work:

| Skill | Why it is first |
| --- | --- |
| `ow-conductor` | Every full-path task is told to use it, and it points at a 19 KB 4.6.0 document |
| `ow-forge` | The brief is step 0 |
| `ow-backend-architecture` | The backend route loads it |
| `ow-vault` | Memory writes have a policy in Go and a stub in the skill |
| `ow-design-intelligence` | Frontend direction has no procedure in that file |
| `ow-research` and `ow-research-paper` | Citation and “what people said” jobs |
| `ow-ml-experiment` | The skip line is the whole policy |
| `ow-ship-safe` and `ow-observability` | The production route loads them |

Write the procedure so the task’s paths win, in the shape `ow-qa` already uses: Preconditions, Procedure, Output, Verify, Skip. Then drop `stub: true`. Do not paste a textbook.

## Disciplines with no route, and the decision

| Discipline | Decision |
| --- | --- |
| Distributed systems, microservices, caching, cloud cost | Fold into `ow-backend-architecture` as a short checklist. No new skill until a task fails without it |
| Authn and authz | Same. The backend checklist should name the boundary. A dedicated skill is a later trial |
| Fine-tuning, MLOps, vision, NLP | Stay inside `ow-ml-experiment` plus the data the user already has. No training skill that invents a metric |
| Literature review and statistics | `ow-research-paper`. Cite what was opened. No DOI invention |
| Slide QA | `ow-presentation` procedure: hierarchy, notes, and a read-through. Not a new deck skill |
| Desktop automation | Do not add. The host and the hook are the boundary |
| Model routing across providers | Do not add. The host selects the model. `model_class` stays a hint |
| A skill per component library | Reject. `ow-component-discovery` plus the registry |

## Two selectors, and one name that is not on a route

`orchestra compile` scores `registries/routing.json`. `orchestra classify` scores `design-resource-graph.json`. `ARCHITECTURE.md` still says the graph has 17 rows. The file has 33 capabilities. A chat that follows the markdown table and a process that follows `classify` can load different skills. 5.5 should make the route table the skill picker and the graph a resource catalog, which is the simplification already in `RESOURCE-SELECTION-5.5.md`.

`ow-focus` is not in `routing.json`. `docs/CONDUCTOR.md` line 61 still tells the conductor to run it before classify. `STUB-TRIAGE.md` says the route table needs the name. The table does not. Keep the skill only if the conductor procedure still calls it. Do not cite the route table as the reason.

`WORKFLOW.md` documents `orchestra classify --task`. The binary has no `--task` flag. Positional text works. The documented command exits. That is F-15 in the baseline.

## What is not missing

TDD, systematic debugging, verification before completion, worktrees, code-review reception, grilling, domain modeling, and to-tickets are on disk. Award-style, unslop, 3D, scroll, mobile Expo, scrape, site recon, and app recon have routes. The gap is procedure text and honest verification, not another catalog.

## Eval gap

No fixture has a live model result for 5.0.0. That is the missing capability that blocks a superiority sentence. It is not a missing skill. See `EVAL-5.5.md`.

# Design Lab protocol

The Design Lab is the checkpoint between "I understood the brief" and "I started writing frontend files". Its job is to make *design-as-you-code* impossible for work a stranger will see.

The gate is a lock, not a warning. While it is `PENDING`, the engine refuses to write anything a browser renders.

## When the lab runs

| Quality bar | Design Lab |
| --- | --- |
| `STANDARD` | **Off** by default. The human can opt in by asking. |
| `PREMIUM` | **On**. Opt out only if the human says "skip the lab". |
| `EXPERIMENTAL` | **On**. Always ship a low-end fallback. |

The bar comes from the capability row in `registries/design-resource-graph.json` (`quality_bar`), then from the brief: internal or throwaway work drops to `STANDARD`, and work a stranger will see rises to `PREMIUM`.

Run `orchestra classify "<your brief>"` to see the bar, the chosen route, and every route that was declined.

## Gate states

| State | Meaning | Frontend writes |
| --- | --- | --- |
| `NOT_REQUIRED` | Bar or task does not call for a lab | allowed |
| `PENDING` | Directions owed, none approved | **blocked** |
| `APPROVED` | A named direction was approved by a named person | allowed |
| `BYPASSED` | The human waived the lab, with a recorded note | allowed |

A bypass is legitimate. A *silent* bypass is not — `Bypass` refuses an empty note.

## What is blocked

Anything the browser renders: `.css`, `.scss`, `.sass`, `.less`, `.html`, `.jsx`, `.tsx`, `.vue`, `.svelte`, `.astro`, `.glsl`, `.frag`, `.vert`, plus token and theme files (`tailwind.config.*`, `theme.*`, `tokens.*`, `globals.*`, `design-system.*`).

Backend code, migrations, notes, and `DESIGN.md` itself stay writable. The human needs something to read before they can approve anything.

## What a direction must contain

Produce **2 or 3** directions. One is a decree; four is a survey.

Every direction needs:

- Visual concept and product type
- Typography — a named pairing **and where it came from**
- Colour world — **and its source**
- Layout language
- Component kit (named, or `custom`)
- **One** motion engine, **and why that one**
- 3D: yes/no, and the library if yes
- Shader: yes/no
- Logo method
- Icon system
- Implementation stack

Typography, colour, and motion are refused without a named source. A claim with no source is a vibe, and vibes are what produce the same purple-gradient page every time.

## Rejection is the useful half

When the human turns a direction down, record **their stated reason**. `Reject` refuses an empty reason.

Rejections are fingerprinted by their actual stack — typography, colour world, layout, component kit, motion engine, 3D, shader, and the sorted implementation stack — not by their name. Renaming "Warm editorial" to "Warm editorial, take two" and re-offering it is refused. Changing the stack is not.

The log lives at `.orchestra/design-lab/rejected-directions.json` and persists across tasks in the same workspace, so the second pass genuinely starts from where the first one ended.

Approvals are recorded at `.orchestra/design-lab/approved-<task-id>.json` with the direction and the approver's name.

## Overrides

- **"skip the lab"** in the brief, or `--skip-visual-gate`, waives the lab for one task.
- **"skip orchestra"** stands the whole contract down for the session.

The gate is a checkpoint, not a cage. The human can override the stack at any point after approval.

## Anti-slop

The lab exists because of a specific failure: research → compare → synthesize is skipped, and the model reaches for its priors. Its priors are Inter with a purple gradient, three equal cards, unmotivated glass, badge soup, motion with no purpose, and `#000000` with no chromatic depth.

Extract principles from references. Never copy branding, assets, copy, trademarks, or source.

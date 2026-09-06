# Designer SOS protocol

Design decision protocol from the Designer SOS Kit (ChatGPT dossier, 2026-09-06). Diagnose the stall **before** opening a tool. One stall, one tool, one timer.

This is evaluation language. It is not a 31st skill and not a dump of the gated PDF.

## S.A.V.E.

1. **Say** the stall in one sentence (what is blocked, on which surface).
2. **Assign** one tool for that stall (Taste, Impeccable, Playwright, GetLayers, a named reference — not all of them).
3. **Validate** in the **real project**, not inside the generator.
4. **Exit** when the stall is resolved or parked. Do not open a second tool “just in case.”

## Keep / Kill / Park

Finish the stall with one of:

- **Keep** — it serves the approved stack; leave it.
- **Kill** — it fights the stack or is generic slop; remove it.
- **Park** — useful later, not this story; note it and stop touching it.

## Six-pass teardown

On showable UI, after Playwright stills exist, run these passes in order. One finding per pass is enough.

1. **Squint** — does one focal point survive with the type unreadable?
2. **Grid** — alignment, columns, accidental equal card grids.
3. **Type** — pairing vs `DESIGN.md`; default Inter is a fail on PREMIUM.
4. **Colour** — 60/30/10 role thinking; token ramps, not a single hex; one background system per project; WCAG contrast as a number.
5. **Space** — density vs the scale in `DESIGN.md`; rasterize heavy blurs if they cost frames.
6. **Cut** — what can leave without losing the story.

A 20-minute rescue sprint is allowed when a ship is blocked. Same six passes, tighter timer.

## Vault / Live / Anatomy

Keep three folders of notes, not three visual worlds:

| Folder | Holds |
| --- | --- |
| **Vault** | Rejected directions, parked echoes, SOS findings |
| **Live** | The approved `DESIGN.md` and tokens in the app |
| **Anatomy** | Extracted principles from a named reference (not copied assets) |

Reference extraction instead of visual copying. Never copy branding, assets, copy, trademarks, or source.

## Where it sits

After Taste → `DESIGN.md` → implement → Impeccable → Playwright. SOS is the teardown. Then **one** consultant packet if the bar is PREMIUM/EXPERIMENTAL. Details: `VISUAL_QA_PROTOCOL.md`.

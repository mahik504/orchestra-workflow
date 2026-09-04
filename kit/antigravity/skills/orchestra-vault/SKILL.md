---
name: orchestra-vault
description: Maintain the private Orchestra workspace and its memory. Jump one note. Lasting files only. Record real outcomes, never synthetic ones. No secrets.
---

# Orchestra vault — 3.1.0

**BRAIN = MEMORY.** The workspace is where Orchestra remembers; it is not a second control plane.

Root: the private workspace created by `kit/init-workspace`, or the path named in `WORKFLOW.md`. Not this public repo.

## Read

Do not open the whole workspace. `routes.md` → one file → the app repo. Catalog files only when asked to catalog.

## Layout

| Path | What |
| --- | --- |
| `projects/<slug>/` | One product. Starts with `idea.md`. Add design, spec, progress, review, ship-post only when they exist. |
| `templates/` | Copy, then fill |
| `memory/decisions.md` | Dated architecture and product decisions |
| `memory/resource-memory.json` | Real resource outcomes only |
| `Preferences.md` | Taste, stack, how the human thinks |
| `kit/` | Host adapters and sync scripts |

## What memory is for

Three layers:

- **Global** — taste and stack in `Preferences.md`; dated yes/no in `memory/decisions.md`.
- **Project** — one product's context in `projects/<slug>/idea.md`.
- **Session** — distilled at the end. Discarded, not filed.

Record after a real job: which resource combination improved the result, what failed, what the human liked or hated, and why a route was chosen.

**Never** write synthetic evaluations. Generated or duplicated rows in `resource-memory.json` are not learning — they are noise that will mislead future routing. Only an executed job writes a row.

Do not file chat transcripts. Do not restate facts the repository already shows.

## Write hygiene

1. Will this matter in a month? If not, say it in chat and do not file it.
2. Prefer editing the existing note over creating a new one.
3. One idea per note, self-contained, headed so a later reader can split it.
4. A map points; it does not duplicate.
5. Delete junk the same turn you notice it.
6. Secrets never here. Never commit MCP configs, `.env`, or tokens.

## Learning

After a like, a hate, a finished task, or a genuinely useful find, append to `Preferences.md` in the **same turn**. Taste is allowed to change; prefer editing the existing bullet. Do not create a second taste file.

## Git

The workspace may have a **private** remote. This public template is a different repo. Do not flip a populated workspace public — its history contains personal files.

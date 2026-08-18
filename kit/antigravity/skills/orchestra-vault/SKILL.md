---
name: orchestra-vault
description: Maintains the Obsidian vault at C:\projects\orchestra-brain. Per-project folders. File only lasting notes. Delete junk the same turn. Learn taste and thinking into Preferences.md.
---

# Orchestra vault

Root: `C:\projects\orchestra-brain`

## Read (save tokens)

Do **not** open the whole brain. Jump `routes.md` → one file → the repo. `STACK.md` only if asked to catalog. `memory/career.md` only for hiring.

## Layout

| Path | What |
| --- | --- |
| `projects/<slug>/` | That product only. Start with `idea.md` (includes **this product's stack** and **kind**). Add architecture/design/research/packet/review/ship-post **only when they exist**. |
| `templates/` | Copy, then fill |
| `memory/` | Global decisions. Career is on this PC + the **private** remote. Public template has `memory/README.md` + `decisions.md`. |
| `memory/local-notes.md` | Gitignored diary nits. Not the public template. Not a second taste file. |
| `kit/sync-vault.ps1` | Private git backup (12h) of this vault. |
| `kit/sync-both.ps1` | 12h: private backup, then public template if files changed. |
| `routes.md` | Keyword → one note. Do not scan every idea.md. |
| `STACK.md` | Installed skills, MCP names, GitHub we have, refused tools. |

No `00-inbox`. No numbered dump folders. No chat recaps. No trending-tool lists.

## Kind (on every idea.md)

`college` = assignments / due-now. `personal` = studio or parked. `hiring-cv` = public internship work.

## Showable titles, not structure

Polish H1s and YAML `title:` so the vault looks professional when shown. **Do not rename** `memory`, `projects/`, `templates/`, `kit/`, or `C:\projects\orchestra-brain`.

## What belongs where

| Kind | Where |
| --- | --- |
| Product idea, stack, open questions | `projects/<slug>/idea.md` |
| Global taste, interests, how they think | `Preferences.md` |
| Lasting yes/no we already decided | `memory/decisions.md` |
| Internship / GitHub / LinkedIn | `memory/career.md` (private backup; omit from public export) |
| Assignment / viva diary | `memory/local-notes.md` (gitignored) |
| Tool we actually adopted | one line in Preferences Stack |
| Chat, “maybe later” tools, session recap | **do not file** |

## Hygiene (every write)

1. Will this still matter in a month? If no, do not file. Say it in chat.
2. Prefer editing `idea.md` / Preferences over a new file.
3. Answer first. One idea per note. Self-contained. Headings so a later reader (or retrieval) can split it.
4. A map points; it does not duplicate. Do not copy pricing, eval numbers, or stack into `START HERE.md`.
5. After a ship, delete stale packets and duplicate drafts in that project folder.
6. If you find empty numbered folders (`00-inbox` … `07-reviews`) or paths to them, delete the folders and fix the path.
7. Secrets never here. Never invent facts that are not in the repo or a note. Never commit `mcp_config.json`, `.env`, or tokens.

## Learn (reinforcement, not ML)

After they love or hate an output, a finished task, **or a useful find**, append **Liked** / **Hated** / **Thinking** in `Preferences.md` the same turn. Taste may change. Prefer editing the existing bullets. Do not file chats.

Do not create a second taste file. `memory/local-notes.md` is only for private diary nits that must not go public.

## Keep updated

When a decision or ship happens, update that product’s `idea.md`, `memory/decisions.md`, and the START HERE state table the same turn. Unpublished `idea.md` still updates on disk even if the public template omits that folder.

## Remotes

- **Private** `mahik504/orchestra-brain`: full backup, including career and unpublished ideas. 12h sync. Do not gitignore career.md here. Do not make this repo public (history leak).
- **Public** `mahik504/orchestra-workflow`: allowlisted template. 12h publish only when those files changed. No empty green-square commits.

## Internet

Search when it will actually improve the work. Not by default-for-fun. Not trending lists in the vault. Adopt into Preferences if we will use it.

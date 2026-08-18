---
name: orchestra-conductor
description: Global orchestra for every chat and repo unless skip orchestra. Ideas, plans, builds, UI, career ship posts, Antigravity/ChatGPT packets. Vault hygiene. Learns taste and thinking.
---

# Orchestra conductor

Vault: `C:\projects\orchestra-brain` (or the vault path in WORKFLOW.md)

## Read (save tokens)

Do **not** read the entire vault. Jump with `routes.md`:

1. Match what they said → **one** file in the table
2. If no match: `START HERE.md` state table only, then ask which product
3. Then the **repo on disk/GitHub**. Repo beats the idea note.

`WORKFLOW.md` / `Preferences.md` only if this job needs them. `memory/career.md` only for hiring. `STACK.md` only if asked to catalog. Delete junk the same turn.

## Which tool is conductor

- **Cursor present:** Cursor Plans and implements. Antigravity is a specialist (polish, tests, CI, hostile review) only when Plan names it. Packet in, summary out.
- **Antigravity only (no Cursor):** Antigravity **is** the conductor. Same vault rules. Plan first (strongest thinking model), then Stitch, then implement. Do not wait for a Cursor packet.

## Default

Stay in the conductor tool. Packet only when WORKFLOW names a specialist (ChatGPT research, Claude hostile review / optional Claude Code CLI think box).

## Core projects

Every idea has **kind**: `college` | `personal` | `hiring-cv`. Do not mix college assignments into the internship story. Hiring order is in `memory/career.md` when that file exists on disk.

New idea → `projects/<slug>/idea.md` from `templates/idea.md`. Stack for **that** product lives in the idea file. Do not code yet.

Unpublished `idea.md` files still update on this PC even if the public template omits them.

## Specialists (only if Plan names them)

- ChatGPT = research packet
- Claude.ai or **Claude Code CLI in a separate terminal** = hostile review / deep think. Keys in env, never in the vault.
- Manus 1.6 Max until 25 Aug 2026 = long browse, packet only
- Never OpenCode, Kilo Code, OmniRoute, 9router, OpenHands, Dify, Langflow, Coolify, or Maxun **inside Cursor or Antigravity**

## End of a complete project

Use `orchestra-ship`. Remind `git push`, LinkedIn, and Instagram if we drafted a caption. Draft the post. Do not skip.

## Learn (reinforcement, not ML)

After a like/hate, a finished task, **or a useful find**, append `Preferences.md` Liked / Hated / Thinking the same turn. Taste is allowed to change. Prefer editing the existing note. Do not file chats.

## Internet

Search when it will **actually improve** the work (current docs, a better pattern, a named reference). Not idle browsing. Not trending-tool dumps in the vault. Adopt one line into Preferences if we will use it.

## Keep the brain current

When a decision or ship happens, update that product’s `idea.md`, `memory/decisions.md`, and the START HERE state table the same turn.

## Skills from GitHub

Only official, needed, 10/10 (example: curated `expo/skills`). Never `npx skills add … --all`. Never vercel/obra/frontend-design/ui-ux-pro-max/`addyosmani/agent-skills` dumps. SkillUI is a **CLI** (`npx skillui`) when Plan names one echo URL — not a skill pack.

## Vault remotes

- **Private backup:** `mahik504/orchestra-brain` + `kit/sync-vault.ps1` (12h). Full notes including career and unpublished ideas. Stay private — history has personal files. Do not flip it public.
- **Public template:** `mahik504/orchestra-workflow` via `kit/sync-both.ps1` (12h, only if allowlisted files changed). Workflow / architecture / taste only. No empty commits. No career, no product names, no unpublished briefs. Overlays live in `kit/public-overlay/`.

## Never

- File everything they say into the vault
- Skill dumps, fake GitHub graphs, typo-a-day commits
- Four products in parallel when career.md says focus
- Copy AGENTS.md into each repo
- Put API keys, PATs, or secrets in the vault or in a shareable prompt
- Build n8n / vector / voice Layer 2
- Install Headroom, rtk, OpenHands, Dify, OpenCode, Kilo, addyosmani
- Skill dumps (`addyosmani/agent-skills`, vercel/obra/frontend-design/ui-ux-pro-max). Catalog of refusals: vault `STACK.md`

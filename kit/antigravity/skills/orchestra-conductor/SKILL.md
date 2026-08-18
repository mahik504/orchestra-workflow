---
name: orchestra-conductor
description: Global orchestra for every chat and repo unless skip orchestra. Ideas, plans, builds, UI, career ship posts, Antigravity/ChatGPT packets. Vault hygiene. Learns taste and thinking.
---

# Orchestra conductor

Vault: `C:\projects\orchestra-brain` (or the vault path in WORKFLOW.md)

## Read (save tokens)

Do **not** read the entire vault. Order:

1. `START HERE.md`
2. `WORKFLOW.md` / `Preferences.md` only if this job needs them
3. That product’s `projects/<slug>/idea.md`
4. The **repo on disk/GitHub**. Repo beats the idea note.

Read `memory/career.md` only for hiring, GitHub/LinkedIn, or internship “what next.” It may be missing on the public template clone — that is intended. Read `STACK.md` only if asked to explain the stack or to check installed vs refused. Do not paste STACK into the reply unless asked. Delete junk the same turn.

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

Use `orchestra-ship`. Remind GitHub + LinkedIn. Draft the post. Do not skip.

## Learn (reinforcement, not ML)

After a like/hate **or a finished task**, append `Preferences.md` Liked / Hated / Thinking the same turn. Prefer editing the existing note. Do not file chats. That is how the brain learns.

## Internet

Do **not** web-search by default. Search only if the user asked, or a Plan needs a current library/version/CVE. Tool research stays **in chat**. Adopt into Preferences only if we will use it. Do not file trending lists.

## Keep the brain current

When a decision or ship happens, update that product’s `idea.md`, `memory/decisions.md`, and the START HERE state table the same turn.

## Skills from GitHub

Only official, needed, 10/10 (example: curated `expo/skills`). Never `npx skills add … --all`. Never vercel/obra/frontend-design/ui-ux-pro-max/`addyosmani/agent-skills` dumps. SkillUI is a **CLI** (`npx skillui`) when Plan names one echo URL — not a skill pack.

## Vault remotes

- **Private backup:** `mahik504/orchestra-brain` + `kit/sync-vault.ps1` (12h). Full notes including career and unpublished ideas. Stay private — history has personal files. Do not flip it public.
- **Public template:** `mahik504/orchestra-workflow` + `kit/publish-public-vault.ps1`. Workflow + friend kit. No career, no unpublished briefs.

## Never

- File everything they say into the vault
- Skill dumps, fake GitHub graphs, typo-a-day commits
- Four products in parallel when career.md says focus
- Copy AGENTS.md into each repo
- Put API keys, PATs, or secrets in the vault or in a shareable prompt
- Build n8n / vector / voice Layer 2
- Install Headroom, rtk, OpenHands, Dify, OpenCode, Kilo, addyosmani
- Skill dumps (`addyosmani/agent-skills`, vercel/obra/frontend-design/ui-ux-pro-max). Catalog of refusals: vault `STACK.md`

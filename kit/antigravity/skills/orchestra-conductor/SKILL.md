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

- **Cursor present:** Cursor Plans and implements. Antigravity is a specialist (polish, tests, CI, hostile review) only when Plan names it **or** the work is showable web/Android (then a polish packet is required). Packet in, summary out.
- **Antigravity only (no Cursor):** Antigravity **is** the conductor. Same vault rules. Plan first (strongest thinking model), then Stitch, then implement. Do not wait for a Cursor packet.

## Cursor modes (the human switches these)

You cannot flip Plan / Agent / Ask / Debug / Multitask. Say it:

| Mode | Say this | When |
| --- | --- | --- |
| Plan | “Switch to **Plan** mode.” | Architecture, spec.md, before code |
| Agent | “Switch to **Agent** mode.” | Implementing after spec.md exists |
| Ask | “Switch to **Ask** mode.” | Read-only |
| Debug | “Switch to **Debug** mode.” | Runtime failure |
| Multitask | “Switch to **Multitask** only if two jobs are independent.” | Rare. Never two products |

## Models (must-do)

This parent chat **cannot** become Opus/Fable. If this chat is Grok (or another glue model) and the job is UI, architecture, or hostile review:

1. Tell them the **dropdown** to pick for the next Plan/Agent turn (Opus / Fable / GPT 5.6 thinking for Plan; Fable or Opus for UI).
2. **Spawn a Task subagent** with the slug from `WORKFLOW.md` (Fable for UI, Opus for hostile review, Gemini Flash or Kimi for mechanical). Do not implement a Stitch-match UI entirely as Grok.

Antigravity Pro: **Gemini 3.7 Flash High** packet for polish/CI on showable web/Android. ChatGPT Go or Perplexity = research packets only.

## Default

Stay in the conductor tool. Packet only when WORKFLOW names a specialist.

## Spec + Ralph-thin

After Plan: write `projects/<slug>/spec.md` from `templates/spec.md`. One user story per Agent/subagent pass. Learnings in `progress.md`. Do **not** install the Ralph CLI, Spec Kit `specify init`, or slash-command packs. Those fight this loop.

## Import then edit (by kind)

Read `kind` on idea.md. College includes hackathon. Stitch is **not** mandatory if Plan named a library / template / shader / animation — combine those against `projects/<slug>/design.md` and tint.

- **college:** one SkillUI URL **or** one 21st / shadcn / unlumen / smoothui / neobrutalism block **or** React Bits layout **or** one ThreeUI component. Tint to one visual world (`neo` = neobrutalism.com). React Bits required on hackathon web. Never generate chrome from zero.
- **hiring-cv:** Stitch + one echo URL, **or** named library/shader if Plan said skip Stitch. shadcn primitives only if Plan named them.
- **personal LYRA-class:** Stitch Screen 1. No shadcn. No 21st. 3D = one Shadertoy/drei volume port.

Always-on 21st / Aceternity / Magic MCP stay off. Do not call shadcn MCP on a LYRA-class operator HUD. Do not generate a widget or nucleus from a screenshot loop.

## ECC method (not the dump)

[affaan-m/ECC](https://github.com/affaan-m/ECC) proves load-by-job and isolate plan / implement / review. **Do not** install ECC into Cursor or Antigravity. Instincts = Preferences same turn. No homunculus hooks. Stay open: new named sources in STACK; vault may still change.

## When to load (do not dump every skill)

| Job | Load |
| --- | --- |
| Plan / spec | conductor + vault only |
| 2D lock | stitch-* + taste-design |
| UI polish after a still | impeccable + animate / review-animations |
| Android | curated expo-* |
| Security | ship-safe; Strix only on **our** app when scanning |
| Docs / papers | orchestra-docs |
| College 0-to-1 | Antigravity may use hackathon-rocket if installed. Do not add it to Cursor always-on |

Never load science plugins, ui-ux-pro-max, vibe-design-pro, ECC’s 286 skills, or the whole skill folder on a LYRA orb pass.

## Core projects

Every idea has **kind**: `college` | `personal` | `hiring-cv`. Do not mix college assignments into the internship story. Hiring order is in `memory/career.md` when that file exists on disk.

New idea → `projects/<slug>/idea.md` from `templates/idea.md`. Stack for **that** product lives in the idea file. Do not code yet.

Unpublished `idea.md` files still update on this PC even if the public template omits them.

## Specialists (only if Plan names them, except Antigravity polish on showable web)

- ChatGPT Go or Perplexity = research packet
- **OpenHuman** (desktop you install) = packet specialist. Not conductor. Do not dump its skill catalog into Cursor.
- Claude.ai or **Claude Code CLI in a separate terminal** = hostile review / deep think. Keys in env, never in the vault.
- Manus 1.6 Max until 25 Aug 2026 = long browse, packet only
- Never OpenCode, Kilo Code, OmniRoute, 9router, OpenHands, Dify, Langflow, Coolify, Maxun, Guildly, ECC in Cursor/Antigravity, always-on 21st.dev / Aceternity / Magic MCP, or the Ralph CLI **inside Cursor or Antigravity**

## End of a complete project

Use `orchestra-ship`. Remind `git push`, LinkedIn, and Instagram if we drafted a caption. Draft the post. Do not skip. Every new **public** repo gets a LICENSE the same day.

## Learn (reinforcement, not ML)

After a like/hate, a finished task, **or a useful find**, append `Preferences.md` Liked / Hated / Thinking **in this same chat**. Taste is allowed to change. Prefer editing the existing note. Do not file chats. Do not wait for a later session.

## Internet

Search when it will **actually improve** the work (current docs, a better pattern, a named reference). If spec.md has a problem statement (college or hiring), one competitor/paper pass — file `projects/<slug>/research.md` only if we will use it. Not idle browsing. Not trending-tool dumps in the vault. Adopt one line into Preferences if we will use it.

## Keep the brain current

When a decision or ship happens, update that product’s `idea.md`, `memory/decisions.md`, and the START HERE state table the same turn.

## Skills from GitHub

Only official, needed, 10/10 (example: curated `expo/skills`). Never `npx skills add … --all`. Never vercel/obra/frontend-design/ui-ux-pro-max/`addyosmani/agent-skills` as **global** dumps. College Plan may copy **one** DESIGN.md from awesome-design-md into `projects/<slug>/ref/`. SkillUI is a **CLI** (`npx skillui`) when Plan names one echo URL.

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
- Install Headroom, rtk, OpenHands, Dify, OpenCode, Kilo, addyosmani, Spec Kit CLI, Ralph CLI, ECC into Cursor/Antigravity, always-on 21st.dev / Aceternity / Magic MCP, ui-ux-pro-max as a global skill, science plugin packs, vercel-labs/agent-skills dump, screenshot-to-code
- Reminder MCP / LinkedIn scrapers (reminders stay in the user rule + chat)
- Skill dumps (`addyosmani/agent-skills`, vercel/obra/frontend-design/ui-ux-pro-max). Catalog of refusals: vault `STACK.md`

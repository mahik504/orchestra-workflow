---
title: Workflow
---

# Workflow

Global for every chat and repo until you say **skip orchestra**.

You talk. Cursor **Plan**s. You follow. Specialists only when Plan names them.

```
you fill templates/plan-brief.md (or paste those answers)
  → Cursor Plan (questions + START HERE + that idea.md, not the whole vault)
  → optional: npx skillui on ONE echo URL Plan named → projects/<slug>/ref/
  → Stitch = 2D screens / app chrome
  → 3D scenes in code (R3F or Spline), not in Stitch
  → Cursor Agent implements (React Bits only for earned web motion)
  → Antigravity only if Plan said so
  → optional Claude Code CLI (separate window) only if Plan said so
  → ship-safe on your app; Strix when Docker is ready and the app is yours
```

## Read the vault (save tokens)

Do **not** read the entire vault. Order:

1. `START HERE.md`
2. `WORKFLOW.md` / `Preferences.md` only if this job needs them
3. That product’s `projects/<slug>/idea.md`
4. The repo on disk

Do not paste `STACK.md` unless asked to explain the stack. `memory/career.md` only for hiring jobs; it lives on this PC and the private remote. Delete junk the same turn.

## Learn (same turn)

After a like/hate or a finished task, append Preferences Liked / Hated / Thinking. Prefer editing existing notes. Do not file chats. That is how the brain learns (reinforcement, not ML). Unpublished `idea.md` still updates on disk even if the public template omits that folder.

## Internet

Do **not** web-search by default. Search only if the user asked, or a Plan needs a current library/version/CVE. Tool research stays in chat unless we adopt it into Preferences.

## Keep the brain current

When a decision or ship happens, update that product’s `idea.md`, `memory/decisions.md`, and the START HERE state table the same turn.

## Who does which part

| Job | Use |
| --- | --- |
| Idea, architecture, stack, “what next?” | **Cursor Plan** |
| Code | **Cursor Agent** |
| It broke | **Cursor Debug** |
| 2D UI / screens / DESIGN.md | **Stitch** |
| Interactive 3D in a web/app screen | **Cursor Agent** + React Three Fiber (or Spline drop-in) |
| Extra polish / tests / CI | **Antigravity** + packet |
| Long cloud browse / wide research (until 2026-08-25) | **Manus 1.6 Max** + packet · not the conductor |
| Research | **ChatGPT** + packet |
| Hostile architecture / deep think | **Claude.ai packet** or **Claude Code CLI** (own window) · only if Plan names it |
| Security on **your** app | **ship-safe** always; **Strix** when we scan |
| Library docs | **Context7** |
| Click the running web UI | **Playwright** |
| Android (Play Store, Apple later) | **Expo (React Native)** — one TS codebase |

Do not switch the whole factory to Fable / Firebase Studio. 3D = R3F/Spline + a strong Agent model.

Do **not** install OpenCode, Kilo Code, OmniRoute, or similar **inside Cursor or Antigravity**. That would compete with this conductor loop.

## Models (pick in the tool you are in)

You have Cursor models, ChatGPT Go, Antigravity Pro, and **Manus 1.6 until 25 Aug 2026**. Use **job**, not brand.

| Job | Where | Model class |
| --- | --- | --- |
| Plan, architecture | Cursor Plan | Strongest **thinking** model in Cursor |
| Implement, 3D, UI match | Cursor Agent | Strongest Agent model |
| Mechanical refactors, tests | Cursor Agent or Antigravity | Fast / Flash |
| Extra UI polish, CI | Antigravity | Gemini 3.1 Pro (or best Gemini) |
| Hostile review / deep think | Claude.ai packet, or Claude Code CLI in a **separate terminal** | Opus-class via your router credits · time-boxed · Plan must name it |
| Long browse / wide research | Manus 1.6 Max (until 25 Aug) | Packet only. Do not let it rebuild the app in Manus cloud |
| Research / product | ChatGPT Go | ChatGPT’s strongest available |
| Tiny college GUI (due now) | Cursor Agent only | Skip Stitch. PySimpleGUI. Ship today |

If Cursor quota dies: packet to ChatGPT Go and/or Antigravity. Same vault.

## Surfaces

| Surface | How |
| --- | --- |
| Marketing / desktop web | React + Vite (or Next if the repo already is) |
| Android + later iOS | Expo. Play Store first. Apple is a later build of the same app. |
| 3D | Stitch does chrome. Scene = R3F (`@react-three/fiber` + `drei`) or a Spline embed. One 3D moment per screen unless Plan says more. |
| Mini web game | Canvas 2D + physics (Rapier or Matter) unless the game *is* 3D |

## Hard rules

1. Fill `templates/plan-brief.md` (or the same questions) before we design or code.
2. No frontend until Stitch chrome is locked, except a Plan-approved 3D/game scene **or a due-now college desktop GUI (PySimpleGUI)** — then skip Stitch.
3. Cursor implements Stitch. It does not invent a generic SaaS look.
4. Expo for mobile until you explicitly switch to Kotlin/Flutter.
5. Secrets never in git or the vault. `ship-safe` on every ship. Router API tokens live in **Windows env** or `~\.claude\settings.json` on this PC — never in this vault.
6. No skill packs. No extra MCP unless a product needs that host/db.
7. Vault: lasting notes only, under `projects/<slug>/`. Each idea has `kind`: college | personal | hiring-cv. Delete junk the same turn. Like/hate/thinking → Preferences.
8. When a project is showable: remind GitHub + LinkedIn and write `projects/<slug>/ship-post.md`. Never fake contribution graphs. Never typo-a-day commits.
9. Six-month internship: follow `memory/career.md` when it exists on disk. Do not run four products in parallel.
10. Tool research stays in chat unless we adopt it into Preferences. Official Expo skills are already installed (curated). Do not add skill packs. `npx skillui` is a CLI, not a pack — only when Plan names one echo URL.
11. Vault backup is **private git** (`kit/sync-vault.ps1` + private GitHub `mahik504/orchestra-brain`). Public template is a **separate** repo (`kit/publish-public-vault.ps1` + `mahik504/orchestra-workflow`). Do not flip the private remote public — history has career and unpublished ideas. Do not build n8n / Supabase / embeddings until an automation must query the vault without us.
12. Installed vs refused tools: `STACK.md`. Do not install OpenHands, Dify, Langflow, Coolify, Maxun, Headroom, rtk, or skill dumps.

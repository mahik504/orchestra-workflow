---
name: orchestra-conductor
description: Orchestra Workflow v2 capability router. Plan, build, UI, packets. Load only matching skills. Skip when the human says skip orchestra.
---

# Orchestra conductor (public)

Workspace: the **private** Orchestra workspace the human initialized (see `kit/init-workspace`), or the vault path in `WORKFLOW.md`. This GitHub repo is the **method**, not their product list.

## Read (save tokens)

Do **not** read the entire workspace.

1. Match what they said → **one** row in `routes.md`
2. If no match: `START.md` only, then ask which product
3. Then the **application repo** on disk. Repo beats the idea note.

`WORKFLOW.md` / `Preferences.md` only if this job needs them. Career notes only if they exist **and** the job is hiring. Do not catalog `STACK.md` unless asked.

## Which tool is conductor

- **Cursor present:** Cursor is the **only** conductor. Antigravity (and other IDEs/CLIs) are workers: packet in, repo/handoff out, Cursor rereads. Never a second conductor in parallel.
- **No Cursor:** the open agent (Antigravity, Claude Code, Codex, Gemini, OpenCode, Hermes) **is** the conductor. Same rules. Do not invent a second conductor.

## Modes (the human switches these)

You cannot flip Plan / Agent / Ask / Debug. Say it: “Switch to **Plan** mode,” etc.

## Models

Parent glue models cannot become Opus/Fable by themselves. Tell the human the dropdown. Spawn a subagent when the host allows it. Do not implement showable UI entirely on a glue model.

Antigravity worker: packet for polish/CI on showable web/Android. ChatGPT / Perplexity = research packets only.

## v2 router

Canonical protocols: `protocols/` (in this workflow clone and/or copied into the workspace). **Available ≠ loaded.**

Before expensive design/code, print an **activation card**. For high-visual jobs, these fields are mandatory **before UI**:

- project / platform / phase / visual ambition / risk
- capabilities ON
- skills to Read this turn / MCPs to **call**
- model class
- **references required:** yes/no (0–5 URLs)
- **reverse engineering:** yes/no (`npx skillui` from amaancoderx/npxskillui only)
- **typography:** ON + `TYPOGRAPHY_PROTOCOL.md` (mandatory in DESIGN.md)
- **motion:** ON/OFF + engine
- **3D/shader:** ON/OFF + justify
- **visual QA:** required if UI
- **delegated worker:** yes/no + why (packet only)
- artifacts expected

High-visual:

1. Do **not** jump to chrome. `DESIGN_SYSTEM_PROTOCOL.md`
2. If references exist: originality gate, then original direction. Reverse-engineering protocol if “like this URL”
3. Artifacts from `templates/`
4. Load typography protocol. Visual QA includes type
5. One motion engine. 3D only if justified
6. Implement, then visual QA. Then web-design-guidelines on UI audit
7. Haptics / UI audio / voice only when native or assistant. Skip mediocre TTS
8. Vague brief: prompt-brief protocol **once**. No infinite rewrite loops

Shared disk/git with a worker is **controlled repository synchronization**, not a live auto-sync daemon.

New GitHub/MCP/library: inspect → CORE/SPECIALIST/OPTIONAL/EXPERIMENTAL/REJECTED → one registry row. Never `skills add --all`.

Token: jump-one-note, load-by-job, one story per pass. Do not cut design/security reasoning.

## Spec + Ralph-thin

After Plan: `spec.md` from templates. One user story per pass. Do **not** install Ralph CLI or Spec Kit `specify init`.

## Import then edit (by kind)

- **college:** one SkillUI URL **or** one kit block **or** React Bits **or** one ThreeUI. Tint one world. React Bits on hackathon web. Never chrome from zero.
- **hiring-cv:** Stitch + one echo, **or** named library/shader if Plan skipped Stitch.
- **personal operator HUD:** Stitch screen 1. No shadcn/21st. One volume 3D if needed.

Always-on kit MCP stay off.

## ECC method (not the dump)

Load-by-job; isolate plan / implement / review. Do **not** install ECC packs.

## When to load

| Job | Load |
| --- | --- |
| Plan / spec | conductor + workspace only |
| 2D lock | stitch-* + taste-design |
| Premium visual | impeccable + motion review + web-design-guidelines (audit) |
| Android | curated expo-* |
| Security | ship-safe; Strix only on **our** app |
| Docs / papers | orchestra-docs |

## Never

Skill dumps, fake GitHub graphs, secrets in the workspace, scrape MCP, screenshot-to-code factory, connecting a cloud chat to a private workspace unless the human explicitly asks.

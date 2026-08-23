# Getting started

**clone → initialize → configure → activate → route → load → execute → verify**

This repository is the **Orchestra Workflow** (public). Your notes and products belong in a **separate private workspace**.

## 1. Clone

```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow
```

## 2. Initialize a private workspace

Windows:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/init-workspace.ps1
```

macOS / Linux:

```bash
chmod +x kit/init-workspace.sh
./kit/init-workspace.sh
```

Default location: `../orchestra-workspace` (sibling of this clone). Override:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/init-workspace.ps1 -Target "D:\work\my-orchestra"
```

```bash
ORCHESTRA_HOME=/path/to/my-orchestra ./kit/init-workspace.sh
```

The workspace is empty on purpose: `projects/`, `memory/`, `Preferences.md`, `routes.md`. Fill them with **your** work.

## 3. Configure

1. Open **both** folders in your agent: this workflow repo **or** just the copied `protocols/` + skills, **and** your private workspace, **and** the application repo you are building.
2. Put API keys in the **agent’s** env / MCP UI — never in the workspace git.
3. Edit workspace `Preferences.md` (taste) and `routes.md` (your slugs).
4. Optional: copy `protocols/` and `registries/` into the workspace if you want a single root. The init script can symlink or copy; default is copy of protocols/templates into the workspace so the private brain stays self-contained.

## 4. Activate

Install skills into local agent directories:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/install-skills.ps1
```

```bash
./kit/install-skills.sh
```

This copies `skills/*` into common locations **if they exist** on your machine:

- `~/.cursor/skills/`
- `~/.claude/skills/`
- `~/.agents/skills/`
- `~/.gemini/config/skills/`

It does **not** run `npx skills add --all`. It does **not** invent logins.

For Codex / OpenCode / Hermes: keep `AGENTS.md` in the app repo or workspace root.

## 5. Route the task

In chat: say what you are building. The conductor should:

1. Print an **activation card** (job, risk, visual ambition, capabilities ON).
2. Open **one** `routes.md` target, then the **app repo**.
3. Load only matching protocols/skills.

High-visual UI: no chrome until the card includes references, reverse engineering, typography, motion, 3D/shader, visual QA, and whether a worker IDE is packeted.

## 6. Load relevant capabilities

Read the protocol files for **this** job. Examples:

- Showable web: `protocols/DESIGN_SYSTEM_PROTOCOL.md` + `TYPOGRAPHY_PROTOCOL.md` + `VISUAL_QA_PROTOCOL.md`
- “Make it like this URL”: `REVERSE_ENGINEERING_PROTOCOL.md` then originality gate
- Auth/payments: `SECURITY_PROTOCOL.md` (ship-safe + Strix on **your** app)
- Vague ask: `PROMPT_BRIEF_PROTOCOL.md` once

## 7. Execute

Implement in the **application repository**. One spec story per pass. Packet Antigravity (or any worker) only when the conductor says so. After a worker returns, the conductor **re-reads the git diff**.

## 8. Verify

UI: visual QA protocol on a running app (Playwright or the agent’s browser).  
Security: ship-safe always; Strix only on your code.  
Never claim “looks good” from a single chat screenshot.

## 9. Persist (small)

Lasting taste → workspace `Preferences.md`.  
Decisions → `memory/decisions.md`.  
Product facts → `projects/<your-slug>/idea.md`.  
Do not file chat logs.

## Skip Orchestra

Say **skip orchestra** in a chat when you want a raw agent with no router.

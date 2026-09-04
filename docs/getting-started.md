# Getting started

**clone → private workspace → doctor → classify → plan → Design Lab (if visual) → implement → verify**

This repository is the **Orchestra Workflow** (public method). Your notes and products belong in a **separate private Brain**. Set `ORCHESTRA_HOME` to that workspace.

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

Default location: `../orchestra-workspace` (sibling of this clone), or `$ORCHESTRA_HOME` if that env var is already set. Override:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/init-workspace.ps1 -Target "D:\work\my-orchestra"
```

```bash
ORCHESTRA_HOME=/path/to/my-orchestra ./kit/init-workspace.sh
```

The workspace is empty on purpose: `projects/`, `memory/`, `Preferences.md`, `routes.md`. Fill them with **your** work.

## 3. Configure

1. Open **both** folders in your agent: this workflow repo, your private workspace, and the application repo you are building.
2. Put API keys in the **agent’s** env / MCP UI — never in the workspace git.
3. Edit workspace `Preferences.md` (taste) and `routes.md` (your slugs).
4. Set `ORCHESTRA_HOME` to the private workspace. `orchestra doctor` prints the Memory and Overlay paths it bound.

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

1. Re-brief (archetype, quality bar, platform, hard constraints).
2. Classify into one capability. Load that route. Leave other routes closed.
3. Open the **app repo** before trusting a brief.

PREMIUM / EXPERIMENTAL visual work: **Design Lab is a write lock**. Do not write frontend files until a stack card is approved. Say `skip the lab` to bypass one task.

## 6. Load relevant capabilities

Read the protocol files for **this** job. Examples:

- Showable web: `protocols/DESIGN_SYSTEM_PROTOCOL.md` + `TYPOGRAPHY_PROTOCOL.md` + `VISUAL_QA_PROTOCOL.md`
- “Make it like this URL”: `REVERSE_ENGINEERING_PROTOCOL.md` then originality gate
- Auth/payments: `SECURITY_PROTOCOL.md` (ship-safe + Strix on **your** app)

## 7. Execute

Implement in the **application repository**. One spec story per pass. Packet a worker only when the conductor says so. After a worker returns, the conductor **re-reads the git diff**.

## 8. Verify

UI: visual QA protocol on a running app (Playwright or the agent’s browser). Zero screenshots is not a Playwright pass.
Security: ship-safe always; Strix only on your code.
Never claim “looks good” from a single chat screenshot. `DONE` / `VERIFIED` needs evidence in the same message.

## 9. Persist (small)

Lasting taste → workspace `Preferences.md`.
Decisions → `memory/decisions.md`.
Product facts → `projects/<your-slug>/idea.md`.
Do not file chat logs.

## Skip Orchestra

Say **skip orchestra** in a chat when you want a raw agent with no router.

# Getting started

The **chat is Orchestra**. After bootstrap, open the app repo plus your private workspace and talk. You do not run `orchestra classify` at the start of every conversation.

This repository is the **Orchestra Workflow** (public method). Your notes and products belong in a **separate private Brain**. Named skills, MCP templates, and the public resource catalog ship with the clone. Your taste, keys, and overlay memory do not.

```
clone → bootstrap (pick hosts) → paste your keys → restart → chat
optional: persist ORCHESTRA_HOME, install orchestra.exe, doctor / classify / plan
```

Reference URLs (awwwards, shadcn/ui, GSAP, R3F, Drei, React Bits, …) live in [`registries/resources.json`](../registries/resources.json). Host extras (Gmail, Stripe, …) are a **checklist**, not a git install. See [`registries/host-stack.json`](../registries/host-stack.json).

## 1. Clone

```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow
```

## 2. Bootstrap (pick hosts)

Windows:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/bootstrap.ps1
```

Non-interactive:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/bootstrap.ps1 -Hosts cursor,antigravity -Target "D:\work\my-orchestra"
```

macOS / Linux:

```bash
chmod +x kit/bootstrap.sh kit/init-workspace.sh kit/install-local-engine.sh
./kit/bootstrap.sh
# or:
./kit/bootstrap.sh --hosts cursor,antigravity --target /path/to/my-orchestra
```

Bootstrap will:

1. Create an empty private workspace if it is missing (`projects/`, `memory/`, `Preferences.md`, `routes.md`).
2. Copy the 30 canonical skills onto the hosts you picked.
3. Copy `AGENTS.md` and the matching adapter (`.cursorrules`, `CLAUDE.md`, Antigravity MASTER-PROMPT).
4. Write MCP **templates** with `REPLACE_WITH_*` placeholders. It will not overwrite a live `mcp_config.json`.
5. Print the marketplace plugin checklist. You click Connect. Orchestra cannot log into Google, Stripe, or Stitch for you.

It does **not** run `npx skills add --all`. It does not copy anyone's Brain.

## 3. Configure the window

1. Open **the app repo** you are building. Add the private workspace as a second folder.
2. Paste API keys in the **agent's** MCP UI — never in git.
3. Edit workspace `Preferences.md` (your taste) and `routes.md` (your slugs).

Restart the IDE. The next chat in that workspace reads `AGENTS.md`. Antigravity: still paste `kit/antigravity/MASTER-PROMPT.md` once per **new** AG chat (VAULT + APP ROOT). See [`kit/antigravity/ALWAYS-ON.md`](../kit/antigravity/ALWAYS-ON.md).

## 4. Optional: persist the engine

Only if you want `orchestra` on PATH (Design Lab write lock, `classify` / `plan` / `doctor`):

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/install-local-engine.ps1 -HomeDir "D:\work\my-orchestra"
```

```bash
./kit/install-local-engine.sh /path/to/my-orchestra
```

Then:

```bash
orchestra doctor
orchestra classify --task "optional: score this brief"
orchestra plan --task "optional: dry-run the pipeline"
```

A fresh clone without `ORCHESTRA_HOME` writes only under that clone's `.orchestra/` directory.

## 5. Route the task (in chat)

Say what you are building. The conductor should:

1. Re-brief (archetype, quality bar, platform, hard constraints).
2. Classify into one capability. Load that route. Leave other routes closed.
3. Open the **app repo** before trusting a brief.

PREMIUM / EXPERIMENTAL visual work: **Design Lab is a write lock**. Do not write frontend files until a stack card is approved. Say `skip the lab` to bypass one task.

Installing Stitch, Fiber, or Stripe does not load them on every prompt.

## 6. Load relevant capabilities

Read the protocol files for **this** job. Examples:

- Showable web: `protocols/DESIGN_SYSTEM_PROTOCOL.md` + `TYPOGRAPHY_PROTOCOL.md` + `VISUAL_QA_PROTOCOL.md`
- "Make it like this URL": `REVERSE_ENGINEERING_PROTOCOL.md` then originality gate
- Auth/payments: `SECURITY_PROTOCOL.md` (ship-safe + Strix on **your** app)

## 7. Execute

Implement in the **application repository**. One spec story per pass. Packet a worker only when the conductor says so. After a worker returns, the conductor **re-reads the git diff**.

## 8. Verify

UI: visual QA protocol on a running app (Playwright or the agent's browser). Zero screenshots is not a Playwright pass.
Security: ship-safe always; Strix only on your code.
Never claim "looks good" from a single chat screenshot. `DONE` / `VERIFIED` needs evidence in the same message.

## 9. Persist (small)

Lasting taste → workspace `Preferences.md`.
Decisions → `memory/decisions.md`.
Product facts → `projects/<your-slug>/idea.md`.
Do not file chat logs.

## Skip Orchestra

Say **skip orchestra** in a chat when you want a raw agent with no router.

---
title: Connect
---

# Connect

How Antigravity, ChatGPT, and optional specialists attach to this vault.

Cursor already sees `C:\projects\orchestra-brain`. Antigravity does not, until you add the folder and paste a prompt.

**Full install + friend (no Cursor):** zip and send only `kit/antigravity/`, or point them at the public template [mahik504/orchestra-workflow](https://github.com/mahik504/orchestra-workflow). Paste `kit/antigravity/MASTER-PROMPT.md`. Do **not** send the private vault, `mcp_config.json`, or `memory/career.md`.

Or run locally: `powershell -ExecutionPolicy Bypass -File C:\projects\orchestra-brain\kit\antigravity\install.ps1`

Then **restart Antigravity**.

## Every new project (you — Cursor is still conductor)

1. Antigravity → Open the **project folder** (ODYSS, AstroVerse, …).
2. **Add Folder to Workspace** → `C:\projects\orchestra-brain`
3. New chat → paste the short box below (daily) **or** the master prompt once if skills/MCP are missing.
4. Wait until it says vault readable **yes**. After that, only paste Cursor **packets**.

### Antigravity daily bootstrap (specialist)

```
MODE: specialist
VAULT: C:\projects\orchestra-brain
APP ROOT: (the folder already open)

You are Antigravity Pro in our orchestra. Cursor is the conductor. You never see Cursor chats.

1) Confirm you can read C:\projects\orchestra-brain\START HERE.md. If you cannot, tell me to Add Folder that vault and stop.
2) Jump routes.md to one file. Do not read the whole vault. START HERE / WORKFLOW / Preferences only if this job needs them.
3) You are a specialist: extra UI polish, tests, CI, or a second implementation pass. Not the planner unless a packet says so.
4) Follow Preferences.md. No generic AI SaaS UI. No skill packs. No fake GitHub graphs.
5) 2D UI matches locked Stitch / projects/<slug>/design.md. 3D = R3F or Spline only if the packet says so. Android = Expo.
6) Secrets never in the vault. Lasting notes only under projects/<slug>/. Like/hate → Preferences the same turn.
7) Search when it will actually improve this packet (current docs, a better pattern). Not idle browsing. Not trending-tool dumps.
8) When you finish, write a short markdown summary I can paste back into Cursor.

Reply: Orchestra connected. Vault readable: yes/no. App root: <path>. Ready for a packet.
Do not start a product. Wait for a packet.
```

If it says vault **no**: Add Folder again, paste the same box again.

Later jobs: paste only the Paste block from `projects/<slug>/packet.md`. Then in Cursor: “Packet is back.”

---

## ChatGPT Go (research only)

New Project → upload `WORKFLOW.md` + `Preferences.md`. Custom instructions:

```
You are a specialist, not the conductor. Only use the pasted packet. Do not invent our stack. Return markdown for C:\projects\orchestra-brain (path in the packet). Do not web-search unless the packet says the fact must be current.
```

Each job: paste the packet Paste block. Do not implement the repo in ChatGPT.

---

## Manus 1.6 (until 25 Aug 2026 only)

Cloud computer + browser. **Not** the conductor. **Not** Stitch. **Not** Expo. After the 25th, stop unless we renew.

Use **1.6 Max** for one bounded job: wide web research, a long read of public docs, or a second opinion on a README/eval. Then paste the summary back into Cursor.

Do **not**: let Manus scaffold a parallel app, use Design View instead of Stitch, use Manus Mobile Dev instead of Expo, or connect Drive/Notion/Instagram. GitHub connector only if we explicitly packet a **public** repo URL — never paste PATs.

New Manus task → paste:

```
You are Manus 1.6 Max in our orchestra. Cursor is the conductor. You never see Cursor chats.

Use only this packet. Do not invent our stack. Do not create a new product repo in Manus cloud unless the packet says to return files as markdown we can paste.

RULES:
- Follow the pasted Preferences. No generic AI SaaS UI. No skill dumps. No fake GitHub graphs.
- 2D UI is Stitch’s job. 3D is R3F/Spline in our repo. Android is Expo.
- Never invent metrics, users, or accuracy.
- Secrets never in your output.
- When done, return markdown: what you found, file paths or URLs, what Cursor should do next.

JOB:
{{paste the packet Paste block}}
```

Then in Cursor: “Manus packet is back.”

---

## Claude Code CLI (optional think box)

The “Claude Code from GitHub + custom router” thing is **Anthropic’s official Claude Code CLI** pointed at a gateway with env vars. You do **not** need OmniRoute, 9router, OpenCode, Kilo, OpenHands, Dify, Langflow, Coolify, or Maxun **inside Cursor or Antigravity** — those would become a second conductor and fight this vault. Catalog: `STACK.md`.

Use the CLI only when **Plan names it** (hostile architecture, a long refine of a packet). Own terminal window. Time-box it. Paste a summary back into Cursor, same as ChatGPT.

Install (once, on this PC, not in the vault):

```text
npm install -g @anthropic-ai/claude-code
```

Then set **Windows User environment variables** (Settings → Environment variables). Never store the token in this vault, in git, or in a shareable prompt.

PowerShell for **one session** (placeholders only — paste the host and token from the provider console on this PC, never into the vault):

```powershell
# Confirm the exact host in your router console. No /v1 suffix — the CLI appends /v1/messages.
$env:ANTHROPIC_BASE_URL = "https://YOUR-ROUTER-HOST"
$env:ANTHROPIC_AUTH_TOKEN = "YOUR-TOKEN"
Remove-Item Env:ANTHROPIC_API_KEY -ErrorAction SilentlyContinue
claude
```

Local-only settings file (still **not** in the vault): `%USERPROFILE%\.claude\settings.json`

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://YOUR-ROUTER-HOST",
    "ANTHROPIC_AUTH_TOKEN": "YOUR-TOKEN"
  }
}
```

Router credits are finite. If the Plan did not name Claude, skip it and stay in Cursor.

---

## Vault remotes

| Remote | Visibility | Script |
| --- | --- | --- |
| [mahik504/orchestra-brain](https://github.com/mahik504/orchestra-brain) | **private** backup | `kit/sync-both.ps1` (12h task `OrchestraBrainVaultSync`) |
| [mahik504/orchestra-workflow](https://github.com/mahik504/orchestra-workflow) | **public** template | `kit/publish-public-vault.ps1` |

The Obsidian Git plugin only runs while Obsidian is open. The Windows task runs `kit/sync-both.ps1`: private vault backup, then public template **if allowlisted files changed**. Push needs git credentials on this PC.

Do **not** flip `orchestra-brain` to public. History contains career and unpublished ideas. The public snapshot is `orchestra-workflow`.

Never commit `mcp_config.json`, `.env`, API keys, or PATs. Never force-push. Never skip hooks.

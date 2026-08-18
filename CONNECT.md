---
title: Connect
---

# Connect

How Antigravity, ChatGPT, and optional specialists attach to a vault that uses this workflow.

Cursor already sees the vault folder if it is in the workspace. Antigravity does not, until you add the folder and paste a prompt.

**Friend (no Cursor):** zip and send only `kit/antigravity/`, or clone [mahik504/orchestra-workflow](https://github.com/mahik504/orchestra-workflow). Paste `kit/antigravity/MASTER-PROMPT.md`. Do **not** send a private vault, `mcp_config.json`, or career notes.

## Every new project (Cursor is still conductor)

1. Antigravity → Open the **app folder** (the repo you are building).
2. **Add Folder to Workspace** → the orchestra vault path.
3. New chat → paste the short box below **or** the master prompt once if skills/MCP are missing.
4. Wait until it says vault readable **yes**. After that, only paste **packets**.

### Antigravity daily bootstrap (specialist)

```
MODE: specialist
VAULT: (path to the orchestra vault)
APP ROOT: (the folder already open)

You are Antigravity Pro in our orchestra. Cursor is the conductor. You never see Cursor chats.

1) Confirm you can read START HERE.md in the vault. If you cannot, tell me to Add Folder that vault and stop.
2) Jump routes.md to one file. Do not read the whole vault.
3) You are a specialist: extra UI polish, tests, CI, or a second implementation pass. Prefer Gemini 3.7 Flash High.
4) Follow Preferences.md. No generic AI SaaS UI. No skill packs. No fake GitHub graphs.
5) 2D UI matches locked Stitch / projects/<slug>/design.md. 3D = R3F or Spline only if the packet says so. Android = Expo.
6) Secrets never in the vault. Lasting notes only under projects/<slug>/.
7) Search when it will actually improve this packet. Not idle browsing.
8) When you finish, write a short markdown summary I can paste back into Cursor.

Reply: Orchestra connected. Vault readable: yes/no. App root: <path>. Ready for a packet.
Do not start a product. Wait for a packet.
```

## ChatGPT Go (research only)

Upload `WORKFLOW.md` + `Preferences.md`. Custom instructions: specialist only; return markdown; search only if the packet says the fact must be current.

## Perplexity (research only)

Same job as ChatGPT Go. Packet in, markdown out. Do not implement the repo. No Perplexity MCP.

## Manus 1.6 (until 25 Aug 2026 only)

Packet specialist for long browse. Not the conductor. Not Stitch. Not Expo.

## Claude Code CLI (optional think box)

Own terminal. Plan must name it. Keys in **Windows env**, never in the vault. Never OpenCode / Kilo inside Cursor or Antigravity.

## Vault remotes

| Remote | Visibility | Script |
| --- | --- | --- |
| Author’s vault | **private** backup | `kit/sync-both.ps1` (12h task `OrchestraBrainVaultSync`) |
| [mahik504/orchestra-workflow](https://github.com/mahik504/orchestra-workflow) | **public** template | `kit/publish-public-vault.ps1` |

The Windows task: private backup first, then this public template **if allowlisted files changed**. No empty commits.

Never commit `mcp_config.json`, `.env`, API keys, or PATs. Never force-push. Never skip hooks.

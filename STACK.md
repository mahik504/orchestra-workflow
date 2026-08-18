---
title: Stack
---

# Stack

Showable catalog of **this** orchestra. Product facts stay in `projects/<slug>/idea.md`. This file maps the loop, what is installed, what GitHub we actually have, and what we refused.

Vault on this PC: `C:\projects\orchestra-brain`

| Remote | Visibility | Role |
| --- | --- | --- |
| [mahik504/orchestra-brain](https://github.com/mahik504/orchestra-brain) | **private** | Full backup (career + unpublished ideas). 12h sync. Do not make public — history has personal files. |
| [mahik504/orchestra-workflow](https://github.com/mahik504/orchestra-workflow) | **public** | Workflow template for friends. No career.md, no unpublished briefs. |

Hiring pins are product repos, not this catalog: [AstroVerse](https://github.com/mahik504/AstroVerse), [AirLens](https://github.com/mahik504/AirLens).

## How the loop works

```
you fill templates/plan-brief.md
  → Cursor Plan (START HERE + that idea.md, not the whole vault)
  → optional: npx skillui on ONE echo URL Plan named → projects/<slug>/ref/
  → Stitch = 2D screens / app chrome
  → 3D in code (R3F or Spline), not Stitch
  → Cursor Agent implements (React Bits only for earned web motion)
  → specialists only if Plan named them
  → ship-safe on our app; Strix when Docker is ready and the app is ours
  → when showable: GitHub + LinkedIn draft in projects/<slug>/ship-post.md
```

| Step | Who | Notes |
| --- | --- | --- |
| Idea / architecture / “what next?” | **Cursor Plan** | Strongest thinking model. Brief required. |
| 2D screens / DESIGN.md | **Stitch** | Chrome only. One visual world per product. |
| Code | **Cursor Agent** | Match Stitch. Repo on disk beats vault. |
| It broke | **Cursor Debug** | |
| Extra polish / tests / CI | **Antigravity** | Packet in, summary out. Not the planner if Cursor is present. |
| Library docs | **Context7** | |
| Click our running web UI | **Playwright** | Not a web-scrape product. |
| Security on **our** app | **ship-safe** always; **Strix** when we scan | |
| Research | **ChatGPT Go** | Packet only. |
| Long browse (until **2026-08-25**) | **Manus 1.6 Max** | Packet only. Not conductor, not Stitch, not Expo. |
| Hostile review / deep think | **Claude.ai packet** or **Claude Code CLI** | CLI = separate terminal. Router host/token via **Windows env**. Never inside the IDE. |
| Android | **Expo** | Play Store first. |

If there is **no Cursor**, Antigravity is the conductor. Same loop. Friend zip is only `kit/antigravity/` — never the private vault, never `mcp_config.json`.

## Skills installed (inspected 2026-08-18)

Same allowlist in Cursor `%USERPROFILE%\.cursor\skills` and Antigravity `%USERPROFILE%\.gemini\config\skills`. Kit copies only `orchestra-*` + `ship-safe`.

| Skill | Job | Source |
| --- | --- | --- |
| `orchestra-conductor` | Global loop unless “skip orchestra” | this vault / kit |
| `orchestra-vault` | Lasting notes only; delete junk | this vault / kit |
| `orchestra-ship` | GitHub + LinkedIn when showable | this vault / kit |
| `orchestra-docs` | Slides / reports / papers | this vault / kit |
| `ship-safe` | Defensive security on **our** apps | this vault / kit |
| `expo-router` | Expo file routing | `expo/skills` (curated) |
| `expo-project-structure` | New Expo layout | `expo/skills` |
| `expo-data-fetching` | Network / loaders | `expo/skills` |
| `expo-native-ui` | Native controls after Stitch is locked | `expo/skills` |
| `expo-upgrade` | SDK upgrades | `expo/skills` |
| `expo-dev-client` | Dev clients | `expo/skills` |
| `eas-app-stores` | Play Store / later Apple | `expo/skills` |
| `stitch-generate-design` | New Stitch screens | `google-labs-code/stitch-skills` |
| `stitch-manage-design-system` | Stitch tokens | stitch-skills |
| `stitch-extract-design-md` | DESIGN.md from our code | stitch-skills |
| `stitch-extract-static-html` | Static HTML snapshot | stitch-skills |
| `stitch-code-to-design` | Code → Stitch | stitch-skills |
| `stitch-upload-to-stitch` | Upload assets | stitch-skills |
| `stitch-react-components` | Stitch → React | stitch-skills |
| `taste-design` | Anti-generic DESIGN.md | stitch-skills |
| `impeccable` | UI polish against a locked direction | `pbakaus/impeccable` |
| `animate` | Motion from scratch | `emilkowalski/skills` |
| `review-animations` | Motion audit | `emilkowalski/skills` |
| `emil-design-eng` | UI polish philosophy | `emilkowalski/skills` |
| `penetration-testing-with-strix` | Scan **our** app | `usestrix/strix` |
| `fix-security-vulnerabilities-with-strix` | Patch after Strix | `usestrix/strix` |
| `ci-security-scanning-with-strix` | CI Strix | `usestrix/strix` |

`%USERPROFILE%\.agents\skills` also has `find-animation-opportunities` and `improve-animations`. Those are **not** on the Antigravity allowlist. Do not mirror the whole `~/.agents` dump.

On-demand CLI, not a skill pack: `npx skillui` when Plan names **one** echo URL.

## MCP (names only — no keys)

**Cursor** user `mcp.json`: `orchestra-brain`, `stitch`.

**Cursor marketplace plugins** (also expose MCP): GitHub, Playwright, Context7. Extra Cursor plugin: **caveman** (terse talk, not a second conductor).

**Antigravity** `mcp_config.json` servers: `orchestra-brain`, `StitchMCP`, `github-mcp-server`, `playwright`, `context7`.

Do not add Expo MCP, Headroom MCP, code-review-graph MCP, or extra hosts unless asked. Example without secrets: `kit/antigravity/mcp_config.example.json`.

## Product kinds (do not mix)

| kind | Meaning | Hiring? |
| --- | --- | --- |
| `hiring-cv` | Tough, public, internship story | yes, in career order |
| `personal` | Studio / parked | no until unparked |
| `college` | Assignments, due-now, basic | **never** the internship flagship |

Hiring order: **AstroVerse**, then **ODYSS or portfolio**. YUMIT parked. AirLens = college.

**Public product repos (what a recruiter can open):**

| Product | Kind | Disk | GitHub |
| --- | --- | --- | --- |
| AstroVerse | hiring-cv | `C:\projects\AstroVerse` | public `mahik504/AstroVerse` |
| AirLens | college | `C:\projects\DVP` | public `mahik504/AirLens` |

Unpublished briefs stay on this PC + the private remote only (not listed as secrets; not exported).

Per-product stack lives in that `idea.md`, not here.

## GitHub we actually have ([mahik504](https://github.com/mahik504))

**Public (what a recruiter or friend can open):**

| Repo | Role |
| --- | --- |
| `mahik504/AstroVerse` | Hiring flagship until something else ships |
| `mahik504/AirLens` | College GUI. Real course work. Not the internship story |
| `mahik504/orchestra-workflow` | Public orchestra **template** (this loop, not a hiring pin) |

**Private (not the hiring pin list):** course/sandbox repos, plus **`orchestra-brain`** (full vault backup).

Missing for hiring: public profile README repo `mahik504/mahik504`. Do not pin tutorial clones.

## Specialists (time-boxed)

| Tool | Until / how | Not |
| --- | --- | --- |
| Manus 1.6 Max | **25 Aug 2026**. Packet for long browse / wide research | Conductor, Stitch, Expo, extra connectors |
| Claude Code CLI | Optional think box. Own terminal. Plan must name it. Env vars on this PC (router host + token placeholders only in notes) | Plugin inside Cursor or Antigravity |
| ChatGPT Go | Research packet | Implementing the repo in ChatGPT |

## Vault backup

| Piece | Status (2026-08-18) |
| --- | --- |
| Local git | `C:\projects\orchestra-brain` on `master` |
| 12h task | `OrchestraBrainVaultSync` → `kit/sync-vault.ps1` → **private** origin |
| Private remote | **yes** · `https://github.com/mahik504/orchestra-brain` · **stays private** |
| Public template | `kit/publish-public-vault.ps1` → `mahik504/orchestra-workflow` |
| Secrets | `.gitignore` blocks `mcp_config.json`, `.env`, tokens, `memory/local-notes.md`. Never commit keys |
| Layer 2 (n8n / vectors / voice) | **not built** |

The Obsidian Git plugin only runs while Obsidian is open. The Windows task is the reliable commit of the **private** vault. Push needs git credentials on this PC.

Local-only gitignored files (`memory/local-notes.md`) are **not** in either remote. Career.md **is** in the private remote and **not** in the public template.

## What is good / what is missing

**Good:** Layer-1 vault; kinds labeled; Stitch + curated Expo + ship-safe + Strix; Playwright for **our** UI; AirLens shipped as college; private vault remote + public template; friend kit is a zip of `kit/antigravity/` only.

**Missing:** Honest AstroVerse hiring pass; profile README; ODYSS or portfolio as a clickable product; `gh` CLI login if the 12h task should push without Credential Manager.

## Refused on purpose (do not install)

Second conductors / other products: OpenHands, Dify, Langflow, Coolify, Maxun, OpenCode, Kilo, OmniRoute, 9router — not inside Cursor or Antigravity.

Skill dumps: `addyosmani/agent-skills`, vercel packs, obra/superpowers, frontend-design, ui-ux-pro-max.

Headroom MCP and rtk still deferred (quota is not failing). Token control is vault hygiene: short notes, read START HERE first, router, don’t dump chats. Optional later if quota hurts.

Extra servers we do not need: Stirling-PDF, Crawl4AI, browser-use (Playwright already clicks **our** UI), code-review-graph MCP, Open-Generative-AI studio.

Garry Tan **gstack** = another software factory for Claude Code. Skip.

## Evaluated GitHub list (2026-08-18)

URL is enough. Do not clone these into `C:\projects` as new apps.

| Repo | Verdict | Why |
| --- | --- | --- |
| [codecrafters-io/build-your-own-x](https://github.com/codecrafters-io/build-your-own-x) | **bookmark** (high) | Learning curriculum. Practice track in career.md (private). Not a Cursor skill. Do not clone the monorepo. |
| [petergyang/no-ai-slop](https://github.com/petergyang/no-ai-slop) | **adopt-thin** | Folded 4 Taste bullets. Not a skill pack. |
| [DietrichGebert/ponytail](https://github.com/DietrichGebert/ponytail) | **adopt-thin** | YAGNI ladder folded into How to work. Do not install the plugin / OpenCode adapter. |
| [MadsLorentzen/ai-job-search](https://github.com/MadsLorentzen/ai-job-search) | **bookmark** (later) | Useful apply *method* (fit → tailor CV). Also has portal scrapers. Not a global skill. Internship still = shipped repos first. |
| [headroomlabs-ai/headroom](https://github.com/headroomlabs-ai/headroom) | **deferred** | Still skipped. Vault hygiene is the token control we have. Optional later if quota hurts. |
| [rtk-ai/rtk](https://github.com/rtk-ai/rtk) | **deferred** | Same as Headroom. Not installed. |
| [garrytan/gstack](https://github.com/garrytan/gstack) | **skip** | 23-role Claude Code factory. Second conductor. |
| [langflow-ai/langflow](https://github.com/langflow-ai/langflow) | **skip** | Other product (agent workflow builder). |
| [langgenius/dify](https://github.com/langgenius/dify) | **skip** | Other product. No docker-compose. |
| [OpenHands/OpenHands](https://github.com/OpenHands/OpenHands) | **skip** | Other coding agent. Would fight this loop. |
| [getmaxun/maxun](https://github.com/getmaxun/maxun) | **skip** | No-code scrape platform. |
| [coollabsio/coolify](https://github.com/coollabsio/coolify) | **skip** | Self-host PaaS. Not our deploy story. |
| [browser-use/browser-use](https://github.com/browser-use/browser-use) | **skip** | Agent web automation. We have Playwright for **our** UI. |
| [Stirling-Tools/Stirling-PDF](https://github.com/Stirling-Tools/Stirling-PDF) | **skip** | PDF server. Orchestra does not need it. |
| [unclecode/crawl4ai](https://github.com/unclecode/crawl4ai) | **skip** | Crawl/scrape server. |
| [tirth8205/code-review-graph](https://github.com/tirth8205/code-review-graph) | **skip** | Extra MCP that writes into every AI tool. Cursor already reads the repo. |
| [Anil-matcha/Open-Generative-AI](https://github.com/Anil-matcha/Open-Generative-AI) | **skip** | Unrelated image/video studio. |
| [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | **skip** | Skill dump. Already on Hated. Never `npx skills add`. |

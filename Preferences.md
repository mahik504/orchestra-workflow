---
title: Preferences
---

# Preferences

Beats generic skills. Global until skip orchestra.

## Taste

- Not vibe-coded. Not Inter + purple + glow + nested cards.
- Specific to the product. Template look = fail. One visual world per product, not a soup of every trend.
- 3D is a feature, not wallpaper.
- Motion earned (Emil). AI must do real work.
- Lead with the point. No throat-clearing (“Here’s the thing,” “Let me be clear”).
- No binary-contrast filler in product copy (“It’s not X. It’s Y.”).
- No faux-insight or puffery (“What nobody tells you,” “pivotal moment,” “a testament to”).
- Concrete details over abstractions. Keep their voice; do not smooth it into generic AI English.

## Visual world (pick **one** in the plan brief)

Hiring-safe first: **swiss** · **editorial** · **apple-ios** · **bento**.  
Accents only (never the whole app): bauhaus geometry, glass, clay, neo. Mixing all of them is a fail.

## Motion (only where it helps)

Use on marketing / hero / portfolio. Skip on dense tools and dashboards except a **smooth loader**.

Allowed when earned: hover, entrance reveal, micro-interaction, stagger / parallel list, parallax (gentle), 3D hero if Plan said 3D.

Web motion primitives: [React Bits](https://reactbits.dev/) when implementing, not a second design system. Emil still reviews the feel.

## Thinking (how I decide — append when I explain myself)

- I guide in human language. Cursor Plans. I follow. Do not dump a spec and disappear.
- Weak on frontend UI/UX, security, CI/CD. Backend from Cursor/Claude is already decent.
- Unique product look matters. Generic AI SaaS is a fail.
- Six months = internship then job. Shipped public repos beat four unfinished apps.
- Tools cannot DM each other. Vault packets are the handshake.
- College due-now GUI is a workflow test: skip Stitch, ship real git.
- I want studio-site craft (hover, reveal, micro) on surfaces that are *shown*, not on every form.
- College vs hiring must stay labeled. Showing the vault to someone should not mix AirLens with AstroVerse.
- The brain learns by appending Liked/Hated/Thinking the same turn — not by dumping chats.
- Read START HERE first, not the whole vault. Search the web only when asked or a Plan needs a current version/CVE.

## Liked (append when something is right)

- AirLens lamp-black / rust desktop look (they called a later pass very close to 10/10)
- AirLens map: click for a bright mark at the exact point plus that city’s AQI, then a real GitHub push
- AirLens India silhouette under the city dots, viva deck in the same lamp-black / rust
- Echo-allowed studios (materials / motion, **not** a clone): [Nothin'](https://www.noth.in/) editorial-clarity, [K95](https://k95.it/en) grid + case-study density, [Wairk](https://wairk.fr/) sharp AI-studio editorial, [Lax Space](https://www.laxspace.co/) 3D-orbit playground

## Hated (append when something is wrong)

- Fake GitHub contribution graphs / dummy commits / typo-a-day green squares
- Skill dumps (frontend-design, obra/superpowers, vercel skill packs, ui-ux-pro-max, addyosmani/agent-skills)
- Cloning Notion / Linear / a studio site as *our* product UI
- AirLens first polish: generic slate/blue FreeSimpleGUI dashboard
- OpenCode / Kilo / OmniRoute **inside** Cursor or Antigravity (would fight the conductor loop)
- OpenHands / Dify / Langflow / Coolify / Maxun as a second conductor or parallel factory

## How to work

understand → inspect repo → reason → implement → test → review.  
Repo beats vault. No rewrite for sport.

Before adding code: skip it if it is not needed; reuse what this repo already has; prefer stdlib / native / an installed dependency over a new library; then the smallest thing that works. Never skip validation, security, or accessibility to look “minimal.”

## Career

See `memory/career.md` on this PC (private backup). Six months: internship via **shipped public repos**, not four unfinished apps. Public hiring repos: AstroVerse, AirLens (college).

## Core products (do not mix looks)

Hiring focus: **AstroVerse**, then **ODYSS or portfolio**. YUMIT parked until a path exists. AirLens is **college**, not the internship queue.

| Product | Kind | Vault | Disk |
| --- | --- | --- | --- |
| AstroVerse | hiring-cv | `projects/astroverse/` | `C:\projects\AstroVerse` |
| ODYSS | hiring-cv | `projects/odyss/` (unpublished; private backup only) | `C:\projects\ODYSS AI` |
| Portfolio + penfight | hiring-cv | `projects/portfolio-penfight/` (unpublished) | none |
| YUMIT | personal | `projects/yumit/` (unpublished) | missing |
| AirLens | college | `projects/airlens/` | `C:\projects\DVP` (GitHub: AirLens) |
| Evocentric | personal | — | `C:\projects\Evocentric` (not core) |
| Ultron orb | personal | — | `C:\projects\ultron-orb` (parked) |

## Stack (global)

- Web: TypeScript, React, Vite or Next (whatever the repo is)
- Android + later iOS: **Expo** (Router, EAS when we ship Play Store)
- 3D: R3F + drei, or Spline
- 2D games: canvas + Rapier/Matter
- Per-product stack lives in that product’s `idea.md`, not here

Adopted tools (not a wishlist):

- Stitch = 2D chrome. Impeccable/Emil = polish. Context7 = library docs. Playwright = click the UI.
- Official Expo skills (curated). Stitch still owns 2D look; native-ui is native controls after Stitch is locked.
- **Strix** on **our** apps when Docker works (`usestrix/strix` skills already installed). Not a third-party pentest toy.
- **SkillUI** CLI on demand (`npx skillui`) when Plan names **one** live URL to echo. Output goes in `projects/<slug>/ref/`, then we lock *our* `design.md`. Not a global skill pack. Not a clone factory.
- **Manus 1.6 Max** until **2026-08-25**: packet specialist for long cloud browse / wide research. Not conductor. Not a second app factory.
- **Claude Code CLI** (own terminal only): optional think / refine specialist when Plan names it. Point it at a router with **env vars on this PC**. Never install OpenCode, Kilo, OmniRoute, or 9router packs into Cursor or Antigravity.
- Vault **private git** + `kit/sync-vault.ps1` (12h) = Layer-2-lite backup. Public template = `kit/publish-public-vault.ps1`. Not n8n / Vapi / embeddings.
- Catalog of what is installed vs refused: [[STACK]].
- Not adopted: Expo MCP, EAS Observe/Simulator, NativeWind-as-religion, daily dummy commits, n8n/Vapi vector brain, GetLayers templates, ui-ux-pro-max, OpenCode/Kilo/OmniRoute as IDE plugins, OpenHands/Dify/Langflow/Coolify/Maxun, gstack, browser-use, Stirling-PDF, Crawl4AI, code-review-graph MCP, addyosmani/agent-skills.
- Headroom MCP and rtk still deferred (quota is not failing). Token control is vault hygiene: short notes, read START HERE first, router, don’t dump chats. Optional later if quota hurts.

## Anti-references

Generic AI SaaS, skill dumps, fake GitHub graphs, unofficial scrapers, Firebase Studio as the builder, template-library landing pages.

# Orchestra Workflow

An intelligent, provider-neutral control plane for agentic development. One conductor. The same skills on Cursor, Antigravity, Claude Code, and jcode. A human gate before showable UI.

[![CI](https://github.com/mahik504/orchestra-workflow/actions/workflows/ci.yml/badge.svg)](https://github.com/mahik504/orchestra-workflow/actions/workflows/ci.yml)
[![Hygiene](https://github.com/mahik504/orchestra-workflow/actions/workflows/hygiene.yml/badge.svg)](https://github.com/mahik504/orchestra-workflow/actions/workflows/hygiene.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/mahik504/orchestra-workflow?sort=semver)](https://github.com/mahik504/orchestra-workflow/releases)

The pin in `VERSION` is **4.1.0**.

```
ORCHESTRA  = CONTROL PLANE
SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES
AGENTS     = EXECUTORS
BRAIN      = MEMORY
REGISTRY   = RESOURCE KNOWLEDGE
```

This repository is the method. Products, keys, and taste stay in a private workspace. Clone it, wire the hosts you use, and talk.

## What it does

A coding agent will load every skill it can see. Orchestra names the skills, MCPs, plugins, and libraries for the job, and leaves the line off when the job loads nothing new. A skill you name wins unless it is unsafe next to the route that already owns the work. A missing tool stops. It does not install itself.

Showable UI waits for a design contract and golden stills. A markdown contract is not visual evidence. Bugs, APIs, and college work skip that gate unless you say otherwise.

```
Understand → Classify → Load the route → Design Lab or a technical plan → HUMAN GATE
  → Implement → Verify on the running app → Review twice → Remember
```

![The human talks to one conductor](docs/diagrams/assets/control-plane.svg)

![Eight stages, with a human gate before implementation](docs/diagrams/assets/the-loop.svg)

The engine keeps a mission, a worker capsule, and a budget derived from the task. Ordinary work stays small. A flagship page can take a few directions, then one contract. Playwright checks the running app. The model dropdown stays yours. Orchestra does not map a model name onto a skill.

## What's new in 4.1.0

4.0.0 already had the eight stages, the capability graph, the design gate, and host sync. 4.1.0 adds selection on top of that.

- The routing note lists what the task loads: skills, MCP, plugins, libraries, and a stack when one was chosen.
- Design domains annotate the existing routes: premium sites, editorial, shops, portfolios, 3D, shaders, SaaS, dashboards, and mission-control. Responsive layout is a gate on the parent route.
- Technology rules keep one scroll owner, one scene, one HTTP stack, and one bundler. Redis, a queue, Supabase, Lenis, Theatre.js, and a scroll rig are selected when the brief needs them.
- `templates/contracts/` is the design pack another agent can implement from. A bug, a typo, or one control does not generate it.
- Review lenses are requirements, build, visual, architecture, security, performance, and simplify. A typo stays with the builder. A visual judge needs a screenshot from the run.

No new skill was added. The canonical set stays 73, under a cap of 80. No new package was added to this repo.

## Architecture

```mermaid
flowchart TD
  human[Human]
  conductor[One conductor]
  graph[Capability graph]
  skills[Skills MCP plugins libraries]
  hosts[Cursor Antigravity Claude Code jcode]
  brain[Private Brain]
  human --> conductor
  conductor --> graph
  conductor --> skills
  conductor --> hosts
  conductor --> brain
```

| Piece | Where |
| --- | --- |
| Contract | [`AGENTS.md`](AGENTS.md) |
| Adapters | [`.cursorrules`](.cursorrules), [`CLAUDE.md`](CLAUDE.md), [`kit/antigravity/`](kit/antigravity/) |
| Routes | [`registries/design-resource-graph.json`](registries/design-resource-graph.json) |
| Design domains | [`registries/design-domains.json`](registries/design-domains.json) |
| Technology rules | [`registries/technology-policy.json`](registries/technology-policy.json) |
| MCP call rules | [`registries/mcp-contracts.json`](registries/mcp-contracts.json) |
| Resource catalog | [`registries/resources.json`](registries/resources.json) |
| Engine | [`runtime/`](runtime/) |
| Skills | [`skills/`](skills/) |
| Protocols | [`protocols/`](protocols/) |
| Figures | [`docs/diagrams/`](docs/diagrams/) |

A large job becomes capsules with one owner: planner, builder, or judge. The acceptance line moves when you change it. Retry needs a remaining budget and a change. The same failure replayed unchanged stops. Specification conflicts and missing credentials stop for you.

A resource moves one step at a time: discovered, installed, available, selected, executed, verified, successful. A row in the catalog is discovered. Listing it is not a record that it ran.

## Routes

The graph has 26 capabilities. Three that must stay distinct:

| Job | Route | What loads |
| --- | --- | --- |
| A public brand page | `premium-website` | impeccable, taste-design, visual-golden-still, GSAP |
| A school or metrics dashboard | `saas-dashboard` | backend-architecture, test-driven-development |
| A spatial scene | `3d-portfolio` | r3f-shaders |

Other jobs have their own row in `AGENTS.md`: a bug, a paper, a deck, one control, a diagram, a TypeScript review, a deployed app, recent public sources, marketing copy, a rendered video, career notes. One control loads no specialist. The component order is the components the project already has, then one catalog, then GSAP only for scroll.

## Design

Premium and experimental visual work writes `DESIGN.md` from [`templates/contracts/`](templates/contracts/), then golden stills, then code. Open-ended premium work uses 23 short cards. A named reference gets three translations. A pasted contract skips the survey. Stills are still required.

The contract covers brand, type, color, space, breakpoints, components, motion, scroll, cursor, 3D, shaders, accessibility, and the budgets for that page. Principles come from references. Logos, copy, and source files stay with the reference.

Defaults the premium route refuses: Inter or Space Grotesk on a creative headline, an equal three-column card grid, a purple gradient as the brand, repeated glass, decorative blobs, and a 3D object that does not explain the product.

## Frontend and 3D

| Need | Choice |
| --- | --- |
| A document | HTML and CSS |
| A small script | JavaScript |
| A production app | TypeScript |
| Components | React |
| Server rendering and SEO | Next.js |
| A client bundle you named | Vite |
| A content site | Astro |
| Scroll, pin, scrub | GSAP |
| Hover, layout, drag | Motion for React |
| Smooth scroll, when asked | Lenis, one scroller |
| A scene in React | React Three Fiber, with Three.js as that renderer |
| Controls for that scene | Drei |
| A custom GPU effect | GLSL |
| The deliverable is a Framer file | Framer, separate from the code stack |

GSAP owns scroll. Motion does not drive the same container. One canvas. Device pixel ratio stays between 1 and 2 unless the contract says tighter. The scene disposes on unmount. Reduced motion and a failed WebGL context keep a local still. WebGL 2 is the target when a scene needs it. This repo does not target a WebGL 3.

The Next.js starter at [`templates/starters/next-react`](templates/starters/next-react) has lint, typecheck, build, and a Lighthouse job that warns.

## Backend and data

Write the shape before the handlers.

| Store | When |
| --- | --- |
| Postgres | Relational data that must survive a restart |
| Supabase | The product note already names it |
| Redis | A cache, a TTL, or a rate limit |
| Express | A small HTTP API the repo does not already serve |
| FastAPI | Python inference or data processing |
| BullMQ | A background job that already needs Redis |
| Node | The TypeScript server |

Say the tables, keys, constraints, indexes, who can call each route, how input is rejected, and which fields are never logged. `supabase` and `postgres-best-practices` load when that product uses them.

## Mobile, research, and documents

Expo skills cover a React Native app: routing, native UI, data fetching, dev client, upgrade, and store submit. An iOS device build needs a Mac and Xcode. Android uses the SDK on the machine that builds it. Swift, Kotlin, and Flutter are outside this skill set.

A paper loads `research-paper` and `humanizer`. Citations are checked. A deck loads `orchestra-docs` and `orchestra-presentation`. Host tools for docx, pptx, pdf, and xlsx stay on the host. They are not copied into `skills/`.

An ML fine-tune loads `ml-experiment` and reports the eval you actually ran. `rea` runs on a binary or APK only when Java is installed and the target is yours, a client's, a CTF, or open source.

## Quality

Playwright walks the running app. The canonical check is the harness in `runtime/` with the project Playwright library and Chromium. A mock screenshot does not pass. Lighthouse is the performance budget for that project. axe-core is the violation list on the running page. Golden stills are the visual contract. Later changes compare against those stills, with masks for clocks and canvases that are not the subject.

Two review passes stay separate: correctness, then simplify.

## Brain

Set `ORCHESTRA_VAULT` to your private notes. The vault holds taste, decisions, and one lesson line per finished session. This clone does not contain it. Durable decisions use [`templates/memory/decision.md`](templates/memory/decision.md). A catalog name is not a result.

## Hosts

![One contract, synced to the hosts you pick](docs/diagrams/assets/four-hosts.svg)

| Host | Role |
| --- | --- |
| Cursor | Usual conductor. Adapter: `.cursorrules` |
| Antigravity | Visual QA. Adapter: `kit/antigravity/` |
| Claude Code | Terminal and review. Adapter: `CLAUDE.md` |
| jcode | Same skills |

`kit/sync-ides.ps1 -Sync` copies `skills/` to Cursor, Antigravity, Claude Code, jcode, and the shared agent skill folders on this machine. It does not install MCP servers. Sign in to GitHub, Playwright, and Context7 in the host UI. Tokens stay there.

Core MCP templates are in [`registries/host-stack.json`](registries/host-stack.json): vault files, Stitch, Playwright, Context7. Optional servers stay optional. Firecrawl and Serena are called when a job already enabled them.

### Limitless Claude

[Limitless Claude](https://github.com/mahik504/limitless-claude) is a separate repository. The published line is 7.0.0. It routes model providers for Claude Code and compatible hosts through a local OmniRoute gateway. Orchestra does not pick the vendor model and does not read that gateway's database. Credentials for providers stay in that project. Use it when you want provider failover on Claude Code. Orchestra still decides which skills load.

## Install

Git, and Go 1.22 if you want the engine. Node is for the hygiene checks and the Next starter.

```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow
```

Windows:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/bootstrap.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync
```

macOS or Linux:

```bash
chmod +x kit/bootstrap.sh kit/install-local-engine.sh
./kit/bootstrap.sh
```

Set `ORCHESTRA_VAULT` to your private notes folder. Restart the IDE. The next chat in this repo reads `AGENTS.md`.

Walkthrough: [`docs/getting-started.md`](docs/getting-started.md).

Say **skip orchestra** to stand down for the session. Say **skip the lab** to bypass the design gate for one task.

### First task

Open the product repo. Say what it is and who it is for. The conductor classifies one route and names what it loads. For a public page, approve the contract and the stills before implementation. For a bug, it loads `systematic-debugging` and skips the lab.

## Test

From `runtime/`:

```bash
go test ./...
```

From the repo root:

```bash
node runtime/tools/hygiene.js
node runtime/tools/check-registry.js
node runtime/tools/check-routing.js
node runtime/tools/check-41.js
```

CI on `main` runs `go test`, `gofmt -l`, and schema JSON checks. Hygiene runs as its own workflow. A design-contract directory can be checked with:

```bash
node runtime/tools/check-design-contract.js templates/contracts spatial
```

## Security

Secrets stay in the host. This repo's hygiene check rejects credential-shaped strings, private vault paths, and personal product names in tracked files. Global package installs from the engine are blocked. Offensive scanning is for an app you own. `ship-safe` loads when a change touches auth, public routes, or secrets.

## Migration from 4.0.0

Pull this tag. The capability ids are unchanged. Old resource rows still load. Sync skills when you want the hosts to match:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync
```

Details: [`kit/ROLLBACK.md`](kit/ROLLBACK.md).

## Rollback

`ORCHESTRA_CONTRACT` changes the label the doctor prints. The 4.1 selector is code in this tree. To restore 4.0.0 behavior, check out that commit and sync skills from it. Do not force-push `main`.

```powershell
$env:ORCHESTRA_CONTRACT = "4.0.0"
```

```bash
export ORCHESTRA_CONTRACT=4.0.0
```

## Contributing

Useful changes are install fixes, a registry row with a license, and a test around the design gate or the routing check. Leave out skill dumps, secrets, and personal notes. See [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Limits

- One conductor. A second plan loses.
- The model dropdown is yours.
- Visual QA needs a running app.
- An MCP you have not signed in to is not healthy.
- GetLayers is purchase-pending.
- iOS device builds need a Mac.
- The saved one-line Go comparison made the same 183-byte change, and Orchestra took longer. That pair does not score a website.

## License

MIT. See [LICENSE](LICENSE).

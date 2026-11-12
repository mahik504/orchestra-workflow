# Capability map 5.5

Recorded 2026-11-12. A discipline is supported when a route or an engineering row already names a skill that exists on disk. “Supported” is not “measured.” Stub skills are marked. New skills are not proposed where a route already points at a file.

## Software engineering

| Job | What the table loads | State |
| --- | --- | --- |
| Bug | `systematic-debugging` | Vendored file exists. Fast path |
| Logic and TDD | `test-driven-development` | Vendored. Also the engineering row |
| Before calling it done | `verification-before-completion` | Vendored. Engineering row, not a keyword route |
| Diff review | `ponytail` | Engineering row. Not for a Design Lab surface |
| Vague new product | `grilling`, then `domain-modeling` | Local folders exist. Upstream now lists `grill-with-docs`. Not diffed |
| Work breakdown | `to-tickets` | Local folder exists |
| Quality bar | `constraint-driven-development` | Named in the engineering table. Folder not re-opened |
| Review feedback | `receiving-code-review` | Vendored from Superpowers |
| Spec | `egm-draft`, `egm-review`, `egm-resolve` | Route exists. STANDARD fix that names the file skips it |
| TypeScript | `ow-typescript` | Stub |
| Backend / SaaS | `ow-backend-architecture`, `test-driven-development` | Backend skill is stub. `fullstack-mastery` is on disk and off the sync list on purpose |
| Full-stack | backend skill, `react-best-practices`, TDD, `webapp-testing` | Route exists. Design lab skipped |
| Next.js speed | `react-best-practices` | Route exists |
| Security-sensitive ship | `ow-ship-safe`, `ow-observability` | Both stub |
| Deploy | the ship route above. `ow-deploy` exists and is stub, and it is not the ship route’s load list | Do not add a second deploy skill |
| CI and supply chain | `ci-security-scanning-with-strix`, `ow-semgrep` | Strix is a skill folder. Semgrep skill has a procedure flag off stub. Neither was run |
| Codebase map | `floorplan`, `ow-diagrams` | Diagrams skill is stub |
| Parallel sessions | `using-git-worktrees`, `ow-wave` | Wave has a procedure. Worktrees are vendored |
| Refactor, performance, legacy, microservices, auth | No dedicated route | Use the backend or TypeScript route plus the task text. A skill per noun is how the catalog gets large without a trial |

## Frontend and interface

| Job | Load | State |
| --- | --- | --- |
| Award-style site | `build-awwwards-quality-sites`, `taste-design`, `impeccable`, `ow-golden-still`, GSAP core, ScrollTrigger, timeline | Golden-still skill is stub. Motion skills are the official GSAP pointers |
| Unslop | `unslop`, `impeccable` | Route exists |
| Echo a layout without copying the brand | `ow-reference-forensics`, `impeccable`, `taste-design` | Forensics skill is stub |
| One control | `none` | Fast path |
| UI assembly | `ow-ui-assembly` | Has a Procedure heading. Not graded |
| Site recon | `ow-site-recon` | Procedure heading. Skip line forbids copying assets, copy, and logos |
| One 3D scene | `ow-3d-engineering` | Stub. `creative-3d-web` stays unloaded |
| Scroll | GSAP trio | One motion runtime |
| Accessibility and browser QA | `ow-qa`, `ow-a11y-auditor` agent | `ow-qa` procedure was opened. It is not a route row. The changelog says it is. That is baseline F-07 |
| Motion decision | `ow-motion-engineering` | Stub. `ow-motion-director` is an agent file |

Component order already stated in older changelog text: the product’s components, then one registry source, not every kit. `ow-component-discovery` is stub. shadcn, Magic UI, and the other kits are resources in the registry, not skills to load together.

## Mobile

Route “Mobile app” loads the official Expo skills: project structure, router, native UI, data fetching, dev client. A marketing site is a negative match. Store release is `ow-store-release`, stub, and not on that route. EAS is a separate skill folder (`eas-app-stores`). Do not load it for a college web page.

## Data, ML, and research

| Job | Load | State |
| --- | --- | --- |
| What people said recently | `ow-research` | Stub. No scraping, no posting |
| Paper | `ow-research-paper`, `humanizer` | Research-paper skill is stub |
| ML experiment | `ow-ml-experiment` | Stub. The skip line forbids a claimed benchmark win |
| Essay | `humanizer`, `ow-research` | Route exists |

Classical ML, vision, NLP, fine-tuning, and MLOps do not get a skill each. The experiment skill, when it has a procedure, should say: data you have, a metric you compute, no invented accuracy. That is the current skip line, expanded into steps later.

## Presentations and documents

| Job | Load | State |
| --- | --- | --- |
| Deck | `ow-docs` | Stub |
| Deck with speaker notes | `ow-docs`, `ow-presentation` | Both stub |
| Docs on disk | `host-provided` | The route says docx, pptx, pdf, and xlsx are not in `skills/`. They are. The license problem is in the security document. The route’s intent, host tools rather than a vendored Anthropic skill, is the right one |
| README and ship post | `ow-ship`, `humanizer` | Ship skill is stub |
| Diagram | `ow-diagrams`, `pretty-mermaid` | Diagrams skill is stub |
| Video | `ow-video` | Stub. Not the GSAP route |

## Recon, scrape, and automation

| Job | Load | State |
| --- | --- | --- |
| App recon | `ow-app-recon` | Procedure heading. Owned, client, CTF, or open source only |
| Scrape | `ow-scrape` | Stops at login, paywall, CAPTCHA, robots, 403, 429. LinkedIn and Instagram are negative matches |
| Browserize a finished session | `ow-browserize` | Procedure heading. Not a route row |
| MCP server | `mcp-builder` | Route exists |
| MCP audit | `ow-mcp-audit` | Procedure heading. Does not upload tool text without a token |

Desktop automation beyond what the host already allows is not a new skill. The hook remains the denylist.

## Career, marketing, and product judgment

`ow-career`, `ow-marketing`, `ow-brainstorm`, `ow-coach`, `ow-showcase` are stubs. Brainstorm is the pressure test before a new product. The user rule is one flagship, not four builds. The skill should stay a note, not a generator of repos.

## Agents

Two review agents are the contract: `ow-reviewer`, then `ow-simplifier`. The others are lenses for a route that needs them: security, accessibility, visual, motion, recon, scrape, eval, product critique, and the EGM spec lenses. A lens that is not named by the route stays unloaded. This session did not measure whether each file changes a diff.

## How a frontend decision should go

This is a procedure for the existing skills, not a new skill.

1. Read the product’s current components and tokens.
2. If the task names one control, stop. Build that control in the existing style.
3. If the task is a screen inside a product, search that repo first, then one component registry.
4. If the task is a new visual identity, look at references for structure only. `ow-reference-forensics` and `ow-site-recon` forbid copying assets, copy, and logos.
5. Motion only when the route is scroll, a transition, or a scene. One runtime. GSAP for scroll. No second engine.
6. 3D only when the task is one scene. Still image fallback. `ow-3d-engineering`, not `creative-3d-web`.
7. Browser QA when the running app can be opened. `ow-qa`. A Lighthouse score is not a design judgment.
8. PREMIUM stills wait for a human yes.

Godly, Awwwards, Dribbble, Behance, Mobbin, Refero, SiteInspire, and Page Flows are reference shelves. They are not skills. Refero stays auth-gated in the older changelog. This pass did not re-check those sites.

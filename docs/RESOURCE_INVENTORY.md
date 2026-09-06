# Resource inventory — Orchestra 3.3.0

Available ≠ loaded. Registered ≠ active. Discovered ≠ verified. Machine-readable rows: `registries/resources.json`. Local OSS metadata (not clones): `registries/oss-index.json`. Health: [`RESOURCE_HEALTH.md`](RESOURCE_HEALTH.md). Combination: [`DESIGN_RESOURCE_ROUTING.md`](DESIGN_RESOURCE_ROUTING.md).

Every catalog row has **trigger** and **skip**. Re-check live prices; do not freeze a checkout in a README.

GetLayers Full Stack lifetime was **$199** on 2026-09-06 at [getlayers.ai/pricing](https://www.getlayers.ai/pricing) (Unlimited lifetime **$139**). MCP requires Full Stack. Status: `PURCHASE_PENDING` / `BLOCKED_UNTIL_PURCHASE`.

Verification date for this table: **2026-09-06**.

## A. ACTIVE SKILLS

40 curated skills in `registries/host-stack.json`. Copied by `kit/sync-ides.ps1` to Cursor, Antigravity, Claude Code, Agents, and jcode.

| Name | URL | Type | Representation | Free capability | Paid capability | Install | Auth | License | Trigger | Skip | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| orchestra-conductor | this repo | skill | skill | control plane | none | sync-ides | NONE | MIT | every Orchestra job | skip orchestra | ACTIVE |
| taste-design | wrapper in this repo | skill | skill | DESIGN.md wrapper | Taste v2 pack not installed | sync-ides | NONE | MIT | visual direction | backend-only | ACTIVE |
| refero-design | https://github.com/referodesign/refero_skill | skill | skill | methodology + public pages | Pro MCP | sync-ides | NONE | site-terms | product UI research | MCP-only jobs without Pro | ACTIVE |
| hallmark | https://github.com/Nutlope/hallmark | skill | skill | post-implement audit | none | sync-ides | NONE | see-upstream | after stills | using it as the generator | ACTIVE |
| unlazy | https://github.com/Leonxlnx/unlazy | skill | skill | evidence/completion method | stop-hooks not installed | sync-ides | NONE | see-upstream | verification ledger | second completion authority | ACTIVE |
| react-best-practices | https://github.com/vercel-labs/agent-skills | skill | skill | React/Vercel practices | rest of the Vercel pack | sync-ides | NONE | see-upstream | React implementation | bulk Vercel dump | ACTIVE |
| web-design-guidelines | https://github.com/vercel-labs/agent-skills | skill | skill | web design guidelines | rest of the pack | sync-ides | NONE | see-upstream | showable web UI | backend-only | ACTIVE |
| gsap-core | https://github.com/greensock/gsap | skill | skill | GSAP core motion | Club plugins if any | sync-ides | NONE | see-upstream | timeline motion named | stacking every animation lib | ACTIVE |
| supabase | https://github.com/supabase/agent-skills | skill | skill | relevant Supabase areas | unused pack skills | sync-ides | NONE | Apache-2.0 | project uses Supabase | non-Supabase apps | ACTIVE |
| postgres-best-practices | https://github.com/neondatabase/postgres-skills | skill | skill | schema/index/query | rest of neon pack | sync-ides | NONE | see-upstream | Postgres/Supabase schema | no database | ACTIVE |
| diagram-generator | https://github.com/jovd83/diagram-generator | skill | skill | grounded Mermaid authoring | none | sync-ides | NONE | see-upstream | architecture docs | invented architecture | ACTIVE |
| pretty-mermaid | https://github.com/imxv/Pretty-mermaid-skills | skill | skill | render Mermaid | none | sync-ides | NONE | see-upstream | after diagram-generator | no Mermaid source | ACTIVE |
| skill-security-check | https://github.com/aliksir/claude-code-skill-security-check | skill | skill | admission review | none | sync-ides | NONE | see-upstream | new external skill | ordinary coding turn | ACTIVE |
| stitch-* / expo-* / strix / impeccable / emil / animate / ship-safe / orchestra-* | this repo + upstream | skill | skill | host adapters and job skills | none extra | sync-ides | NONE | mixed | matching job | wrong platform | ACTIVE |

TOTAL ACTIVE CURATED SKILLS = **40**.

Why these 40 and not 30 or 88: the original 30 stay. Ten additions fill real gaps (Refero method, Hallmark audit, Unlazy ledger, Vercel React, GSAP core, Supabase, Postgres, diagrams, Pretty-Mermaid, skill admission) without a second conductor. The “40” is not a target. Duplicates are not added.

## B. OPTIONAL SKILLS

| Name | URL | Representation | Trigger | Skip | Status |
| --- | --- | --- | --- | --- | --- |
| frontend-design (Anthropic) | https://github.com/anthropics/claude-code | reference | Claude host already has it | do not vendor the marketplace tree | REFERENCE |
| Soup LLM | catalog row `soup-llm` | project_tool | ML fine-tune pack | not a 41st global skill | CURATED_OPTIONAL |

## C. FREE MCPs

| Name | URL | Free capability | Install | Auth | Trigger | Skip | Status | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Magic UI MCP | https://magicui.design/docs/mcp | OSS catalog via `@magicuidesign/mcp` | `npx @magicuidesign/cli@latest install cursor` or example JSON | HOST_CONNECT | named Magic UI component | always-on / Pro claim | OPTIONAL | Pro not claimed |
| Playwright MCP | `@playwright/mcp` | our-app browser | host MCP | HOST_CONNECT | visual QA on our app | foreign-site ops | CORE | |
| Context7 | `@upstash/context7-mcp` | library docs | host MCP | HOST_CONNECT | library docs | guessing APIs | CORE | |
| shadcn MCP | `@shadcn/mcp-server` | component add | host MCP | HOST_CONNECT | named shadcn primitive | operator HUD | OPTIONAL | Cursor already has it on this machine |

## D. AUTHENTICATED MCPs

Configured possible; remaining step is host login. Not HEALTHY until the host confirms.

| Name | Host | Auth type | Manual action | Expected state |
| --- | --- | --- | --- | --- |
| Stitch | Cursor / Antigravity | API key in host | paste key in host UI | HEALTHY after host confirms |
| Mobbin | Antigravity (and Cursor if added) | account | Connect in MCP UI | HEALTHY after Connect |
| GitHub | per host | OAuth | Connect | HEALTHY after Connect |
| Firebase | Antigravity only | Google | Connect for that job | AUTH_REQUIRED until signed in |
| Gmail / Calendar / Drive / Stripe / Zapier | Cursor plugins | OAuth | Connect in Cursor | Cursor only |

## E. AUTH_REQUIRED MCPs

| Name | URL | Free capability | Paid capability | Status |
| --- | --- | --- | --- | --- |
| Refero MCP | https://api.refero.design/mcp | none on the MCP | Pro screen search | AUTH_REQUIRED / PRO |
| 21st.dev MCP | https://docs.21st.dev/mcp | catalog search; 2 installs/day after login | Builder / Builder+AI / extra installs | OPTIONAL_FREE_MCP + AUTH_REQUIRED until login |
| Tailkit MCP | https://tailkit.com/ | none on the MCP | licensed MCP | AUTH_REQUIRED / PAID_ACCESS |
| Firecrawl | host | none without key | keyed crawl | OPTIONAL + key in host |
| Serena | https://github.com/oraios/serena | OSS MCP once installed | none | OPTIONAL, never always-on |
| Supabase MCP | https://mcp.supabase.com/mcp | none without auth | project access | AUTH_REQUIRED |
| vault-memory | filesystem MCP | local vault path | none | must point at YOUR workspace; never commit live path |

## F. FREE COMPONENT SOURCES

| Name | URL | Free | Paid | Trigger | Skip |
| --- | --- | --- | --- | --- | --- |
| shadcn/ui | https://github.com/shadcn-ui/ui | foundation primitives | none | base components | treating it as visual identity |
| React Bits | https://github.com/DavidHDev/react-bits | OSS motion/components | Pro/Ultimate | named motion kit | buy Ultimate |
| Magic UI OSS | https://github.com/magicuidesign/magicui | OSS | Pro | named Magic UI | Pro registry |
| Aceternity | https://ui.aceternity.com/ | buttons, backgrounds, text, cards, forms, navigation, effects, 3D, hero | All-Access | named echo / free category | All-Access claim |
| Kokonut UI | https://github.com/kokonut-labs/kokonutui | gradient/magnetic/particle/command/hold buttons, liquid glass, cards, drawers | none claimed | named Kokonut | collage |
| Motion Primitives | https://github.com/ibelick/motion-primitives | transitions / micro-interactions | none claimed | named on card | stacking with GSAP |
| Cult UI | https://github.com/nolly-studio/cult-ui | animated blocks, shadcn-compatible | none claimed | named Cult | always-on MCP |
| HyperUI | https://hyperui.dev/ | Tailwind utility UI, forms, tables, nav | none claimed | utility fallback | PREMIUM 3D |
| Tailkit public | https://tailkit.com/ | public/free pages | MCP/paid kit | public pattern search | treating MCP as free |
| 21st free catalog | https://21st.dev/ | search + 2 installs/day | Builder | component discovery | always-on / quota burn |

One capability: `interaction-components`. Not one skill per library.

## G. DESIGN RESEARCH SOURCES

| Name | URL | Type | MCP? | Trigger | Skip |
| --- | --- | --- | --- | --- | --- |
| Godly | https://godly.design/ | research | no | high-visual research | bulk scrape |
| Land-book | https://land-book.com/ | research | no | landing research | bulk scrape |
| Lapa Ninja | https://www.lapa.ninja/ | research | no | landing gallery | bulk scrape |
| SaaSFrame | https://www.saasframe.io/ | research | no | SaaS UI research | bulk scrape |
| Refero public | https://refero.design/ | research | MCP is Pro | public styles / DESIGN.md refs | cloning screenshots |

Orchestra may say: go find a relevant reference from the registered design research sources.

## H. FREE PROMPT SOURCES

| Name | URL | Trigger | Skip |
| --- | --- | --- | --- |
| awesome-chatgpt-prompts | https://github.com/f/awesome-chatgpt-prompts | one targeted prompt in UI / frontend / architecture / coding / agents / design | dump whole list into context |
| awesome-prompts | https://github.com/ai-boost/awesome-prompts | second catalog if the first misses | dump whole list |

## I. 3D SOURCES

Project-scoped npm only. Never global frontend installs.

| Name | URL | When |
| --- | --- | --- |
| three.js | https://github.com/mrdoob/three.js | 3D / WebGL / canvas / shader |
| R3F | https://github.com/pmndrs/react-three-fiber | React 3D |
| drei | https://github.com/pmndrs/drei | R3F helpers |
| react-postprocessing | https://github.com/pmndrs/react-postprocessing | cinematic 3D |
| r3f-scroll-rig | https://github.com/14islands/r3f-scroll-rig | scroll-driven 3D |
| GetLayers | https://www.getlayers.ai/ | **PURCHASE_PENDING** — future adapter only |

Triggers: 3D website, WebGL, interactive 3D, portfolio, spatial experience, shader, canvas, 3D interaction.

## J. BACKEND / DATABASE SOURCES

| Name | Trigger | Skip |
| --- | --- | --- |
| supabase skill | project uses Supabase (db, auth, storage, realtime, edge, RLS, migrations, vectors, SSR) | no Supabase |
| postgres-best-practices | Postgres / schema / index / query | no database |
| Supabase MCP | authenticated project | AUTH_REQUIRED |

## K. DEVOPS SOURCES

| Name | URL | Representation | Notes |
| --- | --- | --- | --- |
| devops-skills | https://github.com/arjunprabhulal/devops-skills | quarantined_reference | Do not install all 88. Themes only: CI/CD, Docker, deploy, observability, security, reliability, infra, performance |
| ship-safe / strix / semgrep-adapter | this repo | skills | Prefer these over dumping devops-skills |

## L. SECURITY SOURCES

| Name | Trigger | Skip |
| --- | --- | --- |
| skill-security-check | new external skill admission | ordinary coding |
| strix / semgrep-adapter | owned-app security jobs | scanning systems you do not own |

## M. ARCHITECTURE / DIAGRAM SOURCES

| Name | Role |
| --- | --- |
| diagram-generator | architecture analysis → Mermaid, grounded in the repo |
| mermaid | diagram language |
| pretty-mermaid | optional render |

Eight required diagrams live under `docs/diagrams/01-control-plane.md` … `08-resource-combination.md`. Older numbered notes remain as extra detail.

## N. CLIs

| Name | Role | Skip |
| --- | --- | --- |
| Scrapling | acquire / extract public pages | prohibited bulk scrape |
| agent-browser | interact with foreign websites | `--tools all`; replacing Playwright on our app |
| Playwright | verify our application | foreign-site browsing as the primary tool |
| Magic UI CLI | `npx @magicuidesign/cli@latest install <client>` | Pro install |
| 21st CLI | search / limited install after login | exceeding free quota on purpose |
| jcode | executor against an OpenAI-compatible endpoint | second conductor |

## O. EXECUTORS

Cursor, Antigravity, Claude Code, jcode. One conductor per session: the process the human is talking to. jcode is not an MCP and is not a second Orchestra.

## P. PURCHASE_PENDING

| Name | URL | Status |
| --- | --- | --- |
| GetLayers | https://www.getlayers.ai/ | PURCHASE_PENDING / BLOCKED_UNTIL_PURCHASE. No MCP, no scrape, no fake tools, no HEALTHY claim. After purchase, say **update**. |

## Q. QUARANTINED REFERENCES

Do not install: ECC full pack, Superpowers full runtime, Taste v2 entire pack, DevOps 88, Anthropic marketplace dump, Google skill dump, ui-ux-pro-max dump, `npx skills add --all`, ibelick/ui-skills dump.

They may remain catalog `REJECTED` or `quarantined_reference` rows with skip = never load into runtime.

## Local OSS index

`registries/oss-index.json` stores URL, license, date, capability, categories. `local_copy: false`. Proprietary sites are metadata + URL only. This is resource knowledge, not the Brain.

## Brain vs workflow

| Layer | Lives where | Clones get it? |
| --- | --- | --- |
| Contract, graph, skills, protocols, README | orchestra-workflow 3.3.0 | Yes |
| Product notes, career, Preferences | private Brain | No |
| OmniRoute, MCP tokens, host MCP JSON | user profile | **No** |

# Resource inventory — Orchestra 3.3.0

Available ≠ loaded. Registry presence is not usage. Re-check live prices; do not freeze a stale checkout in a README.

GetLayers Full Stack lifetime was **$199** on 2026-09-06 at [getlayers.ai/pricing](https://www.getlayers.ai/pricing) (Unlimited lifetime **$139**). ChatGPT cited $129 — that was wrong that day. MCP requires Full Stack Lifetime.

## A. Global skills (host-stack allowlist)

40 skills in `registries/host-stack.json`. Copied by `kit/sync-ides.ps1` to Cursor, Antigravity, Claude Code, Agents, and jcode.

Keep the original 30. Added only skills that do not start a second loop:

| Skill | Why it does not fight the 30 |
| --- | --- |
| `refero-design` | Research. MCP stays AUTH_REQUIRED |
| `hallmark` | Audit after implement. Impeccable stays generator-side |
| `unlazy` | Evidence. Points at the ledger. No stop-hook takeover |
| `react-best-practices` | One Vercel skill. `web-design-guidelines` already in the 30 |
| `gsap-core` | Official GSAP core. One motion engine at the gate |
| `supabase` | Backend pack trigger, not a second OS |
| `postgres-best-practices` | Schema/index jobs |
| `diagram-generator` | README/architecture Mermaid authoring |
| `pretty-mermaid` | Render only |
| `skill-security-check` | Admission for **new** third-party skills |

## B. Packs (job-scoped markdown)

`templates/packs/` — not extra global skills.

| Pack | When |
| --- | --- |
| website | Public web / SaaS web |
| ios | Expo iOS |
| android | Expo Android |
| research-paper | orchestra-docs |
| ml-finetune | Soup **project-scoped**, CURATED_OPTIONAL |

## C. MCP states

`HEALTHY` / `OPTIONAL` / `AUTH_REQUIRED` / `BROKEN` / `DISABLED` / `PURCHASE_PENDING`.

| Priority | Servers |
| --- | --- |
| Core | vault-memory, stitch, playwright, context7 |
| P0 later | GetLayers (purchase), Refero (AUTH), Serena (optional), Scrapling |
| P1 when Connect | 21st, Magic UI, Tailkit, agent-browser (foreign sites) |

## D. Paid board (Orchestra fit, not "coolest site")

Do **not** buy React Bits Ultimate / Tailkit / Aceternity / Magic Pro this pass.

| Name | Fit | This pass |
| --- | --- | --- |
| GetLayers Full Stack | 9.7 | Buy later. One name. |
| React Bits Ultimate | 9.4 | No. OSS already in stack |
| Tailkit Unlimited | 9.1 | No |
| Refero Pro | 9.0 | Connect when named |
| Aceternity | 8.8 | Named echo only |
| Magic UI Pro | 8.8 | No. OSS ok |
| 21st Builder+AI | 8.7 | No lifetime. Optional MCP |
| SaaSFrame | 8.0 | Catalog only |
| UI Prompt Library | 7.9 | Catalog only |
| Land-book | 7.6 | Research |
| Godly | 7.5 | Research |
| Lapa | 7.4 | Research |

## E. Free / OSS (register or keep)

shadcn, React Bits OSS, Magic UI OSS, Motion Primitives, Kokonut UI, Cult UI, HyperUI, Three.js, R3F, drei, react-postprocessing, r3f-scroll-rig, GSAP (now free), Scrapling, screenshot-to-code, SkillUI.

**interaction-components** is one capability domain that **routes** across those kits. Not a skill per button library.

## F. GetLayers pending

Status: `PURCHASE_PENDING` / blocked until purchase.

- No MCP connect.
- No fake tools.
- No scraping the paid library.
- After buy Full Stack (~1 week on this operator's plan): say **update** to flip the row.

## G. Ignore list (do not install)

ECC / full Superpowers / anthropics/skills dump / devops-skills all 88 / Taste v2 pack / ibelick/ui-skills / GetLayers MCP now / always-on 21st Magic Aceternity Tailkit MCP / `agent-browser --tools all` / Serena always-on / `npx skills add --all`.

Anthropic `frontend-design`: named Claude-host optional if already present. Do not vendor the tree into public git. Catalog `REFERENCE`.

## H. Host adapters

Cursor, Antigravity, Claude Code, jcode. Same contract. Different syntax. See `docs/manual-setup/`.

## I. Design Lab

23 short cards → one `DESIGN.md`. Named override skips survey. Custom paste wins.

## J. Verification

Ledger + Playwright on our app + Designer SOS. Hallmark after stills. No "zero errors" claim.

## K. 3D

R3F / drei / postprocessing when the card says 3D. Low-end fallback on EXPERIMENTAL. GetLayers 3D layers wait on purchase.

## L. Expo

Official Expo skills in the 40. iOS and Android packs. No unofficial Expo skill packs.

## M. Brain vs workflow

| Layer | Lives where | Clones get it? |
| --- | --- | --- |
| Contract, graph, skills, protocols, README | orchestra-workflow 3.3.0 | Yes |
| Product notes, career, Preferences | private Brain | No |
| OmniRoute, MCP tokens, `~/.claude`, `~/.jcode` | user profile | **No** |

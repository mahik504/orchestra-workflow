# Intelligence lineage — Orchestra 4.1.0

Checked 2026-11-05 with an HTTP HEAD. `200` means the URL responded. `301` means GitHub redirected. `404` is unresolved. Nothing in this file was copied from an upstream skill file.

A star count is not a reason to make that repo the conductor.

## New Orchestra skills

| Skill | Source method | URL check |
| --- | --- | --- |
| orchestra-focus | ADHD-style next action, no repeat questions | `ayghri/i-have-adhd` 200. Caveman stays a warning, not the brief |
| orchestra-presentation | Slide story, type, notes, slide QA | `hugohe3/ppt-master` 200. `anthropics/skills` 200. Host pptx still builds the file |
| orchestra-research | Dated public search and a citation | `mvanhorn/last30days-skill` 200 |
| orchestra-typescript | Type boundaries before a TypeScript edit | `mattpocock/skills` 200. Pack not installed |
| orchestra-observability | Errors, traces, dependency and image scans in the product repo | `getsentry/sentry-javascript` 200 |
| orchestra-video | Rendered mp4 in the product repo | `heygen-com/hyperframes` 200. `Open-Montage-app/OpenMontage` 404, so that name is not an install target |
| orchestra-marketing | Shipping copy without invented metrics | `coreyhaines31/marketingskills` 200 |
| orchestra-career | CV from notes the operator already has. The operator submits | `career-ops-hq/career-ops` 200. The scraper is not installed |

## Already in this repo

| Source | Check | Disposition |
| --- | --- | --- |
| obra/superpowers | 200 | Methods stay in test-driven-development, systematic-debugging, verification-before-completion, using-git-worktrees. The plugin does not start a plan |
| DietrichGebert/ponytail | 200 | ponytail, ponytail-review, ponytail-audit |
| Leonxlnx/taste-skill | 200 | taste-design. Pack not reinstalled |
| pbakaus/impeccable | 200 | impeccable |
| blader/humanizer | 200 | humanizer |
| Panniantong/Agent-Reach | 200 | agent-reach. Own public pages only |
| morluto/rea | 200 | rea |
| vercel-labs/agent-skills | 200 | react-best-practices and web-design-guidelines already point here. The pack is not added |
| microsoft/playwright | 200 | Playwright MCP and the harness |
| GoogleChrome/lighthouse | 200 | Warning job in the Next starter |
| Graphify-Labs/graphify | 200 | Optional map on a repo audit. Not a skill |
| oraios/serena | 200 | Optional symbol tool. Not always on |

## Folded into an existing file

| Source | Check | Where it went |
| --- | --- | --- |
| multica-ai/andrej-karpathy-skills | 200 | ponytail. Think in the brief, smallest complete change, only the requested code, verify before done |
| open-gsd/gsd-core | 200 | Conductor mission. Capsules, one owner, acceptance moves only when the human changes it |
| gsd-build/get-shit-done | 200 | Same mission section. Catalog already treats the old repo as archived |
| open-gsd/get-shit-done-redux | 301 to open-gsd/gsd-core | Not a separate product |
| bmad-code-org/BMAD-METHOD | 200 | Same mission section. Not a second conductor |
| garrytan/gstack | 200 | fresh-context-review and ponytail-review. Not a startup-team simulator |
| JuliusBrussee/caveman | 200 | orchestra-focus. Short output does not replace the brief |
| nextlevelbuilder/ui-ux-pro-max-skill | 200 | taste-design and impeccable. Not a third design system |
| cathrynlavery/diagram-design | 200 | diagram-generator |
| RiantRedwood/Archify | 200 | diagram-generator. Reference only. Mermaid stays the renderer |
| affaan-m/everything-claude-code | 301 to affaan-m/ECC | Read as a pack. Not installed. No extra skill until a named gap shows up |
| addyosmani/agent-skills | 200 | Review habits already in ponytail-review |
| wshobson/agents | 200 | Role ideas already in fresh-context-review. Collection not installed |
| headroomlabs-ai/headroom | 200 | Not installed. Worker capsules already limit what a worker sees |
| Egonex-AI/Understand-Anything | 200 | Optional reference. ponytail-audit plus Graphify stays the map |
| anthropics/skills document skills | 200 | Host docx, pdf, pptx, xlsx. Not copied into skills/ |

## Tools, not skills

| Source | Check | Disposition |
| --- | --- | --- |
| ui.watermelon.sh and WatermelonCorp/watermelon-platform | both 200. GitHub SPDX MIT | REFERENCE row `watermelon-ui`. No MCP |
| microsoft/TypeScript | 200 | The compiler. orchestra-typescript is the instruction |
| clerk/javascript | 200 | The product note picks auth. Not a global skill |
| Deque axe-core, Vitest, Bundlephobia | known public tools | Invoked from web-design-guidelines, test-driven-development, and react-best-practices. Not new skill folders |
| RiantRedwood/Archify | 200 | REFERENCE row `archify`. Mermaid stays the renderer |
| Egonex-AI/Understand-Anything | 200 | REFERENCE row `understand-anything`. ponytail-audit plus Graphify stays the map |
| microsoft/markitdown | 200 | REFERENCE row `markitdown`. orchestra-docs can use it for a conversion |
| DataExpert-io/data-engineer-handbook | 200 | REFERENCE row `data-engineer-handbook` |
| storybookjs/storybook | 200 | REFERENCE row `storybook`. Not a new skill |
| chromatic.com | 200 | REFERENCE row `chromatic`. Paid. Not connected |

## Refused as installed systems

| Source | Check | Why |
| --- | --- | --- |
| thedotmack/claude-mem | 200 | The vault is the only long-term memory |
| ECC, GSD, BMAD, gstack as conductors | above | One planner |
| Headroom the product | 200 | Capsules already bound the worker context |
| Career Ops the scraper | 200 | orchestra-career does not scrape or apply |
| UI UX Pro Max as a global pack | 200 | Two design skills already share one DESIGN.md |
| Watermelon MCP, Sentry MCP | not added | No new global MCP. Core stays the brain, Playwright, Context7, and Stitch |

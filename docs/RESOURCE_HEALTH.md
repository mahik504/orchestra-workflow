# Resource health — Orchestra 3.3.0

Do not call a service HEALTHY without host confirmation in this environment. Unauthorized or unpurchased is not active.

Verification date: **2026-09-06**.

| Resource | Representation | Current tier | Host | Health | Auth | Free capability | Paid capability | Verification | Fallback |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| shadcn/ui | skill + optional MCP | FREE | Cursor MCP present locally | LOCAL_CONFIGURED; EXTERNAL_AUTH_VERIFICATION_REQUIRED for other machines | HOST_CONNECT | primitives | none | Cursor mcp.json lists `shadcn` | copy from docs |
| Magic UI OSS | reference | FREE | all | N/A (URLs) | NONE | OSS components | Pro blocks | catalog row | shadcn-ui |
| Magic UI MCP | mcp | FREE_MCP | Cursor / Antigravity / Claude optional | EXTERNAL_AUTH_VERIFICATION_REQUIRED until host lists the server as connected | HOST_CONNECT | `@magicuidesign/mcp` | Pro | example JSON in repo; live merge is host-local | magic-ui-oss |
| React Bits | reference | FREE OSS | all | N/A | NONE | OSS library | Ultimate | catalog row | motion-primitives |
| Aceternity | reference | FREE components | all | N/A | NONE | public categories | All-Access | catalog row | magic-ui-oss |
| Kokonut UI | reference | FREE OSS | all | N/A | NONE | named controls | none claimed | catalog row | react-bits |
| Motion Primitives | reference | FREE OSS | all | N/A | NONE | micro-interactions | none claimed | catalog row | gsap-core |
| Cult UI | reference | FREE OSS | all | N/A | NONE | animated blocks | none claimed | catalog row | shadcn-ui |
| HyperUI | reference | FREE | all | N/A | NONE | Tailwind utility UI | none claimed | catalog row | shadcn-ui |
| Tailkit public | reference | FREE RESOURCE | all | N/A | NONE | public pages | paid kit | catalog row | hyperui |
| Tailkit MCP | mcp | PAID / AUTH_REQUIRED | none connected | NOT CONNECTED | AUTH_REQUIRED | none on MCP | licensed MCP | no fake free MCP | tailkit public |
| 21st.dev | mcp | OPTIONAL_FREE_MCP | none connected here | EXTERNAL_AUTH_VERIFICATION_REQUIRED | AUTH_REQUIRED until login | search; 2 installs/day | Builder / AI credits | docs.21st.dev/mcp | shadcn-ui |
| Refero skill | skill | FREE method | all hosts via sync | skill files present | NONE | methodology + public pages | Pro MCP | skills/refero-design | refero-public |
| Refero public | research_source | FREE | all | N/A | NONE | public styles / pages | Pro | catalog row | godly-design |
| Refero MCP | mcp | AUTH_REQUIRED / PRO | none connected here | NOT CONNECTED | AUTH_REQUIRED | none | Pro search | do not pretend free | refero-public |
| GetLayers | mcp | PURCHASE_PENDING | none | BLOCKED_UNTIL_PURCHASE | PURCHASE_PENDING | none | Full Stack MCP | no install this pass | named OSS 3D stack |
| Godly / Land-book / Lapa / SaaSFrame | research_source | FREE URL | all | N/A | NONE | public galleries | none claimed | no MCP claimed | each other |
| three.js / R3F / drei / postprocessing / scroll-rig | dependency | FREE OSS | project-scoped | N/A | NONE | OSS | none | catalog + oss-index | skip 3D |
| GSAP | skill + dependency | FREE OSS/package | project-scoped | N/A | NONE | core timelines | club plugins if any | gsap-core skill | motion-primitives |
| Scrapling | cli | FREE OSS | on demand | N/A | NONE | public extract | none | catalog row | web_fetch |
| agent-browser | cli | OPTIONAL | on demand | N/A | NONE | foreign-site interact | none | never --tools all | Playwright |
| Playwright MCP | mcp | CORE | Cursor + Antigravity listed | LOCAL_CONFIGURED; EXTERNAL_AUTH_VERIFICATION_REQUIRED as HEALTHY | HOST_CONNECT | our-app QA | none | host mcp lists playwright | project Playwright |
| Context7 | mcp | CORE | Cursor + Antigravity listed | LOCAL_CONFIGURED | HOST_CONNECT | docs | none | host mcp lists context7 | web docs |
| Stitch | mcp | CORE (keyed) | Cursor + Antigravity | LOCAL_CONFIGURED; EXTERNAL_AUTH_VERIFICATION_REQUIRED | AUTH_REQUIRED | 2D screens | quota | key stays in host, never git | DESIGN.md only |
| Serena | mcp | OPTIONAL | not always-on | NOT CONNECTED here | HOST_CONNECT | repo intelligence | none | large-repo jobs only | grep |
| screenshot-to-code | reference | OPTIONAL | named reconstruction | N/A | NONE | reverse reference | none | catalog row | reverse-engineering route |
| vault filesystem | mcp | CORE | Cursor + Antigravity | LOCAL_CONFIGURED on this machine | NONE | private workspace files | none | path is host-local; never commit | none |
| Mobbin | mcp | OPTIONAL | Antigravity listed | EXTERNAL_AUTH_VERIFICATION_REQUIRED | AUTH_REQUIRED | production screens | plan limits | Connect in host | Refero public |
| Firebase MCP | mcp | OPTIONAL | Antigravity only | EXTERNAL_AUTH_VERIFICATION_REQUIRED | AUTH_REQUIRED | cloud data when signed in | none extra | do not install on Cursor to match | skip |
| Supabase MCP | mcp | AUTH_REQUIRED | Antigravity listed | EXTERNAL_AUTH_VERIFICATION_REQUIRED | AUTH_REQUIRED | none without auth | project | Unauthorized ≠ active | postgres skill |
| jcode / OmniRoute | executor | OPTIONAL | jcode | host-local | local provider | terminal executor | none | not an MCP; not claimed HEALTHY for clones | Claude Code |
| devops-skills dump | quarantined_reference | REFERENCE | none | never loaded | NONE | public markdown map | none | do not install 88 | ship-safe |
| awesome prompt repos | prompt_library | REFERENCE | all | N/A | NONE | targeted search | none | do not dump into context | Orchestra protocols |

## Local host smoke (this operator machine, 2026-09-06)

Inspected **server names only**. Secrets were not copied into git.

**Cursor** MCP server names present after this pass: vault filesystem, `shadcn`, `stitch`, `playwright`, `context7`, `magicuidesign-mcp`. Refero / 21st / GetLayers / Tailkit MCP: **not** present. Magic UI is OPTIONAL (host-connect), not claimed HEALTHY until you confirm it in the Cursor MCP panel.

**Antigravity** MCP server names present: `StitchMCP`, `context7`, vault filesystem, `playwright`, `mobbin`, `firebase-mcp-server`, `supabase`, `chrome-devtools-mcp`, `magicuidesign-mcp`. Doctor on this machine: StitchMCP/context7/playwright HEALTHY; supabase AUTH_REQUIRED; magicuidesign-mcp OPTIONAL. Mobbin/Firebase remain EXTERNAL_AUTH_VERIFICATION_REQUIRED. GetLayers is not present.

**Claude Code / jcode MCP:** not invented. jcode has no Orchestra-claimed MCP layer.

GetLayers: **PURCHASE_PENDING**.

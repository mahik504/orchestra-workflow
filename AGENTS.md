# AGENTS — Orchestra 3.3.6

You are running **Orchestra 3.3.6**. One contract, one conductor, open skills. Not a second toolbox.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

Host rules are adapters. They do not invent a second plan. If two hosts could plan, the chat the human is talking to is the conductor.

## Loop

1. **Forge.** `prompt-forge` writes a short brief: intent, kind, route, acceptance, refs, constraints, out of scope. Then work. Ask only when a product decision is truly ambiguous. If it names a repo or a file, open the repo before trusting any brief.
2. **Classify** one route from `registries/design-resource-graph.json`. Load that route, not every skill. The **graph** (`registries/design-resource-graph.json`) is capability routes. The registry is resource knowledge.
3. **Kind** sets the design bar. `skip the lab` still bypasses it.

| kind | Bar | Lab |
| --- | --- | --- |
| hiring-cv | FLAGSHIP | 23 cards, or 3 translations for a named reference, then contract, then stills |
| personal or client | PREMIUM | One human message approves the contract and the stills |
| college | STANDARD | No lab unless the human says PREMIUM |
| backend, fix, glue | STANDARD | No lab |

<!-- routing:start -->
## Routing

Print `Load:` and `Skip:` before writing. A skill the human names wins. A missing tool stops. It does not install itself.

| Job | Load | Skip |
| --- | --- | --- |
| Award-style site | impeccable, taste-design, visual-golden-still, gsap-core, gsap-scrolltrigger, gsap-timeline | copying a reference site's files; a second motion stack |
| Site to echo | visual-forensics, impeccable, taste-design | copying assets, copy, or a logo |
| Scroll | gsap-core, gsap-scrolltrigger, gsap-timeline | gsap-react unless the motion lives in React components |
| One 3D scene | r3f-shaders | 3D as decoration; a second motion engine |
| SaaS or API | backend-architecture, test-driven-development | the design lab; supabase and postgres-best-practices unless that product uses them |
| Bug | systematic-debugging | the design lab |
| Paper | research-paper, humanizer | invented citations |
| Deck | orchestra-docs | host-provided pptx unless that host already has it |
| README and post | orchestra-ship, humanizer | posting for the human |
| Binary or APK | not ready | rea until Java is installed and the target is yours, a client's, a CTF, or open source |
| Hackathon or college | prompt-forge | the 23-card survey when a reference is named; the design lab for college unless the human says PREMIUM |
| Docs on disk | host-provided | docx, pptx, pdf, and xlsx stay on the host. They are not in skills/ |
<!-- routing:end -->

4. **Implement** only after the lab is cleared for visual work. A DESIGN.md is not visual evidence.
5. **Verify** on the running app. `verification-before-completion` plus the ledger. Then `fresh-context-review`.
6. **Remember** what actually worked. Do not invent usage.

## Engineering

| Job | Skill |
| --- | --- |
| Bug | `systematic-debugging` |
| Logic | `test-driven-development` |
| Before DONE | `verification-before-completion` |
| Diff | `ponytail` on engineering only, never on a Design Lab surface |
| Parallel hosts | `using-git-worktrees` and `protocols/PARALLEL_WORK_PROTOCOL.md` |

## Research and reverse engineering

Agent-Reach reads the operator's own profiles and public pages. Read-only. Cookies stay local. No bulk scraping of other people. No post, like, DM, or apply. No ATS scrapers.

`rea` and APKLab run only on a target the task names as own, client, CTF, or open source.

## Evidence

You may write **DONE / FIXED / VERIFIED / PASSED / SHIPPED** only when observed evidence is in the same message. Intention is not evidence.

## Resources

Available is not loaded. A registry row is not proof the tool ran. New third-party skills pass `skill-security-check` before promotion. Never `skills add --all`.

## Hosts

The same skills and the same core MCPs. Exceptions live in `registries/host-stack.json`. `kit/orchestra-doctor.ps1` checks drift.

| Host | Role |
| --- | --- |
| Cursor | Usual conductor |
| Antigravity | Visual QA |
| Claude Code | Terminal and review. The hook reminds Forge. It does not rewrite the prompt |
| Paseo | Launches workers. Not a conductor |
| Codex, jcode | Same skills |

## Overrides

- **skip orchestra** stands this down for the session.
- **skip the lab** bypasses the design gate for one task.
- `ORCHESTRA_CONTRACT` can pin an older contract. Rollback notes are in `kit/ROLLBACK.md`.

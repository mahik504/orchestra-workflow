# AGENTS — Orchestra 3.3.5

You are running **Orchestra 3.3.5**. One contract, one conductor, open skills. Not a second toolbox.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

Host rules are adapters. They do not invent a second plan. If two hosts could plan, the chat the human is talking to is the conductor.

## Loop

1. **Forge.** `prompt-forge` writes a short brief: intent, kind, route, acceptance, refs, constraints, out of scope. Then work. Ask only when a product decision is truly ambiguous.
2. **Classify** one route from `registries/design-resource-graph.json`. Load that route, not every skill.
3. **Kind** sets the design bar. `skip the lab` still bypasses it.

| kind | Bar | Lab |
| --- | --- | --- |
| hiring-cv | FLAGSHIP | 23 cards, or 3 translations for a named reference, then contract, then stills |
| personal or client | PREMIUM | One human message approves the contract and the stills |
| college | STANDARD | No lab unless the human says PREMIUM |
| backend, fix, glue | STANDARD | No lab |

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

Write DONE, FIXED, VERIFIED, PASSED, or SHIPPED only with observed output in the same message. Intention is not evidence.

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

# Architecture findings

Observed at `60cf9e7`.

## What the repo is

Orchestra is a control plane. `skills/` holds capabilities. `registries/` holds resource knowledge. `runtime/` is the Go engine and the `orchestra` CLI. `kit/` installs that stack onto Cursor, Claude Code, Antigravity, Codex, jcode, and the agents host. The vault is not in this repo.

The eight stages live in `runtime/internal/engine/`. A dry run that needs Playwright records the human gate and still prints the plan. Execute still stops. That behavior is covered by `TestPlan_DryRunRecordsMissingPlaywright` and `TestPhase3E_FlowE_BrowserUnavailableHumanGate`.

## Boundaries that are already clear

- One conductor. Host files (`AGENTS.md`, `CLAUDE.md`, `.cursorrules`, `kit/antigravity/GEMINI.md`) are adapters. The routing block between `<!-- routing:start -->` and `<!-- routing:end -->` is the same text in each.
- One motion owner for scroll (GSAP). One 3D scene skill (`r3f-shaders`). `creative-3d-web` and `fullstack-mastery` are host-local and excluded from the canonical count.
- Receipt formatting prints a line only when that line has something to say.
- The commit-msg hook strips `Co-authored-by: Cursor <cursoragent@cursor.com>`. Repo-local `core.hooksPath` is `kit/git-hooks`.

## Boundaries that are still informal

- Skill ownership is a directory name. There is no `registries/skills.json` with license, commit, or host.
- Run lifecycle (`discovered` through `successful`) and catalog admission are the same vocabulary in comments. A row can look verified because `last_verified` is set.
- Model choice is not a registry. Hosts pick models. Nothing in the repo names a model class.
- `orchestra plan` builds a manifest. There is no `orchestra compile` that emits the same shape as JSON for a fixture.

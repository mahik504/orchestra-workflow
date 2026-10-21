# Orchestra 3.3.2 — implementation plan

**3A + 3A.1 are docs only. 3B (parity) applied. 3C + 3C.1 runtime locked. 3D + 3D.1 visual QA locked. 3E + 3E.1 capability gateway locked. 3F + 3F.1 Brain reconciliation is locked. 3G is a measurement pilot. 3G.1 is a real-host pilot on T1 and T3. 3G.2 through 3G.2.4 corrected that pilot. The 2026-10-15 audit is the status tracker.** Do not implement 3H–3I until the operator says so.

Status, 2026-10-22: the operator is publishing this tree. `VERSION` stays `3.3.1`. The git tag for this publication is a prerelease. Stable `v3.3.2` is still 3I.

Source snapshot: `orchestra-workflow` HEAD `a9a5051`. Public clone is 3.3.1. Runtime still reports `VERSION` 3.3.1. Live hosts match repo conductor text after a surgical one-skill copy (not full `sync-ides`).

Do not: bump `VERSION`, tag stable `v3.3.2`, run full `kit/sync-ides.ps1` (44-skill copy), edit host MCP config, or install OpenDesign into Orchestra. A prerelease tag of this tree is allowed. 3G measures T1 and T3. It does not start 3H.

---

## 3A — research and specification

Seven files in this folder. Pointer from `docs/3.3.1/MIGRATION.md` (not live). Done.

## 3A.1 — spec hardening (this pass)

Docs only. Same seven files. Canonical routing lives in [ORCHESTRA_3_3_2_ARCHITECTURE.md](ORCHESTRA_3_3_2_ARCHITECTURE.md).

Hardened: adaptive independent judge; execution budget; two-layer visual QA; lightweight 3–5 directions; baseline vs 3.3.1; hygiene KEEP/MIGRATE/DELETE/ARCHIVE (classify only); complexity six questions; token ladder; retry-must-change.

No runtime, kit, Brain, MCP, host, Portfolio, OpenDesign, VERSION, or git commit.

---

## 3B — parity (docs / conductor / allowlist)

**Parity only.** Do not sneak execution budget, two-layer QA, or mission files into 3B.

Applied this phase on the operator PC (repo remains `VERSION` 3.3.1):

- Conductor skills on Cursor + Gemini hosts = repo 3.3.1; live allowlist **44**.
- Four job-loaded skills now live: `visual-forensics`, `visual-golden-still`, `fresh-context-review`, `final-showcase`.
- `kit/antigravity/install.ps1` reads `registries/host-stack.json` (same allowlist as `kit/sync-ides.ps1`). No private 30-name list.
- `.cursorrules`, README, CLAUDE.md, adapters, Design Lab diagram: named visual refs = **3 translations**, not skip; stills lock, not stack card.
- Brain overlay “3.2 contract” pointers were left for 3F. 3F updated `routes.md` and `START HERE.md`.
- Do not bump Orchestra `VERSION` to 3.3.2 yet (that is 3I).

---

## 3C — mission / capsule / failure / recovery / budget

**Implemented** (runtime still reports `VERSION` 3.3.1; no tag). Package: `runtime/internal/mission`. Engine bind: `runtime/internal/engine/mission_bind.go`. Durable files: `.orchestra/runs/<mission-id>/mission.json`, `state.json`, and `progress.txt` on FULL only.

- Mission file schema + writer/reader (includes intensity + **execution budget**).
- Worker capsule packager (retrieved → scoped → packaged → executed → returned) + measurable `ScopeReport`.
- Failure classes + derived `retry_budget` (hard cap 4, not a universal N) + **retry must change**.
- LIGHT vs FULL intensity (not a new quality bar).
- FULL **long path** only: progress file + stuck detection. LIGHT does not write `progress.txt`.
- `ErrHumanGateRequired` when-list specified in `AGENTS.md` / conductor skill. Sentinel unchanged.
- `visual_review` and `independent_judge` are **policy data** on the budget. 3C does not run Layer A/B or an independent judge.

`retry_budget` N is derived per run. Hard ceiling is 4 so retries cannot run away. Identical signatures stop before the ceiling.

### 3C.1 — hardening (this pass)

Narrow verification + KEEP fixes. Not a redesign.

- Live conductor skill: Cursor, Gemini, Agents, Claude, jcode hashed equal to `skills/orchestra-conductor/SKILL.md`. Kit Antigravity copy already matched. Did **not** run full `sync-ides`.
- Iterate telemetry no longer leaves blocking `pending_human` on a healed SUCCESS run. `ErrHumanGateRequired` / `ErrMaxIterationsExceeded` identities unchanged.
- Giant constraint/evidence dumps rejected the same way as giant objectives. User objective text is not a structural provider leak.
- Human-needed when-list messages (including material tool/host gap) do not auto-retry.
- `visual_review` / `independent_judge` remain inert policy fields. Stage 8 still owns visual self-heal.

---

## 3D — design intelligence + two-layer visual QA

**Implemented** (runtime still reports `VERSION` 3.3.1; no tag). Package: `runtime/internal/visual`. Engine bind: `PlaywrightVerifier` is the canonical Layer A owner; `visual_review` / `independent_judge` are consumed.

- Design class `ordinary | product_ui | flagship` is **not** a quality-bar proxy. Flagship comes from task language (portfolio, brand, art-directed, 3D showcase), not PREMIUM/EXPERIMENTAL.
- Flagship-only **3–5 lightweight** directions → human/waiver pick → **one** sealed contract. CRUD/API/debug/normal UI: zero extra directions. 23-card survey and 3 named-ref translations unchanged.
- Anti-slop critic: actionable findings. Gradients / rounded cards / dark / common fonts are not automatic fails. Inter/`#000` greps retired as DoD.
- **Layer A owner:** harness `PlaywrightVerifier` + `runtime/tools/layer-a.mjs`. Playwright MCP and Cursor browser are host adapters, not pass/fail authorities. `MOCK_SCREENSHOT_*` cannot pass.
- Default viewports: desktop + mobile. Tablet only when policy/project asks.
- **Layer B:** perception on Layer A PNG captures when `visual_review` is `browser` or `deep`. Cannot pass on source/DOM/assertion.
- Independent judge: `DecideJudge` — off for ordinary; conditional for flagship; required only for security/release or explicit request. `RoleJudge` capsule, no write tools.
- Visual failures become structured findings and reuse Stage 8. Retry needs a concrete change. Identical max-iteration stop preserved (`ErrMaxIterationsExceeded`). Golden stills remain pre-implementation lock 2.

Limitation: 3D’s live chromium path was fail-closed (wrong CLI probe, no `runtime/package.json`, Windows spawn). **3D.1** repairs that; see below.

### 3D.1 — real browser + Layer B policy (hardening)

**Implemented** (still no VERSION bump). Project-local `runtime/package.json` dependency `playwright` (library, not `@playwright/test`). Chromium only. `PlaywrightAvailable` succeeds when `runtime/node_modules/playwright` imports **and** Chromium exists (`LookPath` for `node.exe` on Windows; `cmd.Dir` = `runtime/`).

- Flows C/D/E use `PlaywrightRunner`, not `ScriptedRunner`. `ScriptedRunner` remains a unit double.
- Layer B: `structural` (default) vs optional `ScreenshotEvaluator` when `visual_review` is `browser|deep` and class is product_ui/flagship. No model name in core. Ordinary work: zero image eval.
- 1jehuang/jcode: source-only eval in [ORCHESTRA_3_3_2_JCODE_EVAL.md](ORCHESTRA_3_3_2_JCODE_EVAL.md). Not installed. Not Orchestra’s OmniRoute host named `jcode`.

---

## 3E — capability gateway

- `runtime/internal/capability.Select` after mission refine. Minimum sufficient set that is actually consumed (capsule, research skip, design bind, Layer A/B, judge). Not a second conductor. Not a host-execution router. Not a model registry.
- Provider-neutral. No Claude-primary / secondary ranking. Opaque discovered-model labels on `Input` are ignored; selection matches capability class only.
- Catalog filter only: router/host candidate ids in, caps skills ≤ 8 / resources ≤ 8 / mcp tools ≤ 4. Do not rediscover the 44-skill allowlist or the resource catalog.
- Flow F is selection-only host label. `ExecutePlan` stays a stub. Implement handoff literals stay `antigravity` → `cursor`.
- OpenDesign stays retired (banned MCP name, provider-leak token). Do not add it to `~\.cursor\mcp.json`.
- 1jehuang/jcode stays out of Go. Stance remains KEEP Orchestra / ADAPT later (not this phase) / DEFER swarms-embeddings-Firefox / REJECT a second conductor. See [ORCHESTRA_3_3_2_CAPABILITY_GATEWAY.md](ORCHESTRA_3_3_2_CAPABILITY_GATEWAY.md).

---

## 3F — Brain / overlay reconciliation

Implemented. Policy: [ORCHESTRA_3_3_2_BRAIN.md](ORCHESTRA_3_3_2_BRAIN.md).

- Brain remembers. Runtime executes. Registry verifies. Human decides ambiguous or high-impact conflicts.
- Historical aggregate `2648` with an empty payload is `legacy_aggregate_without_payload`. Current evidence is `len(evaluations)`. No invented rows. No zeroing.
- `classifyMCP` is probe-based. Unprobed is `OPTIONAL`. No `mcp.json` edit. No new MCP manager.
- Missing `last_verified` stays unverified. No catalog date stamps. No GitHub catalog row.
- shadcn MCP stays OPTIONAL. OpenDesign stays retired. 1jehuang/jcode stays deferred. Grok 4.7 stays an opaque label.

3F.1 hardened that boundary. Health does not define a requirement. An unprobed tool cannot satisfy a required capability. Brain failure does not stop a run. Retrieval does not upgrade provenance. See [ORCHESTRA_3_3_2_PHASE_3F_1_HARDENING_REPORT.md](ORCHESTRA_3_3_2_PHASE_3F_1_HARDENING_REPORT.md).

---

## 3G — measurement pilot

The old hygiene list is already covered by 3C–3E tests. This phase adds a comparison harness. It answers whether the benchmark can measure Orchestra on T1 and T3. It does not answer whether Orchestra improves output.

Plan: [ORCHESTRA_3_3_2_PHASE_3G_PLAN.md](ORCHESTRA_3_3_2_PHASE_3G_PLAN.md). Result: [ORCHESTRA_3_3_2_PHASE_3G_REPORT.md](ORCHESTRA_3_3_2_PHASE_3G_REPORT.md).

- T1 and T3 fixtures are enabled. T2 and T4–T10 stay `enabled: false`.
- A trace is `trace_only`. The paired host cell in 3G is `unavailable`. Claim stays `none`.
- No skill-pack installs. No extra MCP/plugins. No Chromium study in this pass.

---

## 3G.1 — real host pilot

Extends the 3G harness. It does not add a host manager, model registry, or conductor change.

Plan: [ORCHESTRA_3_3_2_PHASE_3G_1_PLAN.md](ORCHESTRA_3_3_2_PHASE_3G_1_PLAN.md). Result: [ORCHESTRA_3_3_2_PHASE_3G_1_REPORT.md](ORCHESTRA_3_3_2_PHASE_3G_1_REPORT.md).

- One T1 pair and one T3 pair, or an explicit `unavailable` cell.
- `claim` stays `none`.
- T2 stays disabled.

---

## 3H — real benchmarks

Run A–I from [ORCHESTRA_3_3_2_ACCEPTANCE.md](ORCHESTRA_3_3_2_ACCEPTANCE.md). Compare to **representative 3.3.1** on quality, reliability, speed, token/context cost, human intervention, failures, retries, completion, implementation complexity, operational overhead.

Architecture fails if those products are not materially better **or** if they are slightly prettier at unjustified cost (the 10% / 4× anti-pattern).

Class C stays parked; only the **path** is scored.

Do not use these products as an excuse to rewrite Orchestra mid-benchmark.

---

## 3I — tag / release / hygiene

Only after the operator says so.

- Every leftover 3.3.1/3.2/3.1 component classified KEEP / MIGRATE / DELETE / ARCHIVE (seed list in ACCEPTANCE). Evidence for DELETE. No misleading active docs. No dead runtime modules.
- Then `VERSION` → 3.3.2
- git tag (from **MAHI’s terminal**, GitHub authorship MAHI)
- Public orchestra-workflow gets real template commits, not empty heartbeats

No 3.3.2 tag in 3A–3H. No deletion in 3A.1.

---

## Remaining ambiguities (do not invent in later phases)

1. **Which Playwright surface is Layer A’s owner** — **closed in 3D.** Canonical owner is the engine harness (`PlaywrightVerifier` + `runtime/tools/layer-a.mjs`). Playwright MCP and Cursor browser remain operator/host adapters and must not green-wash a run.
2. **Who runs Layer B when `independent_judge` is off** — **closed in 3D.** Same-runtime critic on Layer A captures. Independent judge uses `Compile(..., RoleJudge)` only when `DecideJudge` activates.
3. **Integer `retry_budget` defaults** — derive per run; no global table until 3H shows a need.
4. Overlay `2648` vs empty `evaluations` — INVESTIGATE in 3F.
5. shadcn MCP broken-at-freeze — INVESTIGATE in 3F.
6. Duplicate motion-surface survivor — 3F.

---

## 3A.1 safety / parity audit

| Check | Result required |
| --- | --- |
| Files written | Only `docs/3.3.2/*` (MIGRATION pointer already exists; leave it) |
| `VERSION` | Unchanged (3.3.1) |
| Runtime / engine / kit scripts | Unchanged |
| Brain / Preferences / MCP / hosts / Portfolio / OpenDesign / Gate 0 | Unchanged |
| One conductor | Preserved |
| LIGHT vs STANDARD/PREMIUM/EXPERIMENTAL | Intensity vs bars; no fourth bar |
| Independent judge | Conditional; not a synonym for FULL |
| 3B | **Not started** |
| Commit | None unless operator asks |

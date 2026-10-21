# Orchestra 3.3.2 — delta from 3.3.1

Status, 2026-10-22: 3B through 3G.2.4 are in this tree. `VERSION` is still `3.3.1`. 3H and 3I are not done. The tables below are the original delta, not a second status tracker. Current status is [ORCHESTRA_3_3_2_AUDIT.md](ORCHESTRA_3_3_2_AUDIT.md).

---

## KEEP

- One conductor.
- Control plane / capabilities / executors / Brain / registry.
- STANDARD / PREMIUM / EXPERIMENTAL quality bars.
- Two-lock Design Lab; golden stills as pre-implementation lock.
- Real-app Playwright **policy** (becomes Layer A).
- Keep / Rework / Kill / Park.
- Discover broadly / activate selectively.
- Job-loaded four visual/review skills (once hosts match the repo).
- Named visual refs → 3 translations (engine + AGENTS).
- Pasted DESIGN.md skips survey only.
- 8-stage engine as the spine.
- `ErrHumanGateRequired`.
- OpenDesign banned as second studio; OD×Antigravity retired.
- 3.3.1 rejections: GSD, ECC dumps, Karpathy packs, always-on Magic UI.

---

## REMOVE

- Mock screenshot files counting as visual QA pass (`MOCK_SCREENSHOT_*` as success).
- Teaching the 3.3.0 named-website **skip** in `.cursorrules`, README “How it works”, and stale adapters.
- “Stack card” as the visual gate in leftover architecture prose.
- Inflated overlay counters with empty payloads.
- Treating duplicate Playwright / Context7 / motion surfaces as three independent owners (collapse; do not keep three truths).
- Treating FULL as “always spawn an independent judge.”
- Treating 3–5 design directions as 3–5 implementations.
- Inter / `#000` greps as Definition of Done.

Nothing in 3A / 3A.1 deletes code. These are specified removals for 3B–3F / 3I.

---

## CHANGE

| From 3.3.1 | To 3.3.2 |
| --- | --- |
| Hosts 3.3.0 / 40 skills | Hosts match repo 3.3.1 allowlist (44) — **3B** |
| Antigravity install hardcoded 30 | Same allowlist as `sync-ides.ps1` — **3B** |
| Independent review sometimes implied, sometimes skipped, or implied for all FULL | Evaluator **policy**; independent judge **conditional** on triggers (not a synonym for FULL) |
| Workers get ad-hoc context | Worker **capsule** + token ladder |
| Retry / loop underspecified | Failure **class + budget + meaningful change** |
| Anti-slop ≈ Inter/#000 grep | Absence-of-decision questions + 2026 cliché calibration |
| Visual QA = screenshots + token greps | **Layer A** Playwright objective pass/fail; **Layer B** perception on A’s screenshots when budgeted; stills stay lock 2 |
| Optional GitHub MCP with no catalog row | Optional with an honest row, or absent |
| Brain “3.2 contract” wording | Match live Orchestra version in overlay docs — **3B**, docs only |
| Capability → model → runtime (underspecified) | Recommend model, reasoning, host, MCP, skill, resource, OpenDesign, Stitch, specialist; human if material choice is missing |
| “Materially better” = any one axis | Baseline vs 3.3.1; reject unjustified 4× cost for 10% prettier |

---

## ADD

- Mission file (run-scoped, compact), including intensity + execution budget.
- Intensity LIGHT vs FULL (not a fourth quality bar).
- Execution budget: `reasoning`, `context`, `research`, `visual_review`, `independent_judge` (`off\|conditional\|required`), `retry_budget` `0..N` (N derived, not a global preset).
- Adaptive evaluator policy (QUALITY × EVIDENCE / COST).
- Two-layer visual QA (Layer A deterministic / Layer B perception).
- Design five stages: research / exploration / contract / implementation / visual verification.
- Failure taxonomy + retry-must-change.
- FULL-only progress file + stuck detection on the **long** path.
- Memory classes (fact/decision/preference/observation/constraint/stale/proposed/verified).
- Specified human-needed situations (including material tool/model ambiguity).
- 3–5 **lightweight** directions **only** for flagship creative; one selected implementation.
- Optional OpenDesign visual worker (later, 3E) — not default, not conductor.
- Release hygiene classes KEEP / MIGRATE / DELETE / ARCHIVE (3I; classify in spec now).
- Baseline comparison axes for 3H (quality, reliability, speed, tokens, human intervention, failures, retries, completion, complexity, overhead).

---

## DEFER

- VERSION bump and `v3.3.2` git tag (3I, operator-gated).
- HarnessXP-style multi-repo contract sync.
- Memory-plugin MCP / measured memory harness.
- Spexcode spec-node product.
- Click-to-source interactive preview IDE.
- Visionaire-class geometry assertions (after Layer A is real).
- Bulk catalog verification of all 126 unverified rows (honesty first, then verify in batches).
- Wiring OpenDesign into host-stack.
- Choosing which Playwright surface is the canonical owner (Cursor-browser vs MCP vs `visual_qa.go`) — **3D**, recorded as ambiguity until then.
- Default integer `retry_budget` tables — derive per run; do not invent globals in this spec.

---

## REPLACE

| Replace | With |
| --- | --- |
| Always-on three-agent theatre | Adaptive roles; independent judge when **triggers** fire |
| FULL = independent judge | FULL = heavy path; judge from evaluator policy |
| Tenet DAG / SQLite job OS | Compact `mission → run → stage → worker → evidence → evaluation → handoff` |
| Spec Kit slash-command zoo | Compact spec/plan/tasks artifacts on STANDARD+ FULL when useful |
| Pixel-clone / SSIM-to-reference as DoD | Human stills + Layer A + Layer B |
| Playwright as art director | Layer A objective; Layer B perception |
| Banned-pattern dump as law | Critic questions; brief wins |
| 3–5 full implementations | 3–5 lightweight directions, one contract, one build |
| `kit/antigravity/install.ps1` private 30-name list | Canonical allowlist |

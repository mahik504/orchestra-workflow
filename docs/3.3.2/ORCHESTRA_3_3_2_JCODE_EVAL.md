# Orchestra 3.3.2 — 1jehuang/jcode evaluation

Source-only. **Not installed** into Orchestra, MCP, host skills, or this PC for 3D.1.

**3E did not change that.** The capability gateway does not import, string-match, exec, or reinstall 1jehuang/jcode. ADAPT later still means a later phase.

This is **not** Orchestra’s existing fifth host named `jcode` (OmniRoute / terminal adapter on the 3.3 host list). That host stays an executor. [1jehuang/jcode](https://github.com/1jehuang/jcode) is a different product: a Rust TUI/CLI coding harness at [jcode.sh](https://jcode.sh/).

3E is the earliest phase that may consider it as an **optional worker** through the capability gateway. 3D.1 does not add parallel workers to the conductor.

---

## VERIFIED FROM SOURCE

Repository inspected at commit `2a4edaa02057ac994a601311c4f03ed450e1b3c9` (GitHub default at eval time). Site: https://jcode.sh/ (fetched 2026-09-21). README, `docs/WINDOWS.md`, and `crates/jcode-app-core/src/tool/browser.rs`.

| Claim | Source |
| --- | --- |
| Rust TUI/CLI coding agent, MIT | README; `Cargo.toml` / `LICENSE` in repo root |
| Platforms include Windows | README badge “Linux \| macOS \| Windows”; `docs/WINDOWS.md` (“Windows 11 x64 supported and manually verified”); PowerShell installer `irm https://jcode.sh/install.ps1 \| iex` |
| `jcode run` headless / non-interactive | README Quick Start: `jcode run "say hello"`; site Auto-poke: “the same mechanism drives headless runs, so `jcode run` keeps iterating” |
| `jcode serve` / `jcode connect` | README Quick Start |
| Swarms | README Swarm section; `docs/SWARM_ARCHITECTURE.md` linked |
| Embedding memory | README Memory section: turn embeddings, cosine graph, optional memory sideagent |
| MCP: stdio servers | README: “jcode currently supports stdio (command-based) servers only; HTTP/SSE entries are recognized and skipped” |
| Built-in browser is **Firefox** (Firefox Agent Bridge), not Playwright | README Browser Automation; `browser.rs` defines `static FIREFOX_PROVIDER: FirefoxBridgeProvider` |
| Browser actions include status/setup/open/snapshot/click/type/screenshot | README Browser Automation list |
| Windows install paths | `docs/WINDOWS.md`: `%LOCALAPPDATA%\jcode\bin\jcode.exe`, user data `%USERPROFILE%\.jcode` |

No Orchestra config, MCP, or skill was changed to install this product.

---

## PROJECT CLAIMS (not independently verified here)

These numbers and qualitative claims are **the project’s measurements**, published on jcode.sh / README. 3D.1 did not re-run their benches, RAM traces, or prompt-cache accounting on this PC.

- Extra PSS ~10.4 MB per additional session; “ten jcode sessions cost about 100 MB”
- Time to first frame 14.0 ms and time to first input 48.7 ms vs other CLIs (Linux PTY sample on their machine)
- Prompt-cache / append-only context engineering as the cost lever
- Terminal-Bench 2.1 confidence-stepping pass-rate deltas (92% vs 88% in-time, their trial)
- jcode bench v1 float-print vs Claude Code on Opus 4.8 (their plot)
- “Most RAM efficient harness” marketing line

Treat those as **project claims**. Do not copy them into Orchestra acceptance as facts.

---

## OUR INFERENCE

jcode is a **peer harness**, not a missing Orchestra subsystem. It already has its own conductor-shaped loop: session, swarm coordinator, auto-poke retries, memory graph, browser tool, MCP catalog, and `jcode run` until todos complete.

Orchestra already owns: one conductor, mission + budget, Stage 8 classified retry, two-lock Design Lab, and harness Layer A (Playwright + Chromium in `runtime/`).

A second copy of those mechanisms would fight the 3.3.2 control-plane rule.

Windows and headless `jcode run` are documented in **source**, so 3D.1 did not perform an isolated binary smoke (no temp install, no PATH change, no Orchestra config).

---

## Per-mechanism stance

| Mechanism | Stance | Why |
| --- | --- | --- |
| Orchestra conductor, mission, budget, Stage 8, Design Lab, harness Layer A | **KEEP ORCHESTRA** | Already the control plane and the Layer A owner |
| Headless `jcode run` as a later optional worker | **ADAPT later (3E gateway)** | Might be cheaper/faster than some current executors for bounded implementation jobs — *if* 3E compares cost on real work. Optional, capability-gated, never default |
| Swarms / multi-agent in one repo | **DEFER** | Orchestra is one conductor. Swarm is a second coordination plane |
| Embedding memory graph | **DEFER** | Orchestra Brain is overlay facts/decisions/preferences, not a vector sideagent |
| Firefox Agent Bridge browser | **DEFER** | Duplicate of Layer A. Canonical owner stays Playwright Chromium in `runtime/` |
| Second conductor, second mission/retry/state, mandatory dependency, always-on parallel core | **REJECT** | Breaks 3.3.2 architecture |

---

## Recommendation

**Optional later, via 3E capability gateway. Not now.**

Do not add 1jehuang/jcode to Orchestra core, MCP, or host skills in 3D.1. Do not confuse it with the OmniRoute host named `jcode`. Isolated smoke is unnecessary: Windows + `jcode run` are already in public source.

If 3E later wants a worker probe: temp directory only, no Orchestra config, no second Layer A.

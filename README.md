# Orchestra V3 — Autonomous Agentic Orchestration Architecture

[![Go Runtime CI](https://github.com/mahik504/orchestra-workflow/actions/workflows/ci.yml/badge.svg)](https://github.com/mahik504/orchestra-workflow/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release: V3.1-Production](https://img.shields.io/badge/Release-V3.1--Production-0A2118.svg)](#)

> **Orchestra V3** is a high-performance, contract-driven agent orchestration kernel engineered for professional AI pair programming and multi-agent development. It coordinates specialized coding models, governs cognitive load, enforces strict human approval gates, and guarantees reproducible, production-grade deliverables across **Cursor**, **Antigravity**, and **Claude Code**.

---

## 1. What Orchestra Is

Modern AI coding agents frequently fail at two extremes: they either produce shallow, repetitive templates (*AI slop*) because prompt context is too generic, or they burn excessive context tokens trying to load dozens of competing tools into one giant session.

Orchestra solves this by acting as an **operating system kernel for AI coding**:
- **Normalized Task Contracts:** Deconstructs complex user briefs into validated JSON schemas.
- **Dynamic Capability Composition:** Evaluates tasks and selects only the minimal sufficient skill set.
- **4-Stage Capability Pipeline:** Moves beyond nominal "skill naming" to enforce *Retrieval → Analysis → Application Directives → Adversarial Verification*.
- **Multi-Agent Handoff:** Provides versioned state transfer between Cursor and Antigravity with cryptographic file integrity and out-of-band conflict detection.
- **Strict Storage Governance:** Maintains a clean, permanent memory layer with complete private/public boundary isolation.

---

## 2. The Problem It Solves

| Failure Mode in Typical Agent Workflows | How Orchestra V3 Resolves It |
|---|---|
| **Context Bloat & Token Exhaustion** | Lazy-loads skills on-demand; enforces Ponytail memory retention to purge transient scratch. |
| **Decorative Skill Usage (Naming without Consulting)** | Router synthesizes mandatory `CapabilityExecutionDirective` objects with banned anti-patterns and verification checklists. |
| **Loss of State across Tools** | Structured `state.json` contract with SHA256 checksums, version increments, and automated resume vectors. |
| **Monotonous AI Slop (Generic Templates)** | Enforces bespoke architectural layouts, typography locks (`Fraunces` / `Plus Jakarta Sans` / `JetBrains Mono`), and bans repetitive 3-column card rows. |
| **Secret Leaks & Private Data Bleed** | Two-tier architecture: Personal Brain (`orchestra-brain`) is air-gapped from the public distribution (`orchestra-workflow`). |

---

## 3. Architecture Evolution: V1 → V2 → V3

```
V1: Monolithic Markdown Prompts & Manual Clipboard Handoffs
    └── 35% goal completion, prompt drift, no verification contracts.

V2: Python Scripts & Modular Rule Injection
    └── 72% goal completion, improved classification, but high context burn and tool bloat.

V3: Compiled Go Runtime & Contract-Driven Multi-Agent Kernel
    └── Modular 4-stage pipeline, lazy capability loading, versioned handoffs, and adversarial verification gates.
```

---

## 4. Architectural Overview

```
                      +----------------------------------+
                      |      User Brief / Task PRD       |
                      +-----------------+----------------+
                                        |
                                        v
                      +----------------------------------+
                      |        Orchestra Classifier      |
                      |   (Type, Visual, Security, Gaps)  |
                      +-----------------+----------------+
                                        |
                                        v
                      +----------------------------------+
                      |         Capability Router        |
                      |  (Minimal Sufficient Composition)|
                      +--------+----------------+--------+
                               |                |
             +-----------------+                +-----------------+
             v                                                    v
+---------------------------+                      +---------------------------+
|    Capability Registry    |                      |      Agent Allocation     |
| - superpowers-planning    |                      | - Cursor: Bulk/Refactor   |
| - taste-design            |                      | - Antigravity: Multi-Tool |
| - impeccable / motion     |                      | - Claude: Deep Reasoning  |
| - semgrep / security      |                      +-------------+-------------+
+------------+--------------+                                    |
             |                                                   v
             +----------------->[ EXECUTION MANIFEST ]<----------+
                                        |
                                        v
                      +----------------------------------+
                      |    Dual-Agent Execution & QC     |
                      |   (State v1 -> v2, SHA256 Check) |
                      +-----------------+----------------+
                                        |
                                        v
                      +----------------------------------+
                      |    Adversarial Quality Gates     |
                      | (Playwright E2E, Lighthouse 13)  |
                      +-----------------+----------------+
                                        |
                                        v
                      +----------------------------------+
                      |      Production Deployment       |
                      +----------------------------------+
```

---

## 5. The 4-Stage Capability Pipeline

Orchestra V3 guarantees that registered tools are actively enforced through an automated 4-stage lifecycle:

1. **Retrieval (`cap.LoadDetails()`):** Reads the exact instructions, schemas, and rule sets from local disk.
2. **Analysis:** Parses the tool guidelines to extract actionable rules, banned anti-patterns, and unit-level constraints.
3. **Application Contract:** Compiles a markdown `Execution Manifest` that is prepended to the implementing agent's working prompt.
4. **Adversarial Verification:** Automatically checks code output against the directives (e.g., automated grep audits for banned font stacks, viewport overflow checks, and security audits).

---

## 6. Agent Allocation Matrix (Cursor vs Antigravity)

Orchestra matches task profiles to agent execution engines:

| Task Characteristics | Recommended Agent | Recommended Model & Effort |
|---|---|---|
| **Multi-File Structural Refactor / Heavy Edits** | **Cursor** | Claude 3.5 Sonnet / High Context |
| **Multi-Tool Research / MCP / Browser E2E** | **Antigravity** | Gemini 2.0 Pro / High Effort |
| **Algorithm Optimization / Deep Logic** | **Claude Code** | Claude 3.7 Sonnet (Thinking) |
| **Rapid Prototyping / Local Smoke Checks** | **Antigravity** | Gemini 2.0 Flash / Standard |

---

## 7. Versioned Handoff Protocol

When switching between development environments (e.g., planning in Antigravity → bulk editing in Cursor → verifying in Antigravity), Orchestra writes a versioned `.orchestra/handoff/state.json`:

```json
{
  "session_id": "sess-84920",
  "version": 2,
  "timestamp": "2026-09-03T16:12:00Z",
  "source_agent": "cursor",
  "target_agent": "antigravity",
  "active_tasks": ["task-redesign-v3"],
  "changed_files": [
    { "path": "src/pages/AboutPage.tsx", "sha256": "4f18a2..." },
    { "path": "src/components/layout/Navbar.tsx", "sha256": "9b72e1..." }
  ],
  "completed_steps": ["step-1-tokens", "step-2-routes"],
  "pending_steps": ["step-3-playwright-qa"],
  "failure_recovery": {
    "can_resume": true,
    "resume_from_step": "step-3-playwright-qa"
  }
}
```

If an external tool or human edits a tracked file out-of-band, `DetectConflicts()` flags the SHA256 mismatch before work proceeds, preventing silent overwrite disasters.

---

## 8. Real Benchmark Validation — TTB Agro Redesign

The Orchestra V3 system was tested by driving a full corporate multi-route overhaul of **TTB Agro India Private Limited**, a B2B agricultural trading house.

### Empirical Evidence Register:

| Metric Category | Target Standard | Measured Result (Verified) | Verification Tool |
|---|---|---|---|
| **Route Coverage** | 100% of public views | **12 / 12 Routes Verified** | Playwright E2E |
| **Multi-Viewport Visuals** | Desktop, Tablet, Mobile | **36 / 36 Screenshots Captured** | Playwright Chromium |
| **First Contentful Paint (FCP)** | < 1.0s (Unthrottled) | **804 ms** | Chrome Performance API |
| **First Paint (FP)** | < 500 ms | **392 ms** | Chrome Performance API |
| **Cumulative Layout Shift (CLS)** | 0.00 | **0.00** | Lighthouse 13.4.1 |
| **Total Blocking Time (TBT)** | < 50 ms | **0 ms** | Lighthouse 13.4.1 |
| **Lighthouse Best Practices** | 100 | **100 / 100** | Lighthouse 13.4.1 |
| **Lighthouse SEO** | 100 | **100 / 100** | Lighthouse 13.4.1 |
| **Lighthouse Accessibility** | > 90 | **93 / 100** | Lighthouse 13.4.1 |
| **Lighthouse Performance** | Throttled Mobile 4G | **66 / 100** | Lighthouse 13.4.1 |
| **CSS Gzip Payload** | < 40 kB | **6.97 kB** (35.7 kB uncompressed) | Vite 6 Build Analyzer |
| **JS Gzip Payload** | < 150 kB | **127.9 kB** (430.8 kB uncompressed) | Vite 6 Build Analyzer |
| **Total Page Transfer** | < 2.0 MB | **1.43 MB** (including all assets) | Network Request Interceptor |
| **Production Uptime & Status** | HTTP 200 on all endpoints | **10 / 10 Tests Passed** | Live Production Smoke Suite |

*Note on Evaluation Scoring:* Automated metrics and engineering checks above are empirical. Qualitative design scores (e.g. brand tone and layout rhythm) are marked as `MODEL-EVALUATED / NOT INDEPENDENTLY VALIDATED` in accordance with Orchestra publication standards.

---

## 9. Installation & Quickstart

### Prerequisites:
- **Go 1.22+**
- **Node.js 18+** & **npm**

### 1. Clone the Public Repository
```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow/runtime
```

### 2. Verify Kernel Status
```bash
go test ./internal/... ./cmd/...
```

### 3. Initialize a New Project Profile
```bash
go run ./cmd/orchestra init --name "my-fintech-app"
```

### 4. Classify and Compose a Task
```bash
go run ./cmd/orchestra route --task "Build institutional trading dashboard with high-contrast serif typography and strict input validation"
```

---

## 10. Privacy & Brain Isolation Model

Orchestra maintains strict data privacy:
- **`orchestra-workflow` (Public):** Contains zero credentials, zero private client requirements, and zero personal notes. Distributed under the MIT license.
- **`orchestra-brain` (Private):** Air-gapped personal Obsidian knowledge vault backed by an automated 12-hour encrypted scheduled synchronization task (`OrchestraBrainVaultSync`).

---

## 11. Project Roadmap

- [x] Compiled Go runtime kernel with sub-second execution
- [x] 4-Stage Capability Pipeline (Retrieval → Analysis → Application → Verification)
- [x] Cryptographic state handoff between Cursor and Antigravity
- [x] Production benchmark verification on TTB Agro
- [ ] Direct IDE extension integration for native VS Code / Cursor status bar
- [ ] Distributed multi-machine task DAG runner with WebAssembly sandbox

---

## 12. License

Distributed under the **MIT License**. See `LICENSE` for details.

# Orchestra V3

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Architecture: Capability--Graph](https://img.shields.io/badge/architecture-Capability--Graph-emerald.svg)](#architecture)
[![Agents: Cursor%20%7C%20Antigravity%20%7C%20Claude](https://img.shields.io/badge/agents-Cursor%20%7C%20Antigravity%20%7C%20Claude-purple.svg)](#cross-agent-portability)

**A portable, capability-driven execution engine and resource orchestration layer for AI-assisted software engineering.**

Orchestra V3 replaces the obsolete "load every AI skill into prompt context" anti-pattern with a dynamic, mathematically sound **Capability Graph** and a compiled Go runtime engine. Rather than guessing which dependencies or scripts to install, Orchestra automatically infers, discovers, acquires, and verifies the exact resources required to achieve a production-grade quality bar.

---

## Table of Contents

- [Why Orchestra V3](#why-orchestra-v3)
- [Architecture Overview](#architecture-overview)
- [Canonical Registries & Graph](#canonical-registries--graph)
- [Design-First Go Execution Engine](#design-first-go-execution-engine)
  - [The 8-Stage Pipeline](#the-8-stage-pipeline)
  - [Multi-Source Research](#multi-source-research)
  - [Automated Visual QA & Verification](#automated-visual-qa--verification)
- [Resource Acquisition Adapters](#resource-acquisition-adapters)
  - [Supported Adapters](#supported-adapters)
  - [Provenance Tracking](#provenance-tracking)
- [Cross-Agent Portability & Quarantine](#cross-agent-portability--quarantine)
  - [The 30-Skill Active Core](#the-30-skill-active-core)
  - [Quarantine Boundary Isolation](#quarantine-boundary-isolation)
- [CLI Reference & Commands](#cli-reference--commands)
- [Quick Start](#quick-start)
- [Verification & Testing](#verification--testing)
- [License](#license)

---

## Why Orchestra V3

Generic AI coding agents dump uncurated skills and bloated prompts into context windows, resulting in degraded reasoning, hallucinations, generic purple-gradient UIs, and unvetted global package installations.

Orchestra V3 introduces a strict division of responsibility:
1. **Agents are Executors**: Cursor, Google Antigravity, and Claude Code execute tasks within their respective IDEs and runners.
2. **Orchestra is the Engine**: A lightweight, compiled Go kernel that manages capability resolution, multi-source research, design system synthesis, conditional dependency acquisition, visual verification, and durable memory.

```
+-----------------------------------------------------------------------------+
|                                 USER INTENT                                 |
+-----------------------------------------------------------------------------+
                                       |
                                       v
+-----------------------------------------------------------------------------+
|                           ORCHESTRA GO RUNTIME                              |
|                                                                             |
|  [Classify] -> [Research] -> [Synthesize] -> [Acquire] -> [Implement] ->   |
|                                                                             |
|                      [Visual QA Gate] -> [Remember]                         |
+-----------------------------------------------------------------------------+
        |                               |                              |
        v                               v                              v
+----------------+             +------------------+           +---------------+
|  Cursor Rules  |             | Antigravity Kit  |           |  Claude Code  |
|  (.cursorrules)|             | (kit/antigravity)|           |  (CLAUDE.md)  |
+----------------+             +------------------+           +---------------+
```

---

## Architecture Overview

Orchestra V3 comprises five decoupled subsystems:

1. **Canonical Registries (`registries/`)**: Machine-readable JSON specifications validated by JSON Schema (`registries/schemas/`). These define all verified tools, repositories, packages, MCP servers, and their capability mappings.
2. **Design-First Engine (`runtime/internal/engine`)**: An 8-stage state machine that enforces multi-source research and design synthesis before code generation.
3. **Acquisition Adapters (`runtime/internal/adapters/acquisition`)**: Kernel adapters that verify, sandbox, and conditionally install resources without polluting the global system.
4. **Agent Host Adapters (`runtime/internal/adapters/`)**: Configuration sync guaranteeing identical capability contracts across Cursor, Antigravity, and Claude Code.
5. **Durable Resource Memory (`memory/resource-memory.json`)**: Persistent JSON memory capturing quantitative outcomes (success rates, latencies, failure modes) for every resource utilized.

---

## Canonical Registries & Graph

Orchestra decouples resource declarations from runtime logic:

- **`registries/resources.json`**: Authoritative inventory of vetted resources across categories (`FRONTEND`, `BACKEND`, `SECURITY`, `MOBILE`, `DESIGN`, `TESTING`). Each entry specifies:
  - `id`: Canonical identifier (e.g., `playwright`, `gsap`, `taste-design`).
  - `canonical_url`: Upstream repository or official documentation link.
  - `representation`: Nature of resource (`skill`, `dependency`, `cli`, `mcp`, `reference`).
  - `acquisition_method`: How the resource is retrieved (`npm`, `git`, `cli`, `mcp`, `web_fetch`).
  - `runtime_method`: Execution scope (`project_scoped_install`, `on_demand_exec`, `global_tool`).
- **`registries/design-resource-graph.json`**: Directed capability graph mapping high-level domains (`visual_design`, `qa_testing`, `animation_motion`) and capabilities (`premium-editorial-web`, `visual-regression`, `physics-springs`) to prioritized resource sequences.
- **Validation**: All registries are validated against strict JSON schemas in `registries/schemas/resource-catalog.schema.json` and `registries/schemas/design-resource-graph.schema.json`.

---

## Design-First Go Execution Engine

### The 8-Stage Pipeline

When given a high-ambition task, the Go runtime router (`runtime/internal/engine`) executes an 8-stage sequential pipeline:

```
[1. Discover]      Inspect workspace state and identify existing tooling
       |
[2. Classify]      Analyze prompt complexity, quality bar, and archetype
       |
[3. Research]      Aggregate external curated design and architecture references
       |
[4. Synthesize]    Compile typography, palette tokens, motion curves, and rules
       |
[5. Design System] Generate contract-bound DESIGN.md specification
       |
[6. Implement]     Acquire approved scoped resources and execute code changes
       |
[7. Visual QA]     Run multi-viewport headless verification and layout audit
       |
[8. Iterate]       Self-heal visual regressions up to max iteration threshold
```

### Multi-Source Research

For `frontend_premium` or high-visual tasks, Orchestra queries curated design indexes (Awwwards, Jiro, Cari, DesignMD) to extract:
- Dominant aesthetic philosophies (Swiss Editorial, Brutalist, Glassmorphism, Industrial).
- Exact semantic color tokens (`--color-bg-base`, `--color-surface-elevated`, `--color-accent-primary`).
- Contrast ratios meeting WCAG AAA specifications.

### Automated Visual QA & Verification

The Visual QA stage (`runtime/internal/engine/stage_visual_qa.go`) evaluates output against rigorous design heuristics:
- **Multi-Viewport Audit**: Desktop (`1440x900`), Tablet (`768x1024`), Mobile (`390x844`).
- **Horizontal Overflow Protection**: Programmatic rejection of mobile horizontal scrollbars (`scrollWidth > clientWidth`).
- **Anti-Pattern Detection**: Prohibits uncalibrated pure black (`#000000`), generic purple-to-blue gradient cards, unstyled default fonts, and missing touch targets (`< 44x44px`).
- **Iterative Self-Healing**: Automatically routes layout defects to `StageImplement` and token style defects to `StageDesignSystem`.

---

## Resource Acquisition Adapters

### Supported Adapters

Orchestra acquires dependencies through sandboxed, policy-enforcing adapters in `runtime/internal/adapters/acquisition/`:

| Adapter | Capability | Installation Policy |
|---|---|---|
| **NPM** | React/Vue packages, animation libraries, UI components | Project-scoped only (`npm install --save`). **Global flags (`-g`, `--global`) are programmatically blocked.** |
| **Git** | Verified reference repos and component sources | Clones to ephemeral cache or submodules with commit pin verification. |
| **CLI** | Compilers, linters, and verification binaries | Checks existence in `$PATH`; instructs non-destructive local installation. |
| **MCP** | Model Context Protocol servers | Verified against approved manifest (`Stitch`, `orchestra-brain`, `playwright`). |
| **Web Fetch** | Documentation lookups and design reference assets | Offline fallback caching; strict URI scheme validation (`https://` only). |

### Provenance Tracking

All acquisitions are immutably logged to `.orchestra/provenance.json`:
- Exact timestamp and installing agent.
- Target workspace path.
- Source URL and upstream package registry version.
- Cryptographic SHA-256 integrity hash of acquired assets.

---

## Cross-Agent Portability & Quarantine

### The 30-Skill Active Core

Orchestra maintains a lean, unified active skill set of **30 verified skills** synchronized across all supported AI coding hosts:
- **Cursor**: Configured via `.cursorrules` and `.cursor/skills/`.
- **Google Antigravity**: Configured via `kit/antigravity/` and workspace settings.
- **Claude Code**: Configured via `CLAUDE.md` and `~/.claude/skills/`.

Active skills span:
- **Core Orchestration**: `orchestra-conductor`, `orchestra-vault`, `orchestra-ship`, `orchestra-docs`.
- **Visual Design**: `taste-design`, `emil-design-eng`, `impeccable`, `animate`, `review-animations`.
- **Security & CI**: `ship-safe`, `semgrep-adapter`, `penetration-testing-with-strix`, `ci-security-scanning-with-strix`.
- **Mobile Development**: `expo-router`, `expo-native-ui`, `expo-data-fetching`, `eas-app-stores`.
- **Visual Synthesis (Stitch)**: `stitch-generate-design`, `stitch-react-components`, `stitch-manage-design-system`.

### Quarantine Boundary Isolation

The legacy 1,598-skill uncurated library remains **strictly quarantined** on disk:
- Physical and logical separation prevents token bloat and unvetted code execution.
- The Go runtime kernel enforces hard path boundaries: any attempt to traverse, index, or load files from quarantined paths triggers an immediate `ErrQuarantinedPath` security violation.

---

## CLI Reference & Commands

The compiled Orchestra CLI (`orchestra` or `orchestra.exe`) provides standard lifecycle commands:

```bash
# Display system overview and command reference
orchestra

# Initialize fresh Orchestra workspace (.orchestra/, memory/, registries/)
orchestra init [directory]

# Exhaustive environment, registry, toolchain, and quarantine diagnostic
orchestra doctor

# Classify task prompt and resolve target archetype
orchestra classify --task "Build interactive WebGL solar system showcase"

# Synthesize capability routing and execution manifest (zero side effects)
orchestra plan --task "Create responsive financial analytics dashboard"

# Execute complete 8-stage pipeline with automated acquisition and QA
orchestra run --task "Build luxury timepiece storefront" --auto-approve

# Run standalone multi-viewport visual verification
orchestra verify --strict --workdir ./my-project

# Synchronize host adapter configurations across Cursor, Antigravity, Claude
orchestra sync

# Query or record resource evaluation outcomes in durable memory
orchestra memory --list
orchestra memory --stats
```

---

## Quick Start

### 1. Prerequisites
- **Go**: Version 1.22 or higher
- **Node.js**: Version 18 or higher (with `npm`)
- **Git**: Version 2.30 or higher

### 2. Clone and Build
```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow/runtime

# Build the runtime binary
go build -o orchestra.exe ./cmd/orchestra

# Run system doctor diagnostic
./orchestra.exe doctor
```

### 3. Generate a Plan
```bash
./orchestra.exe plan --task "Develop an editorial coffee roastery showcase with smooth parallax scrolling"
```

---

## Verification & Testing

Orchestra V3 includes an exhaustive, multi-layered test suite:

```bash
# Run all unit and integration tests
cd runtime
go test -v ./...

# Run static analysis and linting
go vet ./...
```

The test suite validates:
- Complete 8-stage pipeline state transitions and early-halt invariants.
- Programmatic rejection of global `-g` / `--global` npm installation attempts.
- Quarantine path breach detection across backslash, forward-slash, percent-encoded, and 8.3 short-name variations.
- Multi-viewport visual QA oscillation and eventual healing logic.
- Concurrent execution determinism across parallel pipelines.

---

## License

MIT License. See [LICENSE](LICENSE) for details.

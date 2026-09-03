# Orchestra V3
*Portable orchestration for AI coding agents and agentic development environments.*

[![Go Runtime CI](https://github.com/mahik504/orchestra-workflow/actions/workflows/ci.yml/badge.svg)](https://github.com/mahik504/orchestra-workflow/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release: 3.0.0](https://img.shields.io/badge/Release-3.0.0-0A2118.svg)](#)

## 1. Identity
Orchestra V3 is a portable, agent-agnostic, resource-intelligent orchestration layer. It sits ABOVE your chosen coding agent (Cursor, Antigravity, Claude Code, etc.) to manage cognitive load, route specialized capabilities, and enforce strict execution quality. 

## 2. Problem
Normal agent workflows fail at scale. If you give an agent zero context, you get generic "AI slop". If you dump 50 specialized tools, design rules, and security policies into a system prompt, the agent suffers from instruction blindness and burns 40,000+ tokens per turn. Furthermore, agents overwrite each other's files when switching IDEs.

## 3. V1 → V2 → V3 Evolution
- **V1:** Manual jump tables and copy-paste prompt packets.
- **V2:** Monolithic Python router. Effective, but suffered from massive token bloat and environment lock-in.
- **V3 (Current):** Compiled Go kernel. 4-Stage Capability Pipeline. Cryptographic handoffs. Lazy capability loading. Pure separation of public engine vs private memory.

## 4. Architecture
```text
                    ORCHESTRA (Go Kernel)
                        |
          +-------------+-------------+
          |             |             |
      Cursor        Antigravity   Claude Code
       Adapter       Adapter       Adapter
```
Orchestra is one orchestration system with many interchangeable agent/tool adapters. It defines the contract; the agent executes it.

## 5. Workflow
1. `INIT`: Setup isolated workspace and schemas.
2. `CLASSIFY`: Analyze task for visual, security, and complexity demands.
3. `ROUTE`: Retrieve minimal required capabilities (Skills, MCPs).
4. `PLAN`: Generate Execution Manifest.
5. `HANDOFF`: Secure state transfer between agents.
6. `VERIFY`: Adversarial testing (Playwright, Lighthouse, Semgrep).

## 6. Capability Routing
Orchestra distinguishes between: `registered ≠ loaded ≠ used`.
If a task is a backend bug, visual UI skills are never loaded into the context window. If a capability gap is detected, Orchestra actively researches external docs before proceeding.

## 7. Multi-Skill Composition
For complex tasks, capabilities are synthesized. A premium frontend task merges `taste-design` (typography/layout), `emil-design-eng` (micro-interactions), and `Playwright` (verification) into a single coherent execution plan, preventing the agent from blindly copying generic templates.

## 8. Agent Allocation
- **Cursor:** Optimized for bulk component implementation and surgical IDE debugging.
- **Antigravity:** Optimized for multi-step orchestration, visual QA, and architecture enforcement.
- **Claude Code:** Optimized for terminal-driven server-side logic and independent auditing.

## 9. Model / Mode / Effort
Orchestra routes dynamically based on capability metadata (e.g., `visual`, `long-context`, `reasoning`), not hardcoded model names. The runtime matches the task requirements to the host agent's currently available models.

## 10. Token / Cost Philosophy
**Maximum useful quality per unit of token/cost.**
We utilize lazy capability loading and memory distillation to reduce average context size from 40k+ tokens to ~1,500 tokens. However, we *never* sacrifice quality for an arbitrary token quota. If a task requires deep 3D shader research, the tokens are spent.

## 11. Memory Model
- **Public Workflow (This Repo):** The clean Go engine, public schemas, and adapters.
- **Private Brain:** A strictly isolated local vault containing the user's durable preferences, project memory, and session states. A fresh clone creates your own Brain, not mine.

## 12. Verification
Agents cannot self-certify. Orchestra enforces adversarial verification using Playwright (multi-viewport visual QA), Lighthouse (Core Web Vitals), and static analysis tools.

## 13. Resource Ecosystem
Orchestra integrates seamlessly with an ecosystem of MCPs (Model Context Protocol), specialized Markdown Skills, and external libraries. (e.g., GSAP, R3F, Figma-to-Code MCPs).

## 14. Installation
```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow
go build -o orchestra ./runtime/cmd/orchestra
orchestra init
```

## 15. Clean Workspace
`orchestra init` creates a secure `.orchestra/` boundary. Your private preferences and project context never leak into the public workflow repository.

## 16. Examples

### Example 1 — Premium Business Website
**Routing:** `taste-design` + `impeccable` + `Playwright`.
**Action:** The system composes strict typographic scales and asymmetric layout rules before writing any React components, then visually verifies the result.

### Example 2 — SaaS Dashboard
**Routing:** `superpowers-planning` + `web-design-guidelines` + API/DB skills.
**Action:** Focuses on data density, accessibility, and backend integration.

### Example 3 — Creative 3D Portfolio
**Routing:** `taste-design` + `r3f-threejs` + `emil-design-eng`.
**Action:** Evaluates performance budgets, loads WebGL dependencies, and synthesizes 3D models with hardware-accelerated CSS motion.

### Example 4 — Backend Bug
**Routing:** `semgrep-adapter` + `core-logic`.
**Action:** Automatically drops all UI/UX capabilities. Token usage is minimized to focus entirely on stack traces and secure implementation.

### Example 5 — Mobile Application
**Routing:** React Native / Expo capabilities.
**Action:** Enforces native touch targets, safe area insets, and offline-first data fetching.

### Example 6 — Research Project
**Routing:** Academic lookup MCPs + citation parsing.
**Action:** Prioritizes factual extraction over code generation.

### Example 7 — Unknown Technology
**Routing:** `capability-gap-research`.
**Action:** If asked to use an unknown framework, Orchestra halts generation, reads official documentation, and dynamically builds a resource map before planning.

## 17. Benchmark (TTB Agro)
Orchestra V3 was empirically verified against the TTB Agro production deployment (SHA `3da1805`).
- **Results:** 100/100 Lighthouse Best Practices, 0 CLS, 804ms FCP, verified across 36 multi-viewport automated captures. Generic "AI slop" was completely eradicated.

## 18. Limitations
- Handoffs require explicit saving to `state.json`. Concurrent multi-agent editing of the same file without handoffs will cause cryptographic hash conflicts.
- Initial classification adds a ~2-second latency overhead to tasks.

## 19. Contributing
New Capabilities, MCPs, and Adapters can be added to the `registries/` directory. Ensure any visual skills adhere to the abstraction-first synthesis rules (do not force hardcoded CSS files).

## 20. Roadmap
- A/B Testing Raw Agents vs Orchestra V3 on identical tasks.
- Advanced WebGL/Shader capability routing optimization.

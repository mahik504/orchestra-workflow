# Original User Request

## Initial Request — 2026-09-03T20:03:47Z

# Teamwork Project Prompt — Draft

> Status: Launched
> Goal: Craft prompt → get user approval → delegate to teamwork_preview
> Requested team: Full team (Builds, Research, Ops)

Complete architectural rebuild of the Orchestra V3 resource ecosystem, replacing static skill lists with a dynamic capability graph, canonical machine-readable resource registries, and an automated resource acquisition/execution engine in the Go runtime.

Working directory: C:\projects\orchestra-workflow and C:\projects\orchestra-brain
Integrity mode: benchmark

## Requirements

### R1. Canonical Resource Registry & Graph
Audit all 65+ provided resources and existing Notion inventories. Create a canonical machine-readable registry (egistries/resources.json) preserving direct URLs, repository links, and capabilities. Create egistries/design-resource-graph.json mapping capabilities to resources (e.g., premium-website → Awwwards → Taste → GSAP).

### R2. Design-First Go Execution Engine
Update the Orchestra Go runtime router to use a multi-stage design engine: Discover → Classify → Research → Synthesize → Design System → Implement → Visual QA → Iterate. High-visual tasks must trigger extensive multi-source research before implementation.

### R3. Resource Acquisition Adapters
Implement execution adapters in the kernel for npm packages, Git repositories, CLIs, MCPs, and web references. The system must verify, conditionally install (Global/Project/On-Demand), and record provenance without universal global installation.

### R4. Host Synchronization & Cross-Agent Portability
Update Cursor, Antigravity, and Claude adapters/configurations to respect the new canonical capability contract. The 1,598-skill library must remain quarantined. Ensure the private Brain stores durable, structured resource memory (successes/failures).

### R5. Public Release & README
Completely rewrite the public README.md to reflect the new capability orchestration layer. Ensure strict public/private separation (no secrets or personal paths) and push the final orchestra-workflow and orchestra-brain repositories, recording actual GitHub SHAs.

## Acceptance Criteria

### Architecture & Registries
- [ ] egistries/resources.json and egistries/design-resource-graph.json exist and classify all supplied resources.
- [ ] Direct URLs, GitHub repos, and official docs are preserved for all resources.
- [ ] 1598-skill library remains quarantined outside the runtime.

### Runtime Behavior
- [ ] The Go runtime dynamically discovers, acquires, and conditionally activates resources based on task capability requirements.
- [ ] Premium visual tasks demonstrably research multiple sources (e.g., Awwwards, Jiro) and synthesize a design system before implementation.
- [ ] Project dependencies are installed automatically only when justified.
- [ ] Playwright visually verifies the implementation.

### Validation & Tests
- [ ] TTB Agro benchmark executed and recorded using the new design-first routing.
- [ ] Cross-agent test completed (Cursor, Antigravity, Claude Code) proving portability.
- [ ] 10 Design tests completed (Premium landing page, 3D portfolio, etc.) with recorded resource usage.

### Release
- [ ] Public README.md completely updated with diagrams and architecture explanations.
- [ ] Public repo pushed with verified remote SHA and zero leaked private data.
- [ ] Comprehensive 35-point final report (A-P) generated with evidence-based status.

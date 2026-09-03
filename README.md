# Orchestra V3

**A portable capability-orchestration layer for AI-assisted software development.**

Orchestra V3 replaces the monolithic "load every AI skill" pattern with a dynamic, design-first **Capability Graph**. Rather than guessing which scripts to run, Orchestra mathematically infers the correct research, synthesis, implementation, and QA resources based on a task's quality bar.

## Why Orchestra Exists
Generic AI coding tools guess dependencies and write vanilla code. Orchestra enforces an architecture where:
1. **Agents are Executors.** (Cursor, Antigravity, Claude Code)
2. **Orchestra is the Brain.** (The orchestration engine mapping capabilities to execution)

## The New Model: Capability Graph & Automatic Acquisition
Instead of manually typing "use GSAP" or "use R3F":

USER IDEA → CLASSIFY → RESEARCH → DISCOVER CAPABILITIES → RETRIEVE RESOURCES → SYNTHESIZE → IMPLEMENT → QA → REMEMBER

For a **Premium Visual Project**, Orchestra will:
1. Research design references (Awwwards, Taste, Impeccable).
2. Synthesize a Design System (DESIGN.md).
3. Conditionally acquire dependencies (GSAP, Lenis, React Three Fiber) only when needed.
4. Verify visually via Playwright and Lighthouse.

## Portability Contract
The capability graph is universal. Host adapters for **Cursor**, **Antigravity**, and **Claude Code** guarantee that the identical execution manifest is run across any environment, with a strictly synchronized active capability core.

## Resource Management (Zero Context Pollution)
- **Active Core (Small):** Loaded unconditionally (e.g. orchestra-conductor).
- **Curated Catalog (Large):** Loaded on gap-analysis (e.g. 3f-threejs).
- **Quarantine:** Unverified resources remain dormant outside runtime.

## Quality Over Minimum Tokens
Orchestra intentionally spends *more* tokens on premium visual tasks to research and synthesize multiple resources, prioritizing **Useful Quality / Cost** over naive brevity.

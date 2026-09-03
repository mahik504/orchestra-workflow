# Using Orchestra V3: Operator's Guide

This guide explains how to get the highest quality output from the Orchestra V3 kernel, regardless of which agent you use (Cursor, Antigravity, or Claude).

## Before Starting
To get optimal results, provide Orchestra with a strong initial context:
- The core idea or PRD.
- Client brief or target audience.
- Existing PDFs, references, or repositories.
- Desired visual style (if applicable).
- Hard constraints (e.g., "no 3D", "strict security").
- Objective success criteria.

## Routing by Task Size

### For Small Tasks
Use one agent and minimal routing. Bypass the full planning phase.
- **Example:** "Fix the JWT token expiration bug in `auth.ts`."

### For Medium Tasks
Engage the standard `classify -> plan -> execute -> verify` loop.
- **Example:** "Implement the new user settings page matching the existing design system."

### For Large Tasks
Use the full multi-agent orchestration pipeline:
- `Planning`: Architecture definition.
- `Allocation`: Break into subtasks.
- `Parallelizable work`: Assign frontend to Agent A, backend to Agent B.
- `Handoffs`: Use `state.json` to pass context between agents.
- `Verification & Integration`: Run Playwright and unit tests.

## Specialized Routing Domains

### For Premium Frontend
Explain how you want the visual iteration handled.
- **Action:** Explicitly state "take your time on visual iteration, use premium typography, and do not optimize for minimum tokens." Orchestra will load `taste-design` and `impeccable`.

### For Backend
Avoid polluting the context with UI rules.
- **Action:** State "This is a backend task." Orchestra will automatically drop UI/UX capabilities and prioritize `semgrep-adapter` and security protocols.

### For 3D / WebGL
- **Action:** Explicitly request "Use R3F and shaders." Orchestra will evaluate the performance cost vs visual benefit and inject the 3D capability ecosystem.

### For Research
- **Action:** Provide sources and ask Orchestra to verify claims. Orchestra will load documentation MCPs and research skills.

## Deployment & State Management

### Deploying vs Pushing
- **Push:** Commits to the existing Git repository.
- **Deploy:** Triggers a build to the existing Vercel/Netlify project.
- **New Target:** Orchestra will NEVER create a new deployment target unless you explicitly command: `create a new Vercel project`.

### Agent Switching (Handoffs)
You can start a task in **Antigravity** (for architecture and planning), and switch to **Cursor** (for rapid bulk implementation). 
Orchestra manages this via `internal/handoff/state.json`, using SHA256 checksums to ensure neither agent overwrites the other.

### User Feedback & Overrides
If Orchestra generates a plan you dislike, simply say "No, change X to Y." Orchestra will invalidate the current Execution Manifest and generate a new one before writing code.

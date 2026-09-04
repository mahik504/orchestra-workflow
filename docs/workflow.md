# Workflow

The conductor is the process you are talking to. The loop lives in [`WORKFLOW.md`](../WORKFLOW.md) and [`AGENTS.md`](../AGENTS.md). This page is a short pointer, not a second control plane.

```
Understand → Classify → Search graph → Design Lab / Technical plan → HUMAN GATE
 → Implement → Verify on the real app → Correctness review → Simplify review → Remember
```

1. **Classify** into one capability in `registries/design-resource-graph.json`.
2. **Plan** with `orchestra plan` (no side effects).
3. **Design Lab** is a write lock on PREMIUM / EXPERIMENTAL visual work until a stack card is approved. Say `skip the lab` to bypass one task.
4. **Implement** only the approved direction. Project-scoped installs. No global npm.
5. **Verify** on the running app. `DONE` / `VERIFIED` needs evidence in the same message.

Commands: `orchestra classify`, `orchestra plan`, `orchestra run`, `orchestra doctor`.

Say **skip orchestra** to stand the contract down for the session.

Host files (`.cursorrules`, `CLAUDE.md`, `kit/antigravity/MASTER-PROMPT.md`) are adapters. They translate syntax. They do not invent a second plan.

# Orchestra V3 Workflow

Orchestra governs execution via a strict pipeline. Agents may not write code until the pipeline completes.

## 1. Init & Classification
The user executes `orchestra init` to mount the workspace, followed by `orchestra classify [Task]`.
Orchestra evaluates the complexity, visual demands, and security requirements of the task.

## 2. Resource Gap Analysis
If the task requests a framework/technology not present in the local `registries/`, Orchestra halts and executes an external documentation lookup (Capability Gap Research).

## 3. Plan & Execution Manifest
Orchestra generates an Execution Manifest (e.g., `handoff/state.json`). This manifest specifies:
- Which agent should execute the task.
- Which specific skills (e.g., `taste-design`) are active.
- Which rules are explicitly disabled.

## 4. Execution
The target agent (Cursor, Antigravity, Claude) consumes the manifest and implements the changes. 

## 5. Verification
Visual tasks trigger Playwright captures. Backend tasks trigger static analysis. 
Output is either flagged for revision or approved.

## 6. Distillation
The session is aggressively summarized (Ponytail), and only durable project knowledge is committed to the private Brain.

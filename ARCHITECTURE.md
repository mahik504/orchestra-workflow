# Architecture

Orchestra V3 is powered by a pure Go kernel (`runtime/`).

The 4-Stage Capability Pipeline:
1. **Retrieval**: Classify the task and query the `registries/`.
2. **Analysis**: Extract capability metadata.
3. **Execution Manifest**: Compile explicit instructions, dropping any capability that exceeds performance budgets or conflicts with rules.
4. **Verification**: Adversarial testing gates (e.g., Playwright visual QA) run to ensure output matches the schema.

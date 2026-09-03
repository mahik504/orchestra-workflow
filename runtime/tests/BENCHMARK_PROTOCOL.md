# TTB Agro Pre-Implementation Benchmark Protocol

## Overview
This benchmark evaluates the performance of Orchestra V3 compared to a baseline agent on the "TTB Agro" test task. The goal is to measure improvements in robustness, context retention, and goal completion rates.

## Task/PRD
**Scenario:** Build and deploy a multi-service agricultural data processing pipeline for "TTB Agro" utilizing IoT sensors, a central processing API, and an analytical dashboard.
**Complexity:** Medium-High, involving multiple handoffs and context preservation across sub-tasks.

## Baseline Agent
- **Description:** Standard LLM agent with direct ReAct prompting without the explicit router/classifier orchestration.
- **Constraints:** Limited memory window, implicit state management.

## Orchestra Agent (V3)
- **Description:** Agent powered by the Orchestra V3 runtime (kernel, classifier, router, handoff, verification).
- **Constraints:** Strict schema enforcement for Plans, Tasks, Handoffs, and MemoryEntries.

## Quality Rubric
1. **Task Completion:** Are all micro-services functionally correct and correctly integrated?
2. **Context Retention:** Did the agent recall initial configuration requirements during the final integration step?
3. **Handoff Efficiency:** Were handoffs between capabilities (e.g., coding -> testing -> deployment) seamless without loss of data?
4. **Verification:** Did the agent self-correct using VerificationResults when errors occurred?

## Metrics
- **Success Rate:** Binary pass/fail for the overall scenario.
- **Handoff Count:** Number of explicit handoff events.
- **Error Recovery Rate:** Percentage of test failures successfully resolved via self-correction.
- **Time to Completion:** Total execution time (normalized by LLM generation speed).

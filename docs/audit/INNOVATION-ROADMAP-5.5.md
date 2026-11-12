# Innovation roadmap 5.5

Recorded 2026-11-12. An item is an innovation when the current tree does not already state it and the evidence says the absence costs correctness, tokens, or trust. Ordinary bugfixes are in the backlog, not here. Copied competitor features are rejected in the competitive analysis.

## Worth doing

### I-1. Compile tells the truth about a miss

- **Problem.** A task with no keyword returns confidence 1 and route `unspecified`. A host can treat that as a sure classification.
- **Existing approaches.** Keyword tables, and Superpowers’ “load the bootstrap every time,” which does not score a miss.
- **Why that is not enough.** Orchestra already has ambiguity for ties. The miss path forgot to lower confidence.
- **Change.** Set confidence to 0 and `fast: false` when nothing matches. Set `fast: true` when the winning route is L0 and the task contains a path.
- **Benefit.** The fast path becomes a manifest field, not a paragraph the model may skip.
- **Tokens.** Saves the conductor document on the fast path. Cost of the field is a few bytes of JSON.
- **Latency.** The scan is the same 36 routes.
- **Security.** None beyond not launching a premium loop by mistake.
- **Maintenance.** One function and a test.
- **Eval.** A compile fixture. Not a model trial.
- **Priority.** P0, because it is a wrong contract, and P1 for the fast flag.

### I-2. Task text outranks skill text

- **Problem.** SWE-Skills-Bench showed templates moving files to the skill’s paths and dropping the pass ratio by about 17 points on `asyncio-patterns`.
- **Existing approaches.** Longer skills, or no skills.
- **Why that is not enough.** Orchestra’s stubs invite improvisation. A long skill can override the user.
- **Change.** Every procedure that gets written starts with: the task’s paths, acceptance line, and skip line win. The skill is a method, not a spec.
- **Benefit.** Fewer harmful skills if one is loaded.
- **Tokens.** A few lines per skill, instead of a bundle.
- **Latency.** Neutral.
- **Security.** Stops a skill from redirecting writes.
- **Maintenance.** The 36 stubs, done as they are touched, not in one rewrite.
- **Eval.** A fixture where the skill names `src/wrong` and the task names `src/right`.
- **Priority.** P1 for the hot-path stubs.

### I-3. Evidence before a durable lesson

- **Problem.** A skipped visual check is stored as a successful resource outcome. The system can “learn” from a check that did not run.
- **Existing approaches.** Hermes writes skills from experience. Reflexion stores verbal feedback. Both trust the model’s story more than this repo should.
- **Why that is not enough.** Unverified success becomes the next route’s prior.
- **Change.** Memory writes use the real verdict. `NOT_APPLICABLE` and `NOT_RUN` do not increment success.
- **Benefit.** The ledger stops teaching the wrong lesson.
- **Tokens.** None.
- **Latency.** None.
- **Security.** Stops a poisoned “this resource works” row.
- **Maintenance.** One call site.
- **Eval.** A unit test with `CheckNotApplicable` and an acquired resource id. Assert the outcome is not success.
- **Priority.** P0.

## Already in the tree, so not innovations

The worker capsule, disclosure levels, `PromoteSkill`, the route skip line, the design-lab human gate, symlink refusal in `ow-sync.mjs`, and `AdmitDurable`. Using them is execution, not a new subsystem.

## Rejected as innovation theater

| Proposal from the assignment | Why it is not a 5.5 project |
| --- | --- |
| Adaptive context engine | Disclosure levels exist. The missing piece is I-1, not a new engine |
| Capability intelligence graph | `design-resource-graph.json` exists. Two selectors are a problem to simplify, not a graph to grow |
| Capability lifecycle tracking | The discovered → verified lifecycle exists. The bug is false evidence, I-3 |
| Dynamic model routing | The host owns the model. A hint field is enough |
| Visual intelligence as a new model | Routes and `ow-qa` exist. A Lighthouse wrapper is not a design judge |
| Autonomous QA that merges its own fixes | `ow-qa` already says reproduce, fix, re-run. Auto-merge is a release decision |
| Cost-aware orchestration that calls a billing API | No host in this repo exposes a bill to the compiler. Record tokens when the host prints them |
| Latency-aware execution | Not measured. A timer around compile would be a metric, not a product |
| Safe skill authoring that self-promotes | Contradicts SkillsBench’s self-generated arm |
| Intelligent subagent coordination | Five readers ran. Their reports disagreed on scope, and two arrived after the draft. More agents are not the design |
| Context-aware complexity beyond I-1 | The fast flag is the whole feature |

## Release priority

I-3 and the confidence miss go in the first code change, with the false PASS lines and the license removal. I-2 goes in as the stub skills are edited. Nothing in the rejected table is scheduled.

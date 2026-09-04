# Workflow — Orchestra 3.1.0

The contract is `AGENTS.md`. This is the operator's walkthrough.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

## 0. Mount

```
orchestra init      # create a local workspace next to this repo
orchestra doctor    # environment and host parity check
```

`init` writes a **local** workspace. It never points at someone else's vault.

## 1. Understand and re-brief

Paste the PRD, plan, bug, or idea. Orchestra restates it in one paragraph: archetype, quality bar, platform, hard constraints.

If two archetypes genuinely fit, you get **one** question. In an autonomous run with no answer, Orchestra picks the lower-risk archetype and logs `assumed <archetype>, no response`.

Correct the re-brief now. It is cheaper than correcting a build.

## 2. Classify

```
orchestra classify --task "<description>"
```

Produces task type, product archetype, quality bar, platform, visual importance, security level, research depth, and the verification set.

A restaurant site, a school management SaaS, and a 3D portfolio must not resolve to the same route.

## 3. Search the capability graph

```
orchestra plan --task "<description>"
```

Discover broadly, activate selectively. The chosen route loads fully — references, design skills, typography, motion, optional 3D. Other routes stay closed. Bulk vendor skill libraries stay quarantined.

If the graph has no strong match, Orchestra researches the technology and registers a new capability row instead of forcing a wrong archetype.

## 4. Design Lab (visual work)

`PREMIUM` and `EXPERIMENTAL` stop here. You get 2–3 directions and a stack card: typography, color world, layout language, component kit, one motion engine, 3D and shader decisions, logo method, icon system, implementation stack — each with a named source.

Approve, edit, reject, or combine. Rejected directions are logged with your reason so the next pass does not re-offer them.

`STANDARD` skips the lab unless you ask for it. Say **skip the lab** to bypass it for one task.

## 5. Implement

```
orchestra run --task "<description>"
```

Only the approved direction gets built. Implementation libraries install **project-scoped**. References are fetched on demand. Global installs are blocked.

## 6. Verify on the real app

```
orchestra verify
```

Launch the application and exercise it. For UI: screenshots at 2–3 viewports, zero console errors, no horizontal overflow on mobile, basic contrast check. For backend: tests and static analysis.

A description of a passing test is not a passing test.

## 7. Review twice

**Correctness** — bugs, missing behavior, security holes.
**Simplify** — duplication, needless abstraction, dead complexity.

These are separate passes with different questions. Running them together produces neither.

## 8. Remember

Durable outcomes go to your private workspace: which resource combinations improved the result, what failed, what you liked and hated. Transcripts and repository facts do not.

Only real executed jobs write resource memory. Synthetic rows are not learning.

## Overrides

| Say | Effect |
| --- | --- |
| `skip orchestra` | Stand down for the session |
| `skip the lab` | Bypass the Design Lab for one task |
| `ORCHESTRA_CONTRACT=<version>` | Pin a previous contract if a rollout misbehaves |

The graph recommends. You decide. Replacing the stack, framework, or agent mid-flight is expected, not a fight with a lock.

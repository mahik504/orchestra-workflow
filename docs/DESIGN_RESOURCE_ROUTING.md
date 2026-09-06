# Design resource routing — Orchestra 3.3.0

Orchestra synthesizes **one** design system. It does not collage libraries onto one page.

The instruction is:

> Go to the most relevant registered references, find useful patterns, combine the best **compatible** ideas from the available skill/resource/component libraries, synthesize one original design system, and implement it consistently.

Not: use everything.

## Pipeline

1. Understand the product.
2. Classify one capability (`registries/design-resource-graph.json`).
3. Set the quality bar (`STANDARD` / `PREMIUM` / `EXPERIMENTAL`).
4. Read the existing project design system if any.
5. Search registered research sources (Godly, Land-book, Lapa, SaaSFrame, Refero public — not Pro MCP unless authenticated).
6. Search registered component sources only as needed (`interaction-components`).
7. Search motion resources if the card named motion. **One** primary engine.
8. Search 3D resources if the card named 3D.
9. Build one coherent synthesis.
10. Write one `DESIGN.md`.
11. Human gate.
12. Implement.
13. Visual QA (Playwright on **our** app).
14. Accessibility review.
15. Performance review.
16. Simplification review.
17. Verification ledger.

Named override outranks the graph: a named site, skill, MCP, pack, or pasted `DESIGN.md` skips the 23-card survey.

## interaction-components

One routed capability. Not one skill per button library.

Search: shadcn, React Bits, Magic UI (free), Aceternity (free categories), Kokonut UI, Motion Primitives, Tailkit **public** pages, 21st **free** catalog, HyperUI, Cult UI.

Inspect **existing project components** before generating another control.

Button kinds: primary, secondary, destructive, loading, success, disabled, icon, magnetic, shimmer, gradient, ripple, particle, liquid, glass, command, hold, stateful.

Skip when the job is a full page, product, app, 3D scene, or native mobile — classify those routes first.

## Smallest sufficient motion / 3D stack

| Need | Stack |
| --- | --- |
| Simple UI | CSS / Motion primitive |
| Moderate motion | React Bits or Motion Primitives |
| Complex timeline | GSAP (`gsap-core`) |
| 3D | Three.js + R3F + drei |
| 3D + cinematic | add react-postprocessing |
| Scroll-driven 3D | add r3f-scroll-rig |

Do not install every animation library into every project. Do not stack GSAP with React Bits on the same surface unless the card says so.

## Combination examples (compatible, not collage)

### Premium marketing site

Refero public or Godly/Lapa research **+** shadcn **+** React Bits **or** Magic UI OSS **+** Kokonut (one named control) **+** one motion engine (GSAP if timelines).

### 3D portfolio

Refero/public research **+** shadcn for chrome **+** Three.js **+** R3F **+** drei **+** GSAP or CSS for non-GL chrome **+** Playwright QA **+** low-end fallback.

### SaaS product UI

21st free catalog **or** SaaSFrame **+** shadcn **+** Magic UI only if a named extra **+** `react-best-practices` / `web-design-guidelines`.

### Tint / adapt a reference

Screenshot-to-code or reverse-engineering **+** Refero public **+** taste-design **+** existing kits.

Workflow:

REFERENCE → extract design language → preserve structure → modify tokens → extend system → implement new features.

Example: orange → light blue means keep hierarchy, component grammar, spacing, interaction, shape language; replace color **roles**; adapt semantic colors; extend the new tokens into new sections. Do not do a superficial recolor. Do not clone trademarks, logos, proprietary assets, source code, or copy.

### Game-like 3D UI

`3d-portfolio` or `physics-canvas` (simulation vs presentation) **+** R3F/Three/drei **+** one motion system **+** browser QA. GetLayers stays pending.

## Honest access

| Source | Use now | Do not claim |
| --- | --- | --- |
| Magic UI | OSS + free MCP | Pro |
| 21st | search + 2 installs/day after login | Builder / unlimited installs |
| Refero | skill + public pages | free MCP |
| Tailkit | public pages | free MCP |
| Aceternity | public components | All-Access |
| GetLayers | nothing | connected MCP |

## Provenance

Acquired artifacts keep source, version/commit, license, method, date, and resource id. See existing provenance behavior in the engine. `registries/oss-index.json` is metadata, not a second Brain.

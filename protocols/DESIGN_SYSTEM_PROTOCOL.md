# Design system protocol

Mandatory when `visual_ambition` is premium or signature, or kind is hiring-cv / operator-HUD / showable hackathon web.

## Do not

Start implementing chrome from a blank file. Mix five kits. Default Inter + purple gradient. Nested cards. Random glass, blobs, aurora, decorative WebGL. Animate everything. Treat shadcn as the brand.

## Do

1. Classify visual intent (marketing / product / operator / native).
2. Collect **3–5** references (URLs or stills). Log them in `reference-log.md`.
3. Reverse engineer if the brief is “like this” (`REVERSE_ENGINEERING_PROTOCOL.md`).
4. **Originality gate (mandatory if any reference exists):** extract principles only — composition, typography, motion language, interaction, color *relationships*, spacing/grid, component *ideas*. Then write **our** direction. Never copy proprietary branding, logos, assets, copy, trademarks, or source CSS/JS. College one-block echo still tints; it is not a clone.
5. Propose **3** directions in plain words. the operator picks one.
6. Write artifacts (tiny projects may skip some):

| File | What |
| --- | --- |
| `DESIGN.md` | Point of view, anti-refs, type, color, space, motion, 3D yes/no |
| `design-tokens.json` | Color, type scale, space, radii, shadows |
| `motion-spec.md` | Personality, one engine, reduced-motion |
| `component-map.md` | Primitives vs custom vs one kit |
| `reference-log.md` | URLs, what was stolen (tokens not trademarks) |

7. Component strategy: **one** language. College: one kit + React Bits on hackathon web. Hiring: Stitch + one echo **or** named library vs DESIGN.md. operator HUD: Stitch Screen 1, no shadcn/21st, one volume.
8. Implement skeleton → components → motion → 3D if justified.
9. Visual QA protocol. Then polish.

## Typography (first class)

Follow `TYPOGRAPHY_PROTOCOL.md` (specialist capability, not an optional subsection). DESIGN.md must specify:

- Pairing (display vs text vs mono) and **why**
- Scale (steps, not random px)
- Measure (ideal 45–75ch for reading; HUD can be denser)
- Leading / tracking / optical sizing / variable axes if used
- Responsive type (min/fluid/max)
- Localization (no English-only tricks if product is IN)
- Platform: iOS Dynamic Type vs Android scaled sp vs web rem
- Fallback stack

Visual QA must screenshot **type**, not only color. Reject “Inter on everything” unless the product is a dense tool and Plan said so.

## Color / composition

One coherent palette. Contrast for text (WCAG as floor, not the aesthetic). Hierarchy: one focal point per view. Spacing from a scale.

## Motion

Write personality first. One primary system: Motion **or** GSAP **or** React Bits **or** CSS. Lenis only for scroll stories. Lottie only for authored files. Prefer transform/opacity. Honor `prefers-reduced-motion`.

## 3D / shaders

Use when it explains the product, is a signature hero, or is spatial data. Not for table UIs. operator HUD: one Shadertoy/drei volume, not a CSS planet, not a screenshot-loop GLSL.

## Anti-slop

No generic purple SaaS, no icon-tile dashboards as identity, no five libraries, no design-by-registry.

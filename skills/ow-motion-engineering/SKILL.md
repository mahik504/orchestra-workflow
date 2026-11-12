---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-motion-engineering
description: Decide whether something should move, which runtime owns it, and how it respects reduced motion. GSAP owns scroll. Cursor, transitions, and exits get durations. One runtime per surface.
---

# Motion engineering

Use this when the approved direction includes motion. `gsap-core`, `gsap-scrolltrigger`, and `gsap-timeline` still own scroll. This skill decides the rest.

## Before any tween

1. Should it move? If the motion does not explain state, hierarchy, or feedback, leave it still.
2. Who owns it? Scroll is GSAP. A React component's enter and exit can be CSS or the project's existing runtime. Do not add Framer Motion beside GSAP for the same surface.
3. Properties. Transform and opacity. Do not animate layout (`top`, `width`, `height`) on a scroll-linked element.
4. Interrupt. Hover and focus must be able to reverse. A scroll timeline that fights the user is a bug.
5. Reduced motion. `prefers-reduced-motion: reduce` gets the end state, or a crossfade under 200ms. No pin, no parallax.

## Cursor and transitions

A custom cursor is optional and only on a fine pointer. It must not replace the native cursor for buttons and links. View transitions, when the product already uses them, cover route changes. They do not replace the scroll owner.

## Budget

60fps on the target laptop. If a timeline drops frames, cut effects before you cut content. Measure on the running page, not in a screenshot.

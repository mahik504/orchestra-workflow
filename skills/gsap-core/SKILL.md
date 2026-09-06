---
name: gsap-core
description: GSAP core motion. Use when the approved Design Lab motion engine is GSAP. One motion engine per product. Not stacked with React Bits + Motion.dev on the same surface.
license: MIT
---

# GSAP core (Orchestra)

Official skill (MIT): host GSAP plugin / https://gsap.com  
GSAP is fully free including former Club plugins.

## Orchestra rule

The Design Lab picks **one** motion engine. If the card says GSAP, use this skill. If it says React Bits / Motion / CSS, do **not** also load GSAP.

Use for timelines, ScrollTrigger, interruptible UI, SVG morph, reduced-motion via `gsap.matchMedia()`.

Project-scoped `gsap` npm. Not a global install. Not the conductor.

Related official skills exist (timeline, ScrollTrigger, React, plugins). Load those **only** when this job needs them — not as extra global Orchestra skills.

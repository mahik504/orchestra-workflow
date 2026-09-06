# Resource combination

Orchestra synthesizes **one** design system. It does not collage libraries.

```mermaid
flowchart TD
  brief[Brief]
  research[Research sources]
  kits[Component kits]
  motion[One motion engine]
  three[3D stack only if needed]
  synth[One DESIGN.md]
  brief --> research
  brief --> kits
  brief --> motion
  brief --> three
  research --> synth
  kits --> synth
  motion --> synth
  three --> synth
```

Smallest sufficient stack: CSS/Motion primitives → React Bits / Motion Primitives → GSAP for timelines → R3F/Three/drei for 3D → postprocessing for cinematic → scroll-rig for scroll-driven 3D.

Full examples: [`docs/DESIGN_RESOURCE_ROUTING.md`](../DESIGN_RESOURCE_ROUTING.md).

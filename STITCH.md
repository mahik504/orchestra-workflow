---
title: Stitch
---

# Stitch (2D chrome) + 3D (code)

Stitch designs **screens**. It does not design real-time 3D. Cursor (and Antigravity) build 3D **inside** those screens.

```
Plan brief
  → optional SkillUI on ONE echo URL (see below)
  → taste-design → DESIGN.md
  → Stitch: layout, type, color, components, mobile frames
  → lock projects/<slug>/design.md
  → Agent: implement chrome to match Stitch
  → Agent: R3F or Spline only where Plan marked 3D
  → earned motion (hover / reveal / micro) — Emil; React Bits if needed
  → optional Antigravity polish
  → Playwright on web
```

## Echo a live site (SkillUI)

When the plan brief names **one** site we may echo (e.g. Nothin’, K95), extract tokens — do not clone the page:

```text
npx skillui --url https://example.com --format design-md --out "C:\projects\orchestra-brain\projects\<slug>\ref"
```

Ultra mode (screenshots + hover diffs) only if default mode is too thin:

```text
npx skillui --url https://example.com --mode ultra --format design-md --out "C:\projects\orchestra-brain\projects\<slug>\ref"
```

Then lock **our** `projects/<slug>/design.md`. Delete `ref/` after Stitch is locked if it is just a third-party scrape. Never `skillui` Notion/Linear “because it looks good.” Never install SkillUI as a global always-on skill pack.

Our own repo → Stitch: use existing `stitch-extract-design-md` / `stitch-code-to-design`, not SkillUI.

- **R3F** (`@react-three/fiber` + `drei`): interactive 3D in React (orbs, globes, product heroes).
- **Spline**: authored scene, drop into React, less code.
- **Ultron orb** (`C:\projects\ultron-orb`): later craft reference, not a dependency.
- One 3D beat per view unless Plan says otherwise. Motion still follows Emil (earned, short, `transform`/`opacity` on the 2D chrome).

Campus WiFi may block Stitch APIs. Then use https://stitch.withgoogle.com and save shots into `.stitch/designs/`.

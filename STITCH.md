---
title: Stitch
---

# Stitch (2D chrome) + 3D (code)

Stitch designs **screens**. It does not design real-time 3D. Cursor (and Antigravity) build 3D **inside** those screens.

**Stitch is not required every time.** If Plan named a library, template, shader, or animation, combine those against `projects/<slug>/design.md` and tint.

```
Plan brief
  → read kind
  → named source? combine against design.md (skip Stitch)
  → else college: SkillUI / kit block / React Bits / one ThreeUI
  → hiring: Stitch + one echo or named library
  → operator HUD: Stitch screens, no shadcn; 3D = volume R3F
  → lock projects/<slug>/design.md
  → Agent: implement the lock
  → one motion engine — React Bits on college web
  → optional Antigravity polish
  → Playwright on web
```

## Echo a live site (SkillUI)

When the plan brief names **one** site we may echo (e.g. Nothin’, K95), extract tokens — do not clone the page:

```text
npx skillui --url https://example.com --format design-md --out "C:\projects\orchestra-brain\projects\<slug>\ref"
```

Ultra mode only if default mode is too thin:

```text
npx skillui --url https://example.com --mode ultra --format design-md --out "C:\projects\orchestra-brain\projects\<slug>\ref"
```

Then lock **our** `projects/<slug>/design.md`. Delete `ref/` after Stitch is locked if it is just a third-party scrape. Never `skillui` Notion/Linear “because it looks good.” Never install SkillUI as a global skill pack.

## Dual-track import (by kind)

| kind | 2D import | Forbidden |
| --- | --- | --- |
| college / hackathon | One SkillUI URL or one 21st / shadcn / unlumen / smoothui / neo block or React Bits or one ThreeUI. Tint to one world. React Bits required on web. Stitch optional | Mixing five kits. Generating chrome from zero |
| hiring-cv | Stitch + one echo **or** named library against design.md | Kit look as the product |
| personal operator HUD | Stitch. No shadcn. No 21st. Volume R3F | Kit HUD, CSS planet |

Always-on 21st / Aceternity / Magic MCP stay off. `neo` = neobrutalism.com (college / loud marketing). Ditther / Pryzm / notyourtype = Plan-named craft. react-spring = one motion engine option. Do not clone screenshot-to-code.

Academic monochrome workspace only if Plan named that world. Do not mix it onto a hiring ML product.

Our own repo → Stitch: use existing `stitch-extract-design-md` / `stitch-code-to-design`, not SkillUI.

- **R3F** (`@react-three/fiber` + `drei`): interactive 3D in React.
- **Spline**: authored scene, drop into React, less code.
- One 3D beat per view unless Plan says otherwise. Motion still follows Emil (earned, short, `transform`/`opacity` on the 2D chrome).

Campus WiFi may block Stitch APIs. Then use https://stitch.withgoogle.com and save shots into `.stitch/designs/`.

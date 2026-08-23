# Typography protocol

Specialist capability `typography`. **Not** a section you skip inside DESIGN.md.

Activate on every high-visual job (`frontend_premium`, marketing, editorial, operator HUD HUD, hiring site, showable hackathon web). Load this file. Do not load a typography skill dump.

No high-quality **official** single agent-skill beat this protocol + platform docs (Apple HIG, Material 3 type, MDN fonts). Third-party “font pairing” skills are OPTIONAL at best and were not installed.

## DESIGN.md must lock

| Topic | Required |
| --- | --- |
| Pairing | Display vs text vs mono, **why**, max **two** families unless Plan says otherwise |
| Hierarchy | Roles: display / h1–h3 / body / label / numeric; weight contrast, not size-only |
| Scale | Named steps (e.g. 12/14/16/20/24/32). No random px |
| Measure | Body ~45–75ch; HUD/dashboard denser if Plan says |
| Leading / tracking | Body ~1.4–1.6; display tighter; never `line-height: 1` on paragraphs |
| Variable fonts | Prefer one VF file over 6 static weights when licensed |
| Optical sizing | `font-optical-sizing: auto` when the face has `opsz`; don’t fake with scaleX |
| Responsive | `clamp(min, preferred, max)` or platform Dynamic Type / scaled sp |
| Loading | `font-display: swap` or optional; subset; avoid huge Google CSS |
| Localization | IN product: don’t rely on English-only ligature tricks; test Devanagari if copy is Hindi |
| Platform | iOS Dynamic Type + Dynamic Type styles; Android sp + Material type roles; web rem/clamp |
| Fallbacks | System stack after the branded face |
| License | Confirm we may use the files commercially |

## Pairing rules

- Contrast **structure** (serif display + grotesque text, or geometric display + humanist text). Two similar geometrics = mush.
- Don’t pair two display faces.
- Mono only for code, IDs, telemetry — not marketing body.
- Reject Inter + purple as identity unless a dense **tool** Plan named it.

## Platform notes (official)

- Apple HIG — Typography: https://developer.apple.com/design/human-interface-guidelines/typography
- Material 3 type: https://m3.material.io/styles/typography/overview
- MDN variable fonts / `font-optical-sizing`: https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_fonts

Expo/RN: use platform text styles / `PixelRatio`, not copied web `px`.

## Visual QA

Screenshot type at desktop + mobile. Check wrap, widows on heroes, tabular nums for data, contrast on text, reduced-motion doesn’t shrink type to unreadable.

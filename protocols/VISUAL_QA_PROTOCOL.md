# Visual QA protocol

Cannot claim “looks good” without this on a visually significant UI.

Pipeline: Taste → DESIGN.md → implement → Impeccable → Playwright → Designer SOS → **one** consultant packet.

## Required

1. Launch **our** app. Canonical Layer A owner is the Orchestra harness (`orchestra` / `PlaywrightVerifier` / `runtime/tools/layer-a.mjs`) against that running origin, using the project-local Playwright **library** in `runtime/package.json` plus Chromium (`npx playwright install chromium` from `runtime/`). `npx playwright --version` is the Test CLI and does not prove Layer A. Playwright MCP and Cursor browser may be used to *drive* the app; they do not green-wash the run. `playwright-cli` is optional; do not `install --skills -g`. Mock `MOCK_SCREENSHOT_*` files are not evidence.
2. Desktop screenshot
3. Tablet screenshot
4. Mobile screenshot
5. Accessibility (contrast as a number, focus, labels, reduced motion)
6. Interaction (click/type the main path). Also `:active`, `:focus-visible`, empty / error / disabled — not hover-only chrome.
7. Spacing vs DESIGN.md scale
8. Typography vs DESIGN.md
9. Motion (jank, duration, reduced-motion)
10. Hierarchy (one focal point)
11. Responsive wrap/overflow
12. Performance smell (huge images, 3D on /dashboard tables)
13. Fill `protocols/VERIFICATION_LEDGER_PROTOCOL.md`
14. Designer SOS six-pass (`protocols/DESIGNER_SOS_PROTOCOL.md`): Squint, Grid, Type, Colour, Space, Cut; finish Keep / Kill / Park

If a reference still exists: compare **ours vs reference** (language, not clone). Anatomy notes, not copied assets.

## Regression

After the contract and the golden stills are approved, later UI changes compare against those stills. Mask clocks, random ids, and canvases that are not the subject. Animated regions use a paused frame or a fixture. A pixel difference outside the mask fails when it exceeds the project's threshold. Do not treat a perfect score as the default. Lighthouse and axe-core run against the same origin. A phone width, a keyboard pass, reduced motion, and the WebGL fallback are part of the run when the page has those states. A mock screenshot cannot pass.

Run, when loaded: frontend critique (taste-design), **impeccable** after a keepable still, **web-design-guidelines** (fetch live rules). High-end: **one** consultant packet (Antigravity or Opus Task) — not a swarm.

## Fail

Generic SaaS, Inter-default, nested cards, motion everywhere, 3D on CRUD, operator HUD with shadcn.

## Not a substitute

A single screenshot in chat is not verification.

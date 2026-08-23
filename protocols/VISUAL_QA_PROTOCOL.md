# Visual QA protocol

Cannot claim “looks good” without this on a visually significant UI.

## Required

1. Launch **our** app (Playwright MCP or Cursor browser).
2. Desktop screenshot
3. Tablet screenshot
4. Mobile screenshot
5. Accessibility (contrast, focus, labels, reduced motion)
6. Interaction (click/type the main path)
7. Spacing vs DESIGN.md scale
8. Typography vs DESIGN.md
9. Motion (jank, duration, reduced-motion)
10. Hierarchy (one focal point)
11. Responsive wrap/overflow
12. Performance smell (huge images, 3D on /dashboard tables)

If a reference still exists: compare **ours vs reference** (language, not clone).

Run, when loaded: frontend critique (taste-design), **impeccable** after a keepable still, **web-design-guidelines** (fetch live rules). High-end: a **second** pass (Antigravity packet or Opus Task).

## Fail

Generic SaaS, Inter-default, nested cards, motion everywhere, 3D on CRUD, operator HUD with shadcn.

## Not a substitute

A single screenshot in chat is not verification.

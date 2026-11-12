---
name: ow-site-recon
description: Measure a public page into a blueprint of tokens, motion, components, and the calls that page makes. Hand an original build onward. Do not copy assets, copy, or a logo.
---

# Site recon

Extends `ow-reference-forensics`. The output is a measurement, not a clone.

## Preconditions

- The URL is public. No login, paywall, or CAPTCHA bypass.
- Playwright is available for the measurement pass.
- The product repo is the place the notes are written.

## Procedure

1. Open the page in Playwright. Record computed styles, `:root` variables, font families, spacing, and breakpoints.
2. Record `document.getAnimations()`, any GSAP or ScrollTrigger instances you can see, easings, and the `prefers-reduced-motion` path.
3. List the visible components and the routes the public nav exposes.
4. Record the public network calls made during that ordinary visit. Draft a Zod shape only for calls you actually saw.

## Output

Write `SITE_BLUEPRINT.md`, `tokens.json`, `MOTION_SPEC.md`, and a forbidden-identity list (name, logo, copy, unlicensed fonts). Label every line observed, inferred, or unknown.

## Verify

A second person can tell which lines were measured and which were guessed. No asset, sentence, or logo from the source is in the product repo.

## Skip

Stop at a login wall, paywall, CAPTCHA, or a robots disallow. Do not use a stealth browser.

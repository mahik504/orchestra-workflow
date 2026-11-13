---
name: ow-reference-forensics
description: Named live site or screenshot recreation. Measure composition, separate mechanism from identity, write REFERENCE_BENCHMARK.md. Load only on the ow-reference-forensics route. Not a clone factory.
---

# Visual forensics

Job-loaded. Skip unless the brief names a live site, a screenshot recreation, or “feel like X”, or the `ow-reference-forensics` route is the classified capability.

## Skip

- Tint / keep-the-structure token swap → `reverse-engineering`
- No named site, screenshot, or public source
- Intent is to copy branding, assets, copy, or source
- A locked `DESIGN.md` and golden stills already exist and the job is implementation
- Backend, paper, or security with no UI

## Produce

Write `REFERENCE_BENCHMARK.md` with:

1. URL or still source
2. Captured shots (desktop and mobile if the source is live)
3. Measured layout (grid, type roles, density, motion cues)
4. Interactions that matter
5. Mechanism vs identity
6. Translation map (what we will keep as principle)
7. **Forbidden identity** (logo, copy, trademarks, source, distinctive marks we will not ship)

Hand off to Design Lab. Named visual references get **exactly 3** evidence-backed translations, not 23 cards and not a skip to one `DESIGN.md`.

Study the reference as principles: layout, type, grid, hierarchy, color relationships, spacing, interaction, scroll, light, camera, scene composition, loading, and the responsive change. Record those principles. Do not copy branding, assets, copy, trademarks, or source.

## Reverse engineering a page

Write `REFERENCE_AUDIT.md` next to the design contract:

1. What the page is for.
2. The grid, the type roles, and the color relationship you measured.
3. The interaction that matters, and the one you will not copy.
4. A search for an equivalent component already in the product or in the chosen kit.
5. Authorization. A public page can be studied. A logged-in app, a paywall, or someone else's repo is out unless the task names it as own, client, or open source.

`rea` is the binary path, and only after Java is installed and the target is yours, a client's, a CTF, or open source. A website audit is not an invitation to pull private APIs.

## Preconditions

- The task names a public page, a screenshot, or a target the operator owns, a client owns, or that is open source.

## Procedure

1. The task's paths and acceptance line win.
2. Measure grid, type, motion, and the one interaction that matters. Write the forbidden list: logo, copy, assets, and source.
3. Adapt the principle in the product's own components. A copied code block needs a license that allows it.

## Output

`REFERENCE_BENCHMARK.md` with the measurements and the forbidden list.

## Verify

The written notes contain no pasted asset, no copied paragraph, and no unlicensed source block.

## Skip

Skip when the intent is to copy the site, when no source is named, or when a locked contract and approved stills already exist.

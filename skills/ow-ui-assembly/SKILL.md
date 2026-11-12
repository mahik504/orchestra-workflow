---
name: ow-ui-assembly
description: Assemble a screen from the product repo first, then one shadcn registry source per section. One motion runtime. Record the license and the source.
---

# UI assembly

Supersedes a blind search of every component kit.

## Preconditions

- A direction exists. This skill does not invent the visual world.
- The product has tokens. New components use those tokens.
- One motion runtime is already chosen. GSAP owns scroll when the route says GSAP.

## Procedure

1. Search the product repo for an existing component that does the job.
2. If it does not exist, search the shadcn MCP across the namespaced registries in `components.json`.
3. Take one source per section family. Map it to the product tokens.
4. If the sample uses `motion/react` and the page is GSAP, port the motion to `useGSAP` instead of loading a second runtime.
5. Record the package, version, license, and URL in `.orchestra/components.json`.

## Output

Components in the product repo and a provenance file. At most two visual dialects on a page.

## Verify

The page still builds. The provenance file names every imported block. Reduced motion still has a path.

## Skip

Do not load a second motion engine. Do not install a registry the product does not use. 21st.dev only when a key already exists.

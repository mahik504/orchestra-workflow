---
name: r3f-shaders
description: Build one React Three Fiber scene with a still fallback and a frame budget. Use for 3D, shaders, or WebGL. Skip for ordinary page layout.
---

# R3F and shaders

Use one scene when the product needs it. Do not use 3D as a background.

## Scene

- One Canvas for the scene. Lazy-load that module so the rest of the page paints first.
- Wrap the scene in Suspense. Show a still image while it loads.
- The fallback is a local image. It is not a spinner and not a hotlinked file.
- Clamp device pixel ratio to the range 1 through 2.
- When the visitor prefers reduced motion, render the still and do not start the render loop.
- On unmount, dispose geometries, materials, textures, and the renderer.

## Motion

One motion stack for the page. Do not run a second engine inside the scene. Scroll-driven camera work uses the scroll skills.

Experimental work still needs the design gate and golden stills.

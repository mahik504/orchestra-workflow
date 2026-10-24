---
name: r3f-shaders
description: Build a React Three Fiber scene with a fallback and a frame budget. Use for 3D, shaders, or WebGL. Skip for ordinary page layout.
---

# R3F and shaders

Use 3D when the product needs a scene. Do not use it as a background decoration.

| Rule | Detail |
| --- | --- |
| Fallback | A still image or a simple mesh when WebGL fails |
| Budget | 60 fps on a laptop. Clamp device pixel ratio |
| Motion | One engine. Honor reduced motion |
| Cleanup | Dispose geometries, materials, and textures |
| Load | Show a suspense or loading state |

Experimental work still needs the design gate and golden stills.

---
name: ow-3d-engineering
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

## Acceptance

- The scene explains something. A decorative object does not ship.
- Camera, lights, and materials are chosen. They are not leftovers from an example.
- Assets load behind a local still. The rest of the page paints first.
- Dispose geometries, materials, textures, and the renderer on unmount.
- Cap device pixel ratio from 1 to 2 unless the contract names a tighter cap.
- The render loop does not run for a hidden or reduced-motion visitor.
- A phone gets the still, or a scene the contract says is cheap enough.
- If WebGL fails or the context is lost, the still remains and the page still works.
- Do not claim a frame rate you did not measure. Record the budget and the machine.

## Performance budget

Keep one canvas. Pause it offscreen. Prefer a still over a second postprocessing pass. Theatre.js, a scroll rig, and a bloom pass are added only when the contract names them.

## Which tool

| Need | Use | Leave it |
| --- | --- | --- |
| One scene in a React page | React Three Fiber and Drei | A second canvas, or Three.js beside R3F |
| A shader | GLSL or TSL in that scene | A shader on a page with no scene |
| A marketing loop the human will edit visually | Spline, exported, with the still fallback | Spline plus a hand-rolled scene of the same object |
| A film or a print asset | Blender | Blender as a runtime in the browser |
| WebGPU | Only when the contract names it and a WebGL2 fallback exists | A blank page when WebGPU is missing |

The frame budget is 60fps on the machine you measured. DPR stays in `[1, 2]`. Dispose geometries, materials, textures, and the renderer. Reduced motion and a lost context keep the still.

## Preconditions

- The task names one scene, or it names a film asset that stays in Blender.

## Procedure

1. The task's paths and acceptance line win.
2. One React Three Fiber canvas, lazy, with a local still fallback. Blender exports a glTF. It does not run in the page.
3. Dispose on unmount. Do not add a second motion engine.

## Output

The scene module, the still path, and the machine the frame budget was checked on.

## Verify

The page paints the still when WebGL is absent, and the render loop stops under reduced motion.

## Skip

Skip for ordinary layout, for 3D used as wallpaper, and for a second canvas.

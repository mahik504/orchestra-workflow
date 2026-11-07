# Motion

`ow-motion-engineering` owns the motion decision. GSAP owns scroll. One runtime per surface. `gsap-react` loads only when the motion lives in React components.

Cursor, hover, and route transitions stay on transforms. Honor `prefers-reduced-motion`. A pinned section names its start and end.

3D is a separate route. One scene loads `ow-3d-engineering`. The still fallback and the frame budget are in [3d.md](3d.md).

# Technology selection

Choose the smallest stack that the task needs. Confirm the current license before adding a package. Do not install globally.

| Need | Choice |
| --- | --- |
| A document | HTML and CSS |
| A small script | JavaScript, no application build |
| A production application | TypeScript |
| Components | React |
| Server rendering, routing, SEO | Next.js |
| A client dashboard | Vite, when the brief names a client bundle |
| A content site | Astro |
| Scroll, pin, scrub, a cinematic timeline | GSAP |
| Hover, layout, drag | Motion for React |
| Smooth scroll, and only then | Lenis, one scroller shared with ScrollTrigger |
| A scene outside React | Three.js |
| A scene inside React | React Three Fiber, with Three.js as that scene's renderer |
| Controls or loaders for that scene | Drei |
| A custom GPU effect | GLSL |
| A scene-level effect that was asked for | React Postprocessing |
| An art-directed timeline that was named | Theatre.js |
| DOM and 3D sharing scroll, and that was named | R3F Scroll Rig |
| The deliverable is a Framer project | Framer, and no application stack |
| Durable relational data | Postgres |
| The product note names Supabase | Supabase on Postgres |
| Cache, TTL, rate limit, ephemeral state | Redis, and not as the source of truth |
| A small custom HTTP API | Express |
| Python inference or data processing | FastAPI |
| A named background job that already needs Redis | BullMQ |

One scroll owner. One scene. One HTTP stack. One site bundler. Motion for React is not Framer. GSAP and Motion may share a page only when they do not drive the same element or the scroll container.

Redis, a queue, and a second runtime wait for a requirement. A missing library is reported. It is not installed.

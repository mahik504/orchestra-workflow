# Next starter

App Router, TypeScript, and Tailwind. Lint, typecheck, and build run in CI. The Lighthouse workflow prints scores and is allowed to fail. It does not deploy.

```
npm install
npm run dev
```

`npm run typecheck` follows `npm run build`. Next writes route types into `.next` during the build, and the typecheck reads them.

Copy `.env.example` to `.env.local` when the app needs a public site URL.

Motion and 3D are not installed here. Add them when the design contract names GSAP, Lenis, or one React Three Fiber scene.

A Vercel preview is optional. From this directory, after `npx vercel login`, run `npx vercel`. This starter does not deploy from GitHub, and it does not ship a production deploy on every push.

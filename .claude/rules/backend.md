---
description: Server and API input. Validate at the boundary.
globs: "**/*.{ts,js,go}"
alwaysApply: false
---

- Validate server and API input with Zod or the equivalent at the boundary.
- Return one error shape. Do not invent a second backend stack beside `ow-backend-architecture`.
- TypeScript stays strict.

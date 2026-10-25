---
name: backend-architecture
description: Shape a SaaS or API before code. Routes, auth, schema, migrations, errors, and logs. Use for a backend or API. Skip the design lab.
---

# Backend architecture

Write the shape before the handlers. There is no 23-card survey.

## Shape

| Piece | Decide |
| --- | --- |
| Routes | The resource, the verb, and who can call it |
| Auth | How a caller is identified, and what a missing session returns |
| Schema | Tables or collections, keys, and what is unique |
| Migrations | How an existing database moves forward. Do not silently rewrite live data |
| Errors | One error shape. The client can tell a bad request from a server failure |
| Logs | What is logged, and which fields are never logged: tokens, passwords, session cookies |

## Then

`test-driven-development` for the behavior. `supabase` and `postgres-best-practices` only if this product uses them. `ship-safe` when the change touches auth, public routes, or secrets.

Do not add a service, a queue, or a second database until a route needs it.

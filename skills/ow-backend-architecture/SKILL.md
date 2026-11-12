---
name: ow-backend-architecture
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

`test-driven-development` for the behavior. `supabase` and `postgres-best-practices` only if this product uses them. `ow-ship-safe` when the change touches auth, public routes, or secrets. A multi-tenant SaaS also follows `references/saas-checklist.md`.

The product note picks Supabase, Firebase, Clerk, Prisma, or Drizzle. This skill does not.

Do not add a service, a queue, or a second database until a route needs it.

## Data stores

| Store | Use it when | Leave it out when |
| --- | --- | --- |
| Postgres | The data is relational, transactional, and must survive a restart | The page has no durable data |
| Supabase | The product note already names managed Postgres, auth, row level security, storage, or realtime | You only needed Postgres |
| Redis | A cache, a TTL, a rate limit, or ephemeral state is a requirement | You want a source of truth |
| Node | The server is TypeScript and it is not a Python service | The service is FastAPI |
| Express | The HTTP API is small and the repo does not already have Next or FastAPI | Next or FastAPI already owns the route |
| FastAPI | The work is Python inference, data processing, or model serving | The UI server can do the job |
| BullMQ | A background job needs a Redis-backed queue | The work can finish in the request |

Say which tables exist, how they relate, which constraints and indexes the queries need, and how a caller is authenticated and authorized. Validate input. Use one error shape. Do not log tokens, passwords, or session cookies. Name the migration, the test, and where it is deployed.

## Data and auth

Pick from the project note. Do not add a second database because a tutorial used one.

| Piece | When it belongs |
| --- | --- |
| Postgres | Relational data. Neon or Supabase when the note names a host. SQLite when the app is local or a single node. |
| Redis or Upstash | Cache, rate limit, or a short-lived queue. Not the source of truth. |
| Prisma or Drizzle | The note names one. Do not introduce the other beside it. |
| Auth | The note names Clerk, Supabase Auth, or the existing session. Row level security is a policy, not a slogan. |
| Queues | BullMQ only when a job must leave the request and Redis is already justified. |

Write the migration, the constraint, and the test that fails when the constraint is missing.


## Preconditions

- The task names the files or the acceptance line.

## Procedure

1. The task's paths and acceptance line win.
2. Write the route, auth, schema, and error shape before the handler.

## Output

The change, the command that checked it, and what that command printed.

## Verify

Say the work is done only after that command has been run.

## Skip

Skip when the human says skip orchestra. Do not add a second conductor.

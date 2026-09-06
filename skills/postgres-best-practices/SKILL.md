---
name: postgres-best-practices
description: Postgres schema, indexes, query mistakes. Use on Supabase/Postgres schema work. Not an ORM dump.
---

# Postgres best practices (Orchestra wrapper)

Upstream: https://github.com/neondatabase/postgres-skills  

## When

Schema design, slow queries, missing indexes, RLS-adjacent SQL. Skip for frontend-only and Expo-without-db jobs.

## How

1. Normalize what must be consistent; denormalize what is read-heavy and named.
2. Index the filters you actually query.
3. Do not invent a second database next to Supabase "for Orchestra."

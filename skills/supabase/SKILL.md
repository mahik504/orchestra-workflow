---
name: supabase
description: Supabase database, auth, RLS, storage, edge functions. Use when the product actually uses Supabase. Not a second backend framework. Secrets stay in the host.
---

# Supabase (Orchestra wrapper)

Upstream: https://github.com/supabase/agent-skills  
Do **not** install the entire pack blindly. This wrapper is the Orchestra entry.

## When

The idea.md or repo already names Supabase, or the human asks for hosted Postgres + auth. Skip for SQLite-only college GUI, operator HUD without cloud, and papers.

## How

1. Schema and RLS before UI chrome.
2. Keys in `.env` / host MCP. Never git.
3. If MCP is `AUTH_REQUIRED`, record the Connect step. Do not fake HEALTHY.
4. Prefer official CLI/MCP for migrations.

Orchestra remains the conductor. Supabase is a capability.

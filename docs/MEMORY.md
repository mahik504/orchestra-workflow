# Memory

Two stores. They are not the same thing.

**Resource telemetry** is `resource-memory.json`. A row is a resource id, an outcome, and a sha256 of the task text. New rows do not store the raw prompt. Notes are redacted when they look like a key. The file mode is `0600` and the directory is `0700`. On Windows that mode is not a boundary against the same user. The ledger keeps the last 500 rows. `orchestra memory scrub` rewrites an older row whose task text is not already a sha256.

**Vault knowledge** is the private workspace. Taste, decisions, and one product note live there. A lesson draft from `/ow-retro` stays in `.orchestra/handoff/lesson.md` until a person says yes. Nothing in that draft is copied into a skill or into `AGENTS.md` by itself.

`LoadDurable` can read eight facts. The main pipeline does not call it.

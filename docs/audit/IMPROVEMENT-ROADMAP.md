# Improvement roadmap

Order is the 4.5.0 plan. Each item names the file that will prove it.

1. `registries/versions.json` is the only current version. Check scripts read it. Historical docs keep their old numbers and gain a banner.
2. `registries/skills.json` lists every canonical skill once, with owner, license, and `renamed_from` where a skill moves to `ow-*`.
3. `runtime/tools/ow-sync.mjs` copies by hash, skips conflicts, and can restore. PowerShell and bash installers call it.
4. `runtime/tools/ow-doctor.mjs` reports expected, missing, extra, and hash mismatch per host. `--private` is the only mode that looks at the vault.
5. CI validates every registry that has a schema, runs `node --test`, and runs a Windows PowerShell job. Actions stay pinned. `permissions: contents: read`.
6. `registries/routing.json` renders the adapter blocks. Codex, OpenCode, and Kimi get adapters from that render.
7. Orchestra-owned skills rename to `ow-*` after `skill-aliases.json` records the old id. Upstream names stay.
8. `orchestra compile` prints the task manifest as JSON. Golden fixtures live under `runtime/testdata/compile/`.
9. Admission states are separate from run states. A fit score stores the evidence string beside it.
10. `ow-skill-advisor` searches and recommends. It does not run a global install. Hard limits stay in `AGENTS.md`.
11. Design, motion, 3D, reverse engineering, delivery, and workflow skills are `ow-*` files with references, not a pack install.
12. Context budget is a measured before/after from `measure-context.mjs`. Model classes name a job, not a vendor model id.
13. Bookmarks from the operator's export become registry rows or vault notes. Junk is listed as dropped.
14. Benchmarks c1 and c6 run only when the headless agent is actually available. A failed or skipped run stays `not_run` in the README.
15. README, changelog, tag `v4.5.0`, and the GitHub release happen after the checks above are green.

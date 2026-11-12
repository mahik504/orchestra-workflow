# Memory and learning 5.5

Recorded 2026-11-12. The vault stays the memory store. No new database is recommended. Hermes was read as a README and a root listing, not as a system to install.

## What exists

| Piece | Role | Wired to a run? |
| --- | --- | --- |
| Private vault, path from `ORCHESTRA_VAULT` | Preferences, decisions, project `idea.md`. Template under `workspace-template/`. The skill `ow-vault` is `stub: true` | The conductor is told to read it. The Go pipeline does not open it. This session did not open the vault |
| `AdmitDurable`, `SelectDurable`, `MaxDurableFacts = 8` in `reconcile.go` | Class, status, and a cap. Accepts preference, architecture decision, validated workflow, failure lesson, verified resource, project fact. Rejects chat, tool calls, retries, screenshots | `LoadDurable` is not called from `pipeline.go`. `ConflictWinner` is called from a test only |
| `resource-memory.json` | Per-resource success and failure counts | `IterateStage` writes it, and a skipped visual check is stored as success if the verdict string passed in is `"PASSED"` |
| `memory.go` `Store` | SQLite-shaped API | `NewStore` does not open a database and has no callers. `RetrieveByTags` returns an empty slice |
| `ApplyMemoryPolicy` | Last failure can suppress a resource | Reported by the security reader. Not re-opened |

## Classes

| Kind | Store | Treated as truth? |
| --- | --- | --- |
| Personal preference | Vault, user scope | Yes for this operator, until they change it |
| Project convention | Vault project folder, or the product repo’s own docs | The repo file wins over a note when they disagree. `kit/antigravity/ALWAYS-ON.md` line 5 says that if the note and `AGENTS.md` disagree, `AGENTS.md` wins. That line was opened in the red-team pass |
| Validated workflow | Vault, only after a check passed | Yes, with the command and the date |
| Failure lesson | Vault, only after a check failed for a reason that was fixed or accepted | Yes, as a lesson, not as a ban on every similar tool |
| Unverified suggestion | Nowhere durable | No |
| Temporary instruction | The task, this session | No |
| Global engineering knowledge | This public repo, after a human accepts the diff | Only what was reviewed |
| Another project’s facts | Never copied across project folders | No |

## Loop

CAPTURE the command output and the route that was used.

CLASSIFY with the `AdmitDurable` reasons. If the reason is chat or a retry, drop it.

VALIDATE against the check. A skipped check is not a success. A model’s summary of a page is not a citation until the URL was opened.

STORE in the vault for operator facts, and in `resource-memory.json` only for a resource the run actually used. Cap the selection at 8 facts. Do not grow a transcript.

RETRIEVE by the project folder and a few tags. No vector index. If the note and the repo disagree, the repo wins.

APPLY as a constraint on the next task in that project. Do not apply a project note to a different product.

MEASURE on the next similar task. If the lesson does not change the outcome, mark it historical.

PROMOTE a skill, route, or rule only through a reviewed diff. ARCHIVE when a newer decision names `superseded_by`. The field already exists on `Fact`.

## What the system may propose

A shorter skill, a route-row fix, a deleted stub, a registry correction. The proposal is a diff the human can read. It is not applied to the global instructions by the same session.

## What it must not do

- Promote a skill it wrote in the same session. SkillsBench’s self-generated arm scored below the no-skill baseline.
- Edit `AGENTS.md` or the route table from a memory write.
- Copy vault text into the public repo.
- Mix two projects’ notes.
- Treat `RetrieveByTags` as search. It returns nothing.
- Install Hermes, Mem0, or Letta to get this loop. Hermes’s README describes FTS5 session search and skills created from experience. The state modules were not read. The learning claim conflicts with the self-generated-skill result. Orchestra can record a verified lesson without adopting that loop.
- Store secrets, browser cookies, or `.env` values.

## Hermes, specifically

Useful to know, not to vendor: a user-level memory and a session search exist as a product. Orchestra’s equivalent is the vault plus an 8-fact cap plus a resource ledger that must stop lying about skipped checks. The installer `curl | bash` and `iex (irm ...)` is rejected. The “only agent with a learning loop” sentence in their README was not audited against other products.

## First code changes, when authorized

1. `recordIterateMemory` must receive the real verdict. `NOT_APPLICABLE` stores nothing, or stores a non-success outcome that `ApplyMemoryPolicy` will not treat as a win.
2. Delete or quarantine `memory.go`’s unused `Store` so a reader does not think keyword memory works.
3. Call `SelectDurable` only from a command the operator runs on purpose, such as `orchestra memory list`, and only on vault files under the project. Do not inject the vault into every compile.
4. Give `ow-vault` a one-page procedure that matches `AdmitDurable`. It is a stub today.

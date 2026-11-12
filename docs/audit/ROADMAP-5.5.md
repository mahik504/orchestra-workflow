# Roadmap 5.5

Recorded 2026-11-12. Orchestra’s pin was 5.0.0 at commit `24ad53c3ca4e4098491da7ec7f29683670c09cf3` when this page was written. The release commit moves the pin to 5.5.0. This page stays the plan from before that commit.

## What 5.5.0 is

A correctness release. The doctor, verify, sync, acquire, and compile commands stop reporting success for a check that did not run. The public tree no longer contains the Anthropic document skills whose license text forbids distribution. The route table, the changelog, and the conductor document tell the same version and the same routes. Hot-path stubs gain a short procedure in which the task wins. The claim against other harnesses stays none.

## What 5.5.0 is not

A new brain. A model router. A larger skill catalog. A second conductor. A benchmark win. An install that overwrites personal config. A ten-layer architecture.

## Sequence

1. **License.** B-07. Human confirms the four document skill trees leave the public repo. No tag before that.
2. **False PASS.** B-01, B-02, B-03, B-04, B-05, B-06, B-08 in one commit. Tests in the same commit.
3. **Hooks and inspect.** B-13, B-14.
4. **Words and routes.** B-09 for the hot-path stubs, B-10 so `ow-qa` is a real row or the changelog stops saying it is, B-11 fast-path flag, B-12 one pin, B-17 delete the unused SQLite store.
5. **Shelf.** B-16 move the unloaded stacks.
6. **Evidence, optional for the tag.** B-15. If it is not run, the tag still says the claim is none. If it is run and `PromoteSkill` returns false, the tag still says none.

B-18 waits on a re-read. B-19 is maintenance, not the release.

## Gates

| Gate | Pass condition |
| --- | --- |
| Correctness | B-01 through B-06 and B-08 have tests that fail on the old behavior |
| Tests | `go test ./...` in `runtime/` and the Node tool and hook tests exit 0 on the release commit. CI on that commit is success. A previous green run on `24ad53c` does not count for a later commit |
| Security | B-07 done. B-13 and B-14 done or explicitly deferred in the changelog with the residual risk written down. No new always-on MCP. No `curl \| bash` installer |
| Architecture | `ARCHITECTURE.md` describes the two speeds and does not add a second pipeline. `docs/CONDUCTOR.md` matches the pin and is not a second route table |
| Skills | Hot-path stubs listed in B-09 either have the five sections or still say `stub: true`. No skill added to a Load line without a `PromoteSkill` record |
| Hosts | Go and Node host maps list the same set, including Codex. Missing directories are `NOT_RUN` |
| Install | The master prompt refuses to overwrite a file that is not already an Orchestra adapter. The prompt is in `MASTER-INSTALL-PROMPT.md`. It was not executed |
| Docs | Changelog, `VERSION`, `registries/versions.json`, and plugin manifests share one pin. This roadmap is not that pin until the code lands |
| Eval | `compare.go` stays disabled unless a saved trial makes `PromoteSkill` return true. The trial file is in the repo |
| Performance | `AGENTS.md` size in `docs/CONTEXT-BUDGET.md` matches the file. The measure script’s rewrite is reviewed, not left as a drive-by |
| Consistency | `git grep` does not find `4.6.0` or `4.5.0` as the contract pin in the live docs |

## After 5.5.0

Run T1 and T16 from `EVAL-5.5.md` before any sentence about beating a raw host. Diff upstream Pocock and Superpowers copies before refreshing them. Re-read the junction test before claiming the quarantine check was proven.

## Rollback of a bad 5.5.0

The previous tag is `v5.0.0`, peeled commit `24ad53c`. A bad 5.5.0 is reverted by checking out that commit or by a revert commit. Sync backups from `ow-sync.mjs` restore host copies that were replaced. Personal files that were not Orchestra adapters should not have been touched. If they were, the backup stamp is the restore path, and the release failed its install gate.

# Security and governance 5.5

Recorded 2026-11-12. This is a threat model of the method repo. Nothing was exploited. No `.env`, no live MCP config, and no vault directory was opened. Findings are IMPLEMENTED in source unless a line says otherwise. Local tests of the hooks were NOT_RUN.

The operator is one person on a local IDE and CLI. This repo is not a multi-tenant server. The trust boundary is: the model and any fetched page are untrusted, the operator’s approval is required for destructive writes, and the public git tree must be legal to redistribute.

## Assets

- The public method repo and its MIT license.
- Host config the installer might touch: `~/.cursor`, `~/.claude`, `~/.claude.json`, `~/.gemini`, `~/.codex`, `~/.agents`, `~/.jcode`.
- The private vault path, named by `ORCHESTRA_VAULT` and required to stay out of public files. `hygiene.js` treats that folder name and its absolute path as leak patterns. The script was not run.
- Resource-memory JSON under `.orchestra/memory/` or `ORCHESTRA_MEMORY_PATH`.
- Credentials in `.env`, `id_rsa`, `.pem`, `credentials.json`, and MCP json files. The hook’s read denylist names those shapes.

## Actors

- The operator, who can approve a hook `ask`.
- The model, which can be prompted by repository text, web pages, and skill files.
- A dependency or skill author. The Go module has no third-party `require`. Node locks Playwright at `1.63.0` and uses caret ranges for `ajv`, `lighthouse`, `mermaid`, and `@axe-core/playwright`. The lockfile exists. A vulnerability scan was NOT_RUN.
- A malicious page returned to `Acquire` or `Inspect`.

## Controls that exist

`kit/hooks/guard.mjs`, default profile `standard`:

- Denies `npx skills add --all` and `skills add -g -y` when `--all` or the global yes-install matches the regex as written.
- Denies `git push` with `--force` or a standalone `-f`.
- Denies `curl` or `wget` piped to `sh` or `bash`.
- Denies shell reads of a path containing `.env`, and reads of `.env`, `id_rsa`, `.pem`, `credentials.json`, and live MCP json.
- Asks on global npm installs, `git reset --hard`, and `git push`. Profile `minimal` turns `ask` into `allow`. Profile `strict` turns `ask` into `deny`.
- Profile `off`, or an id listed in `ORCHESTRA_DISABLED_HOOKS`, allows the action before any rule runs. The hook test file asserts that `off` allows `skills add --all`. That test was not run.
- `stop` denies a completion claim without command-shaped evidence only in `strict`. Other profiles allow it and attach a follow-up.
- Edit always allows. A note can mention a BOM, an `sk-` or `ghp_` shaped string, or a `SKILL.md` without frontmatter. The edit is not blocked.
- Subagent start and stop write `.orchestra/handoff/subagent.json`. Strict mode denies only when the payload sets `wrote_outside`.

`runtime/internal/runner/runner.go` starts a process without a shell. The security reader inspected it. This writing pass did not re-open the file, so the metacharacter rules stay REPORTED: the executable name rejects `|;&$\`><`, and arguments are rejected only for null bytes. A test in `runtime/tests/m3_challenger_test.go` is reported to pass a semicolon as a git argument and expect no sentinel file. That test was NOT_RUN.

`ow-sync.mjs` `refuseEscape` throws on a symlink and on a real path that leaves the skill root. `listSkillIds` skips a top-level symlink. `ComputeSkillDirHash` in `host_sync.go` uses `os.ReadFile`, which follows a symlink. The two hashers are not the same function. The Go hasher was not re-opened in the writing pass. The split is REPORTED by the security reader and is consistent with `refuseEscape` being only in the Node file, which was opened.

`acquisition/web.go` blocks localhost, loopback, private, link-local, and `169.254.169.254` on the request and on redirects, including DNS answers. If `LookupIP` returns an error, the URL is not rejected. That branch was reported by the security reader and was not re-opened line by line. The offline-fixture success path was opened (F-05 in the baseline).

`npm.go` rejects global installs. `git.go` rejects a URL that starts with `-`. `mcp.go` blocks a command string that is exactly `cmd.exe`, `powershell.exe`, `bash`, or `sh`. A full path or `pwsh` does not match an equality check. That equality gap is REPORTED and was not re-opened.

`kit/git-hooks/commit-msg` strips one trailer, `Co-authored-by: Cursor <cursoragent@cursor.com>`. Other `Co-authored-by` lines stay. REPORTED from the security reader. The hook file was not re-opened.

Hygiene CI scans tracked files for secret-shaped tokens and for the vault path strings. It was not run locally. The remote Hygiene conclusion on `24ad53c` was success, which means that job’s scanners did not fail. It does not mean every secret shape was covered.

## Findings

### S-01. False assurance

F-02, F-03, F-04, and F-06 in `BASELINE-5.5.md`. A doctor or verify line that says PASS when the check did not see the thing it names trains the operator to trust a broken gate. Severity: high for a control plane whose product is correct reporting. No remote attacker is required.

### S-02. Offline fixture as a successful acquire

F-05. A failed fetch becomes evidence with `source: network_fallback_fixture`. A later stage can treat that JSON as a page. Severity: high if acquire results are stored as verified resources. The resource lifecycle says a row stays discovered until a run records evidence. A nil error is evidence-shaped.

### S-03. Inspect has a weaker network policy than Acquire

`onboard/inspect.go` uses an `http.Client` with a 12 second timeout, a 5-redirect cap, and no private-IP check. Body cap is 64 KB. Prompt-injection consequence: the page is labeled untrusted in a comment, and classification still reads it. Severity: medium. The fetch is operator-triggered by a URL the operator submitted. A redirect to a metadata address is the residual risk.

### S-04. Hook bypasses

The `--all` regex does not cover `skills add some-pack --all` or `pnpm dlx skills add --all`. The pipe-to-shell regex does not cover PowerShell `iex`. Profile `standard` does not stop a false DONE. Profile `off` stops nothing. Severity: medium. The operator can still approve or refuse. The hook is not a sandbox.

### S-05. MCP allowlist is a substring

`mcpAllowed` returns true when `kit/mcp-core.json` is missing. When the file exists, it checks that the server name appears as a quoted string. A name that is a substring of an allowed name can pass. Severity: low to medium, depending on server names. The core file was not parsed in this pass.

### S-06. Proprietary skills in a public repository

`skills/docx/LICENSE.txt` and `skills/pdf/LICENSE.txt` say the materials may not be extracted from Anthropic’s Services, retained outside them, or distributed to any third party. They are tracked in this public repository. `THIRD_PARTY_NOTICES.md` does not name them. `pptx` and `xlsx` were not opened. Severity: high for redistribution, not for remote code execution. This pass does not delete the files. A release should remove them, or replace them with a host-provided pointer and a notice a lawyer has read. This document is not legal advice.

### S-07. Empty parity and a Codex blind spot

`ResolveHostDirs` omits `~/.codex/skills`. `ow-sync.mjs` includes it. `orchestra sync` can report full parity while Codex was never compared, and while every directory was absent. Severity: medium. It is an integrity bug, not a secret leak.

### S-08. Junction test that does not fail

The security reader reported `TestChallenger_Quarantine_SymlinkOrJunctionAnalysis` logs a junction evasion and does not call `t.Error`, and that a failed `mklink` still passes the test. The test file was not re-opened. Status: REPORTED. If the report is accurate, the test cannot catch a quarantine escape. Confirm before treating it as a proven bypass.

### S-09. Supply chain

Go code in `runtime/` uses the standard library only (`go.mod` has no `require`). That removes a class of module attacks. Node dependencies are locked for CI via `npm ci`. Caret ranges can move when the lockfile is regenerated. `kit/mcp-core.json` pins filesystem `@2.3.0`, Playwright MCP `@0.0.85`, Context7 `@4.6.0`, and Stitch MCP `@0.9.0`, and the lint fails `@latest`. No CVE list is included. A scanner was not run.

`skills/docx` is the supply-chain item with a license conflict. Superpowers skills carry `LICENSE.superpowers` and are named in `THIRD_PARTY_NOTICES.md`. That part of the notice file was opened only at the top. The mapping for Superpowers was visible.

### S-10. Memory isolation

Resource memory is a JSON file on disk. `SerializePublic` is reported to refuse a public artifact that contains the private vault folder name unless `publicSafe` is set. The function was not re-opened. The vault skill says the vault root must not be pasted into a public file. `ConflictWinner` is not on the execute path, so two contradictory notes are not resolved by the pipeline. Last-write or last-failure overlay behavior is a separate function (`ApplyMemoryPolicy`), reported and not re-opened. Severity: medium if a public file is generated from vault text. The hygiene pattern is the backstop that was not re-run.

## Threat model, short

| Threat | What exists | Residual |
| --- | --- | --- |
| Prompt injection from a fetched page | Comment plus SSRF checks on the web adapter | Inspect is weaker. Offline fixture returns success. No HTML sanitizer |
| Secret in git | Hygiene patterns, hook read deny, edit note | Edit is allow. Hygiene was not re-run locally. Remote Hygiene succeeded |
| Destructive shell | Hook denylist and ask list | Bypassable regexes. Profile `off` disables them. Runner is not a shell, arguments are loose |
| Skill install of the internet | Hook refuses one `skills add --all` shape. Advisor skill says not to install unless the human insists | Other installers, including gstack’s clone-into-`~/.claude`, are outside the regex |
| Symlink escape during sync | Node `refuseEscape` | Go hash follow is a different function. Junction test status is REPORTED |
| Cross-host desync | Sync hashes and backups | False PASS when directories are missing. Codex omitted from the Go map |
| Private memory in a public commit | Hygiene string match, `publicSafe` flag | Not re-tested. Vault was not opened |
| Unlicensed redistribution | Notices for vendored MIT skills | Anthropic document skills contradict their own license text |
| False completion | Strict stop hook | Default profile allows the claim |

## Governance for 5.5

- PASS is allowed only when the check’s input existed and the predicate was true. A missing file, a missing directory, and a skipped visual check are `NOT_RUN` or `NOT_APPLICABLE`, and they are not stored as `OutcomeSuccess`.
- The public tree drops the Anthropic document skills or replaces them with a one-line pointer. `THIRD_PARTY_NOTICES.md` lists every remaining third-party tree.
- Install never overwrites a file that does not already look like an Orchestra adapter. Bootstrap already claims this for `.gemini/mcp_config.json`. The claim was not re-tested.
- New skills are not promoted by the agent that wrote them. `PromoteSkill` stays the gate.
- Hook profile `standard` should deny, not merely note, a completion claim with no command output, if the product wants the contract “evidence before DONE” to be enforced. Today it is enforced in prose and in `strict` only.
- No new always-on MCP. `kit/mcp-core.json` stays an explicit list.
- Rollback of a sync is the backup stamp `ow-sync.mjs` already writes. Uninstall is delete the copied adapter only when the file matches an Orchestra hash. That behavior is described in the changelog and was not re-run.

## Release gate

A 5.5.0 tag is blocked while S-06 is true, while S-01 prints PASS, and while a comparison claim is enabled without a saved trial. Tests of the hook and of `verify` exit codes belong in that gate. They were NOT_RUN in this research.

# Resource selection 5.5

Recorded 2026-11-12. Selection is a decision procedure. It is not a new index. The repo already has `registries/resources.json` (208 entries), `registries/design-resource-graph.json` (33 capabilities, 42 domains, version 5.0.0), `registries/skills.json`, and `registries/host-stack.json`. Presence in a registry is not usage. `ARCHITECTURE.md` says a resource counts as used only when a run recorded evidence.

## Order of search

1. The product repo the user is in. Existing components, routes, and conventions.
2. This method repo: the route table, then the one skill, then the registry row if the skill names a resource.
3. The private vault, for preferences and project facts. Not for public code. Not pasted into a commit.
4. Official documentation of the library the task already uses.
5. A component registry, one source for the section, when the product does not have the control.
6. GitHub, for a named open-source tool the task asks for. Pin a commit. Read the license.
7. A paper or a primary page, when the claim depends on a result. Record the URL.
8. A design reference, for composition only, when the task is a new visual identity.

Stop at a login wall, a paywall, a CAPTCHA, a robots denial, a 403, or a 429. That is the scrape skill’s rule. It applies to every search, not only to `ow-scrape`.

## When to reuse, compose, or reject

| Situation | Action |
| --- | --- |
| The product already has the component | Reuse it |
| One registry primitive plus the product’s tokens fits | Compose |
| The skill is `stub: true` and the task is a bug in a named file | Do not load the stub. Fast path |
| The skill’s text conflicts with the task | The task wins |
| The resource is `creative-3d-web` or `fullstack-mastery` | Reject. The route already names the replacement |
| The resource is `docx`, `pdf`, `pptx`, or `xlsx` inside this repo | Reject for redistribution. Use the host’s tool |
| A second motion runtime, a second HTTP stack, or a second bundler | Reject. Technology policy already says so |
| An MCP for a job a CLI or a file read does | Reject the MCP |
| The license forbids distribution, or the maintainer is unknown and the install is global | Reject |
| The trial for a new skill does not exist | Keep it off the default Load line |
| Two routes tie | Ask. Do not pick the longer keyword |

## MCP and plugins

Native tools, the shell, and a direct API come first. An MCP is justified when the host cannot do the call safely any other way and the server is in `kit/mcp-core.json`. A missing core file must not mean “allow every server.” Today `mcpAllowed` returns true when the file is missing. That is a code defect, not a selection rule.

Always-on MCPs are rejected. Playwright for a visual route is the one the compiler adds, and only for `frontend-design-high`. A backend bug does not get a browser server.

Plugin install of Superpowers, gstack, GSD, ECC, or Compound Engineering is a second conductor. The hook should keep refusing a blanket `skills add`. The human can still insist. The advisor skill already says that. The agent does not insist for them.

## Graph and indexes

`design-resource-graph.json` is the capability catalog the classifier reads. `orchestra compile` reads `routing.json`, not the graph. Two selectors can disagree. A later change should make one of them the authority for “which skill,” and the other a catalog of resources that skill may use. This pass does not merge the files.

A new relationship graph, vector index, or embedding service is not justified. The 208-row catalog is already more than the eval has exercised. `Promote` refuses a jump from discovered to successful. Keep that.

## Provenance

A selected external resource records: URL or repo, commit or version, license name, and the date it was opened. `Acquire` must not record an `offline_fixture` as a successful fetch. Today it can. Until that returns an error, a network failure is not a source.

## Confidence

| Signal | Confidence |
| --- | --- |
| File found in the product repo | Use it |
| Route keyword, unique winner | Use the route |
| Tie or no keyword | Ask |
| Registry row with no `run:` or `evidence:` | Discovered only. Do not call it verified |
| README claim with no file opened | Unknown. Do not select on it |
| Skill with `stub: true` | The name may be right. The body is not a procedure yet |

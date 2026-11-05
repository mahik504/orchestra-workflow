# Orchestra 4.0.0

[![CI](https://github.com/mahik504/orchestra-workflow/actions/workflows/ci.yml/badge.svg)](https://github.com/mahik504/orchestra-workflow/actions/workflows/ci.yml)
[![Hygiene](https://github.com/mahik504/orchestra-workflow/actions/workflows/hygiene.yml/badge.svg)](https://github.com/mahik504/orchestra-workflow/actions/workflows/hygiene.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/mahik504/orchestra-workflow?include_prereleases&sort=semver)](https://github.com/mahik504/orchestra-workflow/releases)

A control plane for agentic development. One conductor. The same skills on Cursor, Antigravity, Claude Code, and jcode. A human gate before showable UI.

Orchestra is a markdown contract the agent already knows how to read, plus an optional Go engine for the parts a contract cannot enforce. Clone it. Keep products in a private workspace. This repository is the method.

The pin in `VERSION` is **4.0.0**.

```
ORCHESTRA  = CONTROL PLANE
SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES
AGENTS     = EXECUTORS
BRAIN      = MEMORY
REGISTRY   = RESOURCE KNOWLEDGE
```

Available is not loaded. A registry row is not proof a tool ran. A discovered model name does not pick a skill.

---

## What you get

- One contract, [`AGENTS.md`](AGENTS.md). Host files are adapters. If two hosts could plan, the chat you are talking to conducts.
- A routing table. The agent prints `Load:` and `Skip:` before it writes. A skill you name wins. A missing tool stops. It does not install itself.
- A design gate. Flagship and premium visual work stays unwritten until you approve a contract and stills. College work and bugs skip the lab unless you say otherwise.
- A private vault for taste, decisions, and one lesson line per session. That vault is not in this clone.
- Core MCPs: Playwright for your running app, Context7 for current library docs, and GitHub on each host after you sign in. The token stays in the host.
- A Next.js starter at [`templates/starters/next-react`](templates/starters/next-react) with lint, typecheck, build, and a Lighthouse job that warns.

Architecture: [`ARCHITECTURE.md`](ARCHITECTURE.md). Figures: [`docs/diagrams/`](docs/diagrams/).

![Control plane: the human talks to one conductor](docs/diagrams/assets/control-plane.svg)

```
Understand → Classify → Load the route → Design Lab or a technical plan → HUMAN GATE
  → Implement → Verify on the running app → Review twice → Remember
```

![Eight-stage loop with a human gate before implementation](docs/diagrams/assets/the-loop.svg)

The engine wraps that loop with a mission, a worker capsule, and a budget derived from the task. Ordinary work stays cheap. A flagship page can take a few directions, then one contract. Playwright checks the running app. A model label does not change which skills are selected.

---

## Advantages

- The agent names the skills for the job before it writes, so a bug does not load a design survey and a page does not load every kit.
- Showable UI has a gate: contract, then stills, then code. A markdown contract is not visual evidence.
- The same skill folders sync to the hosts you pick. Marketplace packs are not the install.
- Memory lives in your private workspace. This public repo stays a method.
- GitHub, Playwright, and Context7 are part of the working set once you sign in. The agent can read the repo, review a diff, and open a pull request.

## Disadvantages

- A one-line Go fix was slower with Orchestra than with the host alone. Both arms made the same 183-byte change. That pair does not score a website.
- There is no study in this repository that shows Orchestra produces better pages than a raw host. A structural difference is not a score.
- The design gate adds a human step on flagship and premium UI. Say **skip the lab** when you want it off for one task.
- The engine does not choose a model for you. Opaque model names are ignored on purpose.
- A commit made inside Cursor can attach the wrong author. Push this repo from your own terminal.
- Agents do not silently push `main`. GitHub access is for reading, review, and pull requests. You approve the push.

---

## Compared with earlier releases

| | 3.3.1 | 3.3.5 | 3.3.6 | 4.0.0 |
| --- | --- | --- | --- | --- |
| Contract | Design gate, named references take 3 translations | One contract, skills are the folders in `skills/` | Same contract, plus a routing table and a required `Load:` / `Skip:` line | Same architecture. Eight routed skills. 73 canonical skills. Cap 80. No second conductor. No new global MCP |
| Checks | Governance release | Pin moved. CI and Hygiene were red on that commit | Humanizer is files in the clone. CI reports `gofmt -l` and does not rewrite | Same checks. The routing checker matches the job cell exactly |
| Design bar | Premium visual work uses the lab | Kind sets the bar. College skips the lab unless you say PREMIUM | Same bar. The table names which skills load | Same bar. One component source in a fixed order. A catalog is not a motion runtime |
| Starter | None in this line | None | Next.js App Router starter with lint, typecheck, build, and a warning Lighthouse job | Same starter. No new packages in this method repo |

Older tags remain: `v3.3.1`, `v3.3.0`, `v3.2.0`. `v3.3.5` stays on the commit that moved the pin. Rollback notes: [`kit/ROLLBACK.md`](kit/ROLLBACK.md).

## Compared with a raw chat and other harnesses

| | What it does | Where Orchestra differs |
| --- | --- | --- |
| Raw Cursor, Antigravity, or Claude Code | The host writes from the prompt and whatever skills happen to be in context | Orchestra prints the skills for this job, keeps one conductor, and blocks showable UI until you approve stills |
| Superpowers | A methodology plugin. Useful as a worker for planning habits | It does not replace `AGENTS.md`. If it starts a second plan, the contract wins |
| Ralph Loop | Repeats until a stop condition | It must not bypass the design gate or become the conductor |
| This repo | The method, the skill folders, and the engine | Your products, keys, and taste stay in a private workspace |

---

## Install

```bash
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow
```

Windows:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/bootstrap.ps1
```

macOS or Linux:

```bash
chmod +x kit/bootstrap.sh kit/init-workspace.sh
./kit/bootstrap.sh
```

Set `ORCHESTRA_VAULT` to your private workspace. Sign in to GitHub, Playwright, and Context7 on Cursor, Antigravity, and Claude Code. Paste keys in the host UI. Restart. Talk.

Walkthrough: [`docs/getting-started.md`](docs/getting-started.md).

Say **skip orchestra** to stand down for the session. Say **skip the lab** to bypass the design gate for one task.

### Install prompt

Copy this into a new agent chat after the clone is open.

```
You are installing Orchestra 4.0.0 from this repository. AGENTS.md is the contract. You are the conductor for this chat. Do not start a second plan.

Do this:
1. Read AGENTS.md, ARCHITECTURE.md, and docs/getting-started.md.
2. Run kit/bootstrap.ps1 on Windows, or kit/bootstrap.sh on macOS or Linux. Wire the hosts I name. Default: Cursor, Antigravity, Claude Code.
3. Copy skills/ onto those hosts with kit/sync-ides.ps1 -Sync. Do not run a marketplace install-all.
4. Tell me to set ORCHESTRA_VAULT to my private notes folder. Do not invent the path. Do not print secrets.
5. Tell me to sign in to GitHub, Playwright, and Context7 on Cursor, Antigravity, and Claude Code. Tokens stay in the host. GitHub is for reading the repo, reviewing diffs, and opening pull requests. Do not push main yourself.
6. Before you write any file, print Load: and Skip: from the routing table in AGENTS.md. A skill I name wins. A missing tool stops. Do not install it.
7. Keep products out of this repo. This clone is the method.

Limits you must repeat back once: a model name does not pick a skill. There is no study in this repository that shows Orchestra produces better pages than a raw host. On a saved one-line Go fix, Orchestra was slower and the diff was the same 183 bytes. That pair does not score a website. Showable UI waits for a contract and stills unless I say skip the lab.
```

---

## Limits

- `VERSION` is 4.0.0. There is no study in this repository that shows Orchestra produces better pages than a raw host.
- The saved one-line Go pair tied on the diff. Orchestra took longer.
- Research in the Go engine is offline unless you opt in.
- The engine enforces structure. It does not supply taste.
- Visual QA needs a running app.
- An MCP you have not signed in to is not healthy.
- GetLayers is purchase-pending. It is not connected.
- `docx`, `pptx`, `pdf`, and `xlsx` stay on the host. They are not in `skills/`.
- A binary or APK is `not ready` until Java is installed, and only for a target you own, a client target, a CTF, or open source.

This repo does not contain product briefs, live keys, or a populated project tree.

---

## Layout

```
AGENTS.md            the contract
.cursorrules         Cursor adapter
CLAUDE.md            Claude Code adapter
ARCHITECTURE.md      system design
protocols/           job rules, including the design lab
registries/          resources, graph, host stack
runtime/             Go engine
skills/              skills the hosts sync
templates/shelf/     skills kept in the repo and off the host list
templates/starters/  Next.js starter
kit/                 bootstrap, sync, doctor, rollback
docs/                setup and figures
VERSION              4.0.0
```

---

## Contributing

Useful changes: install fixes, a registry row with a license, a test around the design gate or the routing check.

Leave out skill dumps, always-on UI-kit servers, secrets, personal notes, and benchmarks with no measurement. See [`CONTRIBUTING.md`](CONTRIBUTING.md).

---

## License

MIT. See [LICENSE](LICENSE). Copyright (c) 2026 Mahaveer Singh Gehlot.

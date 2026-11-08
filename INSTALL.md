# Install

Orchestra 4.6.0 is a folder of skills, registries, and a Go CLI. It does not install a global pack.

## This clone

```powershell
git clone https://github.com/mahik504/orchestra-workflow.git
cd orchestra-workflow
powershell -NoProfile -ExecutionPolicy Bypass -File kit/install-git-hooks.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync
```

`-Sync` copies skills onto Cursor, Claude Code, Antigravity, Codex, jcode, and the agents host. A skill you already edited is left in place.

On macOS or Linux, `kit/sync-ides.sh` is the dry run and `kit/sync-ides.sh --apply` writes. `kit/bootstrap.sh cursor /path/to/product` stamps adapters into a product repo.

## A product repo

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\path\to\orchestra-workflow\kit\bootstrap.ps1 -Hosts cursor -Target D:\work\my-app
```

The product gets `AGENTS.md` between orchestra markers when the file is new. An existing unmarked `AGENTS.md` is not overwritten. Copy `templates/passport/project.json` to `.orchestra/project.json` and fill the name.

## What to paste once

In a host that does not read `AGENTS.md`, paste this:

```
You are running Orchestra. AGENTS.md is the contract. One conductor. Load the skills the route names. A missing tool stops. It does not install itself. Say skip orchestra to stand down.
```

## First machine

Paste this into Cursor, Antigravity, Claude Code, Codex, Kimi, or OpenCode. It is the setup, not a second conductor.

```
You are running Orchestra. AGENTS.md in the orchestra-workflow clone is the contract.
Clone https://github.com/mahik504/orchestra-workflow.git if it is not already on disk.
Run kit/sync-ides.ps1 on Windows, or kit/sync-ides.sh on macOS or Linux, and show the plan.
After I confirm, run the apply form. That copies skills, agents, commands, and the Orchestra adapters, and it adds the core MCP servers that do not need a key.
Ask before installing Obsidian. On Windows the command is: winget install --id Obsidian.Obsidian -e
On macOS the command is: brew install --cask obsidian
After I confirm, create a private vault from workspace-template/ and open that folder in Obsidian.
Learning means Preferences.md, decisions.md, and lessons.md. Write one when I like something, hate something, or make a decision. Read the matching note at the start of the next task.
Do not run npx skills add --all. Do not commit a key. Do not install a second conductor.
```

Plugin installs, for a machine that does not use `ow-sync`, are in the README. Using the plugin and `ow-sync` together copies each skill twice.

## Uninstall

Delete the skill folders Orchestra installed. The manifest is `.ow-manifest.json` in each host skill directory. Folders that are not in that manifest are yours. Leave them.

Remove the block between `<!-- orchestra:start -->` and `<!-- orchestra:end -->` from the product's `AGENTS.md`, `CLAUDE.md`, or `.cursorrules`.

The doctor:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/orchestra-doctor.ps1
```

Add `-Private` only on the machine that has the vault. A clean machine does not need it.

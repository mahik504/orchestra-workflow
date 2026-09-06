# Rollback — Orchestra 3.3.0

Three levels, smallest first. Use the smallest one that fixes the problem.

## 1. One session — `skip orchestra`

Say **skip orchestra** in the chat. The contract stands down for that session only. Nothing on disk changes. Plain text. Not a slash command.

Use when: a single task fights the routing and you just want the raw model.

## 2. One task — `skip the lab`

Say **skip the lab**. The Design Lab gate is bypassed for that task; the rest of the contract still applies.

Use when: you already know the visual direction and want implementation now.

## 3. A bad rollout — pin the contract

Set the environment variable:

```powershell
$env:ORCHESTRA_CONTRACT = "3.2.0"    # PowerShell
```

```bash
export ORCHESTRA_CONTRACT=3.2.0      # bash
```

Hosts and the CLI read `ORCHESTRA_CONTRACT` first and fall back to `VERSION` when it is unset. Pin **3.2.0** or **3.1.0**. Public tag `v3.2.0` is `9d6900d`. 3.1 zip: `orchestra-workflow-3.1.0-a4405e0-*.zip` if you archived before the 3.2 bump.

Unset it to return to the current contract:

```powershell
Remove-Item Env:\ORCHESTRA_CONTRACT
```

## 4. Last resort — the backup branch

Every contract migration creates a branch named `backup/pre-<version>-<date>` before it edits anything.

```
git branch --list "backup/*"
git switch backup/pre-v3.1-2026-09-04
```

Never force-push over `main` to undo a rollout. Land a fix-forward commit instead.

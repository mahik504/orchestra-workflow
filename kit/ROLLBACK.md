# Rollback — Orchestra 3.1.0

Three levels, smallest first. Use the smallest one that fixes the problem.

## 1. One session — `skip orchestra`

Say **skip orchestra** in the chat. The contract stands down for that session only. Nothing on disk changes.

Use when: a single task fights the routing and you just want the raw model.

## 2. One task — `skip the lab`

Say **skip the lab**. The Design Lab gate is bypassed for that task; the rest of the contract still applies.

Use when: you already know the visual direction and want implementation now.

## 3. A bad rollout — pin the contract

Set the environment variable:

```powershell
$env:ORCHESTRA_CONTRACT = "3.0.0"    # PowerShell
```

```bash
export ORCHESTRA_CONTRACT=3.0.0      # bash
```

Hosts and the CLI read `ORCHESTRA_CONTRACT` first and fall back to `VERSION` when it is unset. Pinning an older version makes the contract files for that version authoritative without reverting the repository.

Unset it to return to the current contract:

```powershell
Remove-Item Env:\ORCHESTRA_CONTRACT
```

Use when: a contract rewrite shipped and hosts started behaving worse. This is faster and safer than a git revert of a whole workspace, and it does not touch your private notes.

## 4. Last resort — the backup branch

Every contract migration creates a branch named `backup/pre-<version>-<date>` before it edits anything.

```
git branch --list "backup/*"
git switch backup/pre-v3.1-2026-09-04
```

Never force-push over `main` to undo a rollout. Land a fix-forward commit instead.

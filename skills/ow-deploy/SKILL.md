---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-deploy
description: Ship the app the product already builds. Pick the host from the project note, run its preview, then production. No new platform, no invented URL, no global CLI login from this skill.
---

# Deploy

Use this with `ow-ship-safe` when the task is to put a build where people can open it.

`vercel-labs/agent-skills` has deploy notes and no SPDX license on the GitHub API, so this file is original. It does not copy that repo.

## Steps

1. Read the project note for the host (Vercel, Cloudflare, a VPS, Expo). If it names one, use that.
2. Confirm the build command and the output directory from the repo, not from memory.
3. Preview first. Production only after the preview URL loads and the human says to promote it.
4. Secrets stay in the host's env UI. Do not write them into the repo or the vault.
5. Record the real URL in the ship note. If the deploy did not happen, say so.

## Skip

Do not add a second host to "be safe." Do not install a global CLI with a token flag. Do not scan a third-party URL. Observability stays `ow-observability`, and it is not a global Sentry install.

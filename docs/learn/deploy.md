# Deploy

`ow-deploy` reads the host the project note already names. Preview first. Production after the preview URL loads and the human says to promote it.

Secrets stay in the host's env UI. The ship note records the URL that actually opened. `ow-ship-safe` is the defensive check. It does not scan a third-party site, and it does not install Sentry globally.

`vercel-labs/agent-skills` had no SPDX license on the GitHub API on 2026-11-08, so the deploy steps in this repo are original.

---
name: ship-safe
description: Defensive security and CI/CD baseline for the user's own apps. Use when adding auth, APIs, env vars, GitHub Actions, deploy, or when they ask to harden, scan, or ship. Not for attacking third-party systems.
---

# Ship safe

You do not know this area well. Be concrete, boring, and complete. No theater.

**Only the user's apps and repos.** If they point at a live URL they do not own, refuse.

## Before ship (minimum)

- [ ] Secrets only in env / GitHub Actions secrets, never committed
- [ ] Auth on every mutating API; IDs authorized as well as authenticated
- [ ] User input validated at the boundary; parameterized queries; no raw HTML from users
- [ ] CORS, cookies (`HttpOnly`, `Secure`, `SameSite`), and HTTPS assumed in production
- [ ] Rate limit login and other expensive routes
- [ ] Error messages do not leak stack traces to clients
- [ ] Dependency install is lockfile-based
- [ ] CI runs lint, tests, and build on every PR

## Android (Expo) extra

- [ ] Tokens in `expo-secure-store`, not plain AsyncStorage
- [ ] No API keys in the app binary; use your backend
- [ ] Auth session expires; logout actually clears storage
- [ ] Play: release signing / Play App Signing, not a debug keystore in git

Do this on every Play Store build. Strix still only on apps they own, when Docker works.

## CI/CD default (GitHub Actions)

Keep one workflow until the app needs more:

- checkout with `fetch-depth: 0` if a security scan needs diffs
- setup the real runtime (Node/Python) with cache
- `npm ci` / equivalent
- lint + test + build
- fail the job on failure

Do not add matrix pyramids, OIDC, or self-hosted runners unless they asked.

## Security scans

When they want a real scan of **their** code or **their** deployed app:

1. Use Strix skills (`penetration-testing-with-strix`, then `fix-security-vulnerabilities-with-strix`).
2. For PR gating, `ci-security-scanning-with-strix`.
3. Write findings into the vault at `projects/<slug>/security.md`.
4. Fix in code. Re-scan. Do not ship “we will fix later” on auth/injection issues.

If Docker/Strix is not installed yet, say so and give the install steps. Do not fake a pentest.

## API connections

- One client module per external API
- Timeouts and retries with jitter
- Never log tokens
- Webhook routes verify signatures

## Output

When you add security or CI, also add a short review note in the vault so Claude/Antigravity can audit it later.

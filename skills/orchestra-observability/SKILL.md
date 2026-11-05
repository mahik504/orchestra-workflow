---
name: orchestra-observability
description: Error tracking, traces, and dependency or image scans inside a product repo that is actually deployed. Not a global Sentry install. Not a scan of a site you do not own.
---

# Observability — 4.0.0

Load with `ship-safe` when a real app is going out. Source methods: Sentry, OpenTelemetry, Renovate, and Trivy. They are steps in that product repo, not global MCPs and not skills in this method repo.

## In the product repo

1. Sentry, or the tracker that repo already uses, for product errors. One tracker.
2. OpenTelemetry only when traces must stay vendor-neutral. It does not replace the product tracker.
3. Renovate or the existing update bot for dependencies. Trivy, or the existing scanner, for images and lockfiles when the repo ships a container or a deployed service.
4. Secrets stay in the host's secret store. Do not commit a DSN that is a credential. A public client key still follows that product's docs.

## Do not

- Add these packages to the Orchestra method repo.
- Scan a URL the operator does not own.
- Claim the tracker is receiving events until a test event was observed.

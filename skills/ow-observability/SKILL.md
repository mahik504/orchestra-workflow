---
name: ow-observability
description: Error tracking, traces, and dependency or image scans inside a product repo that is actually deployed. Not a global Sentry install. Not a scan of a site you do not own.
---

# Observability

Load with `ow-ship-safe` when a real app is going out. Source methods: Sentry, OpenTelemetry, Renovate, and Trivy. They are steps in that product repo, not global MCPs and not skills in this method repo.

## In the product repo

1. Sentry, or the tracker that repo already uses, for product errors. One tracker.
2. OpenTelemetry only when traces must stay vendor-neutral. It does not replace the product tracker.
3. Renovate or the existing update bot for dependencies. Trivy, or the existing scanner, for images and lockfiles when the repo ships a container or a deployed service.
4. Secrets stay in the host's secret store. Do not commit a DSN that is a credential. A public client key still follows that product's docs.

## Do not

- Add these packages to the Orchestra method repo.
- Scan a URL the operator does not own.
- Claim the tracker is receiving events until a test event was observed.

## Preconditions

- The product repo is the one being deployed, and it already has an owner.

## Procedure

1. The task's paths and acceptance line win.
2. Use the tracker that repo already has. Add a second tracker only when the contract names it.
3. Keep secrets in the host store. Send one test event before saying the tracker works.

## Output

The config path and the observed test event.

## Verify

The test event was seen in the tracker, or the result says it was not.

## Skip

Skip for the Orchestra method repo, for a URL the operator does not own, and for a global install.

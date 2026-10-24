---
name: resource-scout
description: Look up one to three resources for the current job. Use when the user names a tool, a style, or asks what we have for a task.
---

# Resource scout

Read `registries/resources.json`. Match trigger tags. Return 1 to 3 rows.

Each row includes the name, the tag, and the load command.

| Tag | Load |
| --- | --- |
| job | Read the skill |
| pack | `kit/install-pack.ps1` into this repo only |
| shelf | `kit/shelf.ps1 clone` after a security check |
| refused | Decline and cite the reason |

Do not load every match. Do not install a pack the user did not name.

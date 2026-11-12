---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-resource-scout
description: Look up one to three resources for the current job. Use when the user names a tool, a style, or asks what we have for a task.
---

# Resource scout

Read `registries/resources.json`. Match trigger tags. Return 1 to 3 rows.

Each row includes the name, the tag, and the load command.

| Tag | Load |
| --- | --- |
| job | Read the skill |
| pack | `kit/install-pack.ps1` into this repo only |
| advisory | Name the conflict. Install only if the human insists, and record it |

Do not load every match. Do not install a pack the user did not name.

## Evaluate before you add

A bookmark or a gallery is a candidate, not a capability. For each candidate write the fit evidence: the trigger it matches, the license, and whether it is maintained. `ScoreFit` in the engine refuses a score with an empty evidence string. Admission moves one step at a time: discovered, candidate, evaluated, approved. It does not jump to installed because a URL was interesting.

Public places to look, when the task is visual: Godly, Awwwards, SiteInspire, Lapa, Mobbin, Refero, Codrops, Unsplash, Pexels, Lucide, Phosphor, Fontshare, Fontsource. Name the asset and the decision. Do not dump the list into the prompt.

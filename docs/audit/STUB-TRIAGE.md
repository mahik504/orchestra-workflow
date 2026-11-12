# Stub triage

Orchestra-owned skills that do not yet have Preconditions, Procedure, Output, Verify, and Skip. `stub: true` keeps the skill check honest. Expanding one of these means adding those sections and removing the stub flag.

| Skill | Bytes before stub line | Verdict | Missing |
| --- | --- | --- | --- |
| ow-3d-engineering | 2603 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-backend-architecture | 3109 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-benchmark | 1064 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-brainstorm | 943 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-career | 676 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-coach | 827 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-component-discovery | 1104 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-conductor | 1169 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-deploy | 1118 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-design-intelligence | 1368 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-diagrams | 949 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-docs | 1417 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-focus | 803 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-forge | 1447 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-fresh-review | 943 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-golden-still | 936 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-marketing | 751 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-ml-experiment | 649 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-motion-engineering | 1457 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-observability | 1130 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-onboard | 934 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-presentation | 1178 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-reference-forensics | 2222 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-research | 1007 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-research-paper | 968 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-resource-scout | 1237 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-review-panel | 882 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-ship | 1157 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-ship-safe | 2747 | keep | ## Preconditions, ## Procedure, ## Verify, ## Skip |
| ow-showcase | 720 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify |
| ow-skill-advisor | 1617 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-skill-security-check | 687 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-store-release | 627 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-typescript | 868 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-vault | 3694 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |
| ow-video | 799 | keep | ## Preconditions, ## Procedure, ## Output, ## Verify, ## Skip |

Expanded in this pass: `ow-eval`, `ow-semgrep`, `ow-site-recon`, `ow-app-recon`, `ow-scrape`, `ow-ui-assembly`, `ow-browserize`, `ow-qa`, `ow-wave`, `ow-author`, `ow-mcp-audit`.

`ow-focus` stays. The forge route still loads it. Merging it into `ow-forge` would drop a name the table uses.

`gsap-core` stays the Orchestra policy pointer to the official GSAP skill. One motion runtime per surface.

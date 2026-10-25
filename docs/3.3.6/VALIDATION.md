# Validation for 3.3.6

Local checks on this tree. GitHub Hygiene and CI run when the maintainer pushes the commit. This file does not claim they are green before that push.

| Check | Result |
| --- | --- |
| `node runtime/tools/hygiene.js` | 0 failures. Host-stack 65. Routing 4 copies, 5 document prompts |
| `go test ./...` in `runtime` | 679 passed, 19 packages |
| `kit/orchestra-doctor.ps1` | 0 problems with `ORCHESTRA_VAULT` set |
| UI pair | One private pair. `claim` is `none`. Method: `COMPARISON-METHOD.md` |

Do not write a better-than-raw sentence from this table.

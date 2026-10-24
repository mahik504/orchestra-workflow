# Validation for 3.3.5

| Check | Pass |
| --- | --- |
| `go test ./...` in `runtime` | No new failures |
| `node runtime/tools/hygiene.js` | 0 failures |
| `node runtime/tools/check-host-stack.js` | 0 problems |
| `kit/orchestra-doctor.ps1` | 0 problems |
| 3H | Three tasks, raw host versus this tree. `claim` stays none until a repeat agrees |

Do not write a better-than-raw sentence in the README until that table has numbers.

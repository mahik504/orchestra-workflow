# Verification ledger (Unlazy-style)

Prove the job with a file/evidence ledger. Do not install a second Superpowers pack. `superpowers-planning` in the allowlist already covers planning.

## When

Any claim of `DONE / FIXED / VERIFIED / PASSED / SHIPPED`, and every `orchestra verify` / Playwright pass on showable UI.

## Ledger

For this story, list:

| Claim | Evidence (path, command, URL, shot) | Result |
| --- | --- | --- |
| App launched | URL or command output | pass / fail / skip |
| Main path clicked | which control, what happened | |
| Viewport stills | 1440 / 768 / 390 paths | |
| Console | zero errors, or the errors | |
| Contrast | ratio vs `DESIGN.md` | |
| Tests | command + exit code | |
| SOS teardown | `DESIGNER_SOS_PROTOCOL.md` Keep/Kill/Park | |

Empty cells are skips. Say skip **before** success.

## Rules

- Intention is not a row.
- Another agent’s summary is not a row.
- Code inspection is not “app launched.”
- One screenshot in chat is not three viewports.
- Specialist workers fill `templates/verify-ledger-packet.md`. The conductor integrates once.

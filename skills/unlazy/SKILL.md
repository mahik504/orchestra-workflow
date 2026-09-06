---
name: unlazy
description: Evidence-first completion. Use before claiming DONE/FIXED/VERIFIED. Points at Orchestra's verification ledger. Do not install Unlazy stop hooks.
---

# Unlazy (Orchestra wrapper)

Upstream: https://github.com/Leonxlnx/unlazy

Orchestra already has `protocols/VERIFICATION_LEDGER_PROTOCOL.md`. This skill reminds the agent that **intention is not evidence**.

## When

Before writing `DONE / FIXED / VERIFIED / PASSED / SHIPPED`. After visual or backend work that could be faked from source reading.

## How

1. Fill the ledger: command output, test result, screenshot, browser state, or git state in the **same message**.
2. Report failures and skips first.
3. Do **not** install Unlazy Claude Code stop hooks. Orchestra owns completion.

Another agent's summary is not evidence.

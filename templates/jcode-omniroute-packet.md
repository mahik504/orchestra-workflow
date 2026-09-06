# jcode + OmniRoute packet (Antigravity / operator)

jcode is a **fifth executor** under Orchestra 3.2. It is not the brain. Do not replace Cursor, Antigravity, or Claude Code. Do not claim it works without command output in this packet’s handoff.

Upstream: https://github.com/1jehuang/jcode

```
MODE: specialist
KIND: jcode-omniroute
ORCHESTRA: 3.2.0  (point jcode at orchestra-workflow/AGENTS.md)
VAULT: <ORCHESTRA_HOME>
OMNIROUTE: http://localhost:20128

JOB:
1. Confirm OmniRoute is listening on 20128 (operator starts it).
2. Install/test jcode against that OpenAI-compatible endpoint.
3. Load Orchestra 3.2 AGENTS.md as the contract (same as other hosts).
4. Run one bounded task. Return evidence (version, request that succeeded or the error).

DO NOT:
- Merge extra Claude disk skills into the public 30
- Scrape LinkedIn
- Declare VERIFIED without stdout
```

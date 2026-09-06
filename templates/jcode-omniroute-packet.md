# jcode + OmniRoute packet (Antigravity / operator)

jcode is a **fourth executor** under Orchestra 3.3. It is not the brain. Do not replace Cursor, Antigravity, or Claude Code. Do not claim it works without command output in this packet’s handoff.

Upstream: https://github.com/1jehuang/jcode

```
MODE: specialist
KIND: jcode-omniroute
ORCHESTRA: 3.3.0  (point jcode at orchestra-workflow/AGENTS.md)
VAULT: <ORCHESTRA_HOME>
OMNIROUTE: http://localhost:20128

JOB:
1. Confirm OmniRoute is listening on 20128 (operator starts it).
2. Install jcode. Set `default_provider` to that OpenAI-compatible profile once.
3. After that, `jcode` (no profile flag) is the daily command. Models: `auto/best-coding`, `auto/best-reasoning`, `auto/best-fast`, `auto/best-free`.
4. Load Orchestra 3.3 AGENTS.md as the contract (same as other hosts).
5. Run one bounded task. Return evidence (version, request that succeeded or the error).

DO NOT:
- Merge extra Claude disk skills into the public allowlist
- Scrape LinkedIn
- Declare VERIFIED without stdout
```

# jcode — auth only

jcode is a fourth **executor**, not a second brain.

1. Install jcode. Point `default_provider` at your OpenAI-compatible gateway (this operator uses OmniRoute on `http://localhost:20128/v1`). Clones use **their** endpoint.
2. After that, type `jcode` in the app folder. Models: `auto/best-coding`, `auto/best-reasoning`, `auto/best-fast`, `auto/best-free` if your gateway exposes them.
3. Skills copy: `~/.jcode/skills` via `kit/sync-ides.ps1`. Contract: this repo's `AGENTS.md`.
4. Packet: `templates/jcode-omniroute-packet.md`. Evidence (stdout) before anyone claims it works.
5. Restart the IDE once so PATH includes jcode.

Do not commit OmniRoute tokens. Do not merge extra disk skills into the public allowlist.

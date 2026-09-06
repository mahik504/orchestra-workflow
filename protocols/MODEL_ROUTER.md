# Model router

Do **not** bake vendor IDs as the only source of truth. Cursor’s dropdown still belongs to the operator. Conductor **names a class** and a current slug from WORKFLOW if one exists.

Default: **one primary** model in the conducting chat, **one consultant** packet when the bar needs it, **one verifier** (Playwright / ledger). Not a 23-option swarm of full DESIGN.md files. Combinations happen at the Design Lab (23 short cards, then one contract), not by spawning every specialist.

| Class | Use | Current mapping (2026-09, may change) |
| --- | --- | --- |
| reasoning-high | architecture, spec.md, trade-offs | Opus or GPT 5.6 thinking in **Plan** |
| frontend-design-high | UI, motion, Stitch match, 3D hero | Fable (`claude-fable-5-thinking-high`) or Opus; spawn Task if parent is Grok |
| coding-fast | glue, mechanical, tests | Grok / Gemini 3.7 Flash High / Kimi |
| coding-precise | tricky backend, migrations | Opus after spec |
| research-web | current docs, named papers, public pages | Scrapling / WebFetch; ChatGPT Go / Perplexity **packet**; Context7 for libraries |
| design-critic | SOS teardown after stills | One Antigravity or Opus packet (`templates/design-critic-packet.md`) |
| security-high | threat model, authz, Strix triage | Opus Task |
| creative | visual directions, copy tone | Fable or ChatGPT packet; Higgsfield not core |
| debugging | runtime | Cursor **Debug** + evidence |
| documentation | papers, decks | orchestra-docs + reasoning-high if academic |
| harness-optional | bounded terminal job on OmniRoute | jcode packet; evidence required |

Parent Grok **cannot** become Opus. Spawn Task + tell the operator the dropdown.

Quota dead → same vault, packet to ChatGPT Go and/or Antigravity Gemini 3.7 Flash High.

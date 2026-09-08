# jcode — human actions

jcode is a fourth **executor**, not a second conductor and not a second Brain.

Current Orchestra claim: **MCP support is not invented.** If your jcode build grows MCP later, document it then. Until then, skills + `AGENTS.md` + your provider config.

## 1. Binary and PATH (FREE NOW)

1. **Service:** jcode
2. **Purpose:** terminal executor
3. **Free/paid:** upstream binary; your model gateway may be paid
4. **Authentication:** provider key stays in **your** jcode config, never git
5. **Setup:** install jcode; put it on PATH; restart the IDE once
6. **MCP URL:** none claimed
7. **Expected state:** `jcode` runs in the app folder
8. **Verification:** stdout from a real command. No stdout = not verified
9. **Recovery:** fix PATH; do not wrap jcode in a second Orchestra

## 2. Provider / OmniRoute (AUTH REQUIRED locally)

1. **Service:** OpenAI-compatible `default_provider`
2. **Purpose:** model routing for jcode
3. **Free/paid:** depends on your gateway
4. **Authentication:** your endpoint + key
5. **Setup:** point `default_provider` at **your** gateway. Clones use their endpoint. Packet: `templates/jcode-omniroute-packet.md`
6. **MCP URL:** n/a
7. **Expected state:** OPTIONAL / host-local. Not claimed HEALTHY for every machine
8. **Verification:** a coding turn that prints stdout
9. **Recovery:** switch provider; fall back to Claude Code

## 3. Skills (FREE NOW)

1. **Service:** Orchestra allowlist
2. **Purpose:** same curated skills as other hosts (host-stack allowlist)
3. **Free/paid:** free
4. **Authentication:** none
5. **Setup:** `kit/sync-ides.ps1` copies to `~/.jcode/skills`
6. **MCP URL:** n/a
7. **Expected state:** skill dirs match `registries/host-stack.json`
8. **Verification:** list `~/.jcode/skills`
9. **Recovery:** re-run sync; do not merge extra disk skills into the public allowlist

## 4. Magic UI / 21st / Refero / Tailkit / GetLayers

Do not invent jcode MCP entries for these. Use the same honesty as other hosts: free URLs and skills now; MCP only on hosts that actually support it.

GetLayers: **PURCHASE_PENDING** everywhere.

Do not commit provider tokens.

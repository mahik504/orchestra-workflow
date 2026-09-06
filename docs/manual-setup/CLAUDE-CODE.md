# Claude Code — human actions

Global adapter: `~/.claude/CLAUDE.md` (from this repo’s `CLAUDE.md` at bootstrap). Skills: `~/.claude/skills` after `kit/sync-ides.ps1`. A random folder still reads the global adapter unless it ships a competing `CLAUDE.md`.

Do not install marketplace dumps. Preserve Claude’s native MCP commands.

## 1. Core MCP (FREE NOW / AUTH REQUIRED)

1. **Service:** Playwright, Context7, vault filesystem, Stitch (if you use them here)
2. **Purpose:** QA, docs, files, 2D screens
3. **Free/paid:** mixed
4. **Authentication:** keys in host env
5. **Setup:** Claude Code MCP config / `claude mcp add`. Token in env, never git
6. **MCP URL:** see `mcp_config.example.json`
7. **Expected state:** only the servers you added
8. **Verification:** `claude mcp list` (or the current Claude UI equivalent)
9. **Recovery:** remove and re-add the server

## 2. Magic UI (FREE NOW)

1. **Service:** Magic UI MCP
2. **Purpose:** free catalog
3. **Free/paid:** FREE MCP
4. **Authentication:** none
5. **Setup:** `npx @magicuidesign/cli@latest install claude` if that client is still listed upstream, else stdio `@magicuidesign/mcp`
6. **MCP URL:** stdio
7. **Expected state:** OPTIONAL
8. **Verification:** listed in Claude MCP
9. **Recovery:** Magic UI OSS URLs

## 3. 21st.dev (AUTH REQUIRED)

Add the official CLI/MCP only from current 21st docs. Login locally. 2 free installs/day. Do not buy Builder. Do not invent a Claude-only endpoint.

## 4. Refero (PAID BUT RESOURCE-AVAILABLE)

`claude mcp add --transport http refero https://api.refero.design/mcp` only with a Pro token in host env. Otherwise use the skill + public pages.

## 5. Tailkit (PAID BUT RESOURCE-AVAILABLE)

Do not add a Tailkit MCP without a paid license. Public pages stay references.

## 6. GetLayers (PURCHASE PENDING)

Do not `claude mcp add` GetLayers until purchase + **update**.

## 7. Serena / agent-browser (OPTIONAL)

Serena: large-repo intelligence, never always-on. agent-browser: foreign sites only. Never `--tools all`. Playwright stays for our app.

Stand-down: type **skip orchestra** in the chat.

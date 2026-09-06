# Antigravity — human actions

Open the **app repo**. Add Folder your private workspace. Paste `kit/antigravity/MASTER-PROMPT.md` once. Fill MODE, workspace path, APP ROOT. Skills: `~/.gemini/config/skills` via `kit/sync-ides.ps1`.

Antigravity MCP config is **not** Cursor’s file. Merge into the Antigravity MCP Store / `mcp_config.json` without deleting unrelated servers. Never commit that file.

Keep **science** and **data-agent-kit-plugin** off as Global.

## 1. Core MCP (FREE NOW / AUTH REQUIRED)

1. **Service:** Playwright, Context7, Stitch, vault filesystem
2. **Purpose:** our-app QA, docs, 2D screens, private files
3. **Free/paid:** Playwright/Context7 free; Stitch keyed
4. **Authentication:** Stitch key / Google as required by that server
5. **Setup:** MCP Store → Connect. Keys in the Store, never git
6. **MCP URL:** same public endpoints as `kit/antigravity/mcp_config.example.json` (placeholders)
7. **Expected state:** listed. HEALTHY only after the Store shows connected
8. **Verification:** Store status. Configured ≠ HEALTHY
9. **Recovery:** reconnect; rotate keys in the Store

## 2. Magic UI (FREE NOW)

1. **Service:** Magic UI MCP
2. **Purpose:** free component catalog
3. **Free/paid:** FREE MCP. Not Pro
4. **Authentication:** none for free MCP
5. **Setup:** add stdio `npx -y @magicuidesign/mcp@latest` if the Store does not list it. Do not invent a Cursor-only installer as the only path
6. **MCP URL:** stdio as above
7. **Expected state:** OPTIONAL
8. **Verification:** Store lists it after you add it
9. **Recovery:** use Magic UI OSS URLs

## 3. 21st.dev (AUTH REQUIRED)

1. **Service:** 21st.dev
2. **Purpose:** OPTIONAL_FREE_MCP
3. **Free/paid:** search free; 2 installs/day; no Builder purchase
4. **Authentication:** login if the Store supports the official MCP
5. **Setup:** add only if the Store or official docs list Antigravity. **Do not invent support**
6. **MCP URL:** official 21st docs only
7. **Expected state:** OPTIONAL or undocumented
8. **Verification:** if the Store has no 21st entry, leave it off
9. **Recovery:** shadcn / public catalogs

## 4. Refero (PAID BUT RESOURCE-AVAILABLE)

1. **Service:** Refero MCP
2. **Purpose:** Pro screen search
3. **Free/paid:** MCP Pro; public pages free
4. **Authentication:** Bearer if Pro
5. **Setup:** Connect in MCP Store only with a real token
6. **MCP URL:** `https://api.refero.design/mcp`
7. **Expected state:** AUTH_REQUIRED without Pro
8. **Verification:** absent MCP is correct without Pro
9. **Recovery:** public Refero + Godly / SaaSFrame

## 5. Tailkit (PAID BUT RESOURCE-AVAILABLE)

Same rule as Cursor. MCP paid. Public pages remain references. Do not Connect without a license.

## 6. GetLayers (PURCHASE PENDING)

Do not Connect. Do not add a GetLayers server. After purchase, say **update**.

## 7. Host-only servers

Firebase, Supabase, Mobbin, GitHub: Connect when that job needs them. Unauthorized is AUTH_REQUIRED, not active. Do not install Gmail on Antigravity.

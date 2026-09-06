# Cursor — human actions

Clone `orchestra-workflow`. Run `kit/bootstrap.ps1` and pick Cursor. Skills sync via `kit/sync-ides.ps1`. Do not overwrite unrelated MCP servers; merge.

## 1. Core MCP (FREE NOW / AUTH REQUIRED for keyed servers)

1. **Service:** Playwright, Context7, vault filesystem, Stitch
2. **Purpose:** our-app QA, library docs, private workspace files, 2D screens
3. **Free/paid:** Playwright/Context7 free; Stitch needs your key; vault is local
4. **Authentication:** Stitch API key in the host UI
5. **Setup:** Settings → Cursor Settings → Tools & MCP. Connect. Paste keys in the host UI, never in git
6. **MCP URL:** Stitch `https://stitch.googleapis.com/mcp` (key in host). Others are stdio npx as in `mcp_config.example.json`
7. **Expected state:** listed in the host MCP panel. HEALTHY only after the host shows them connected
8. **Verification:** host MCP panel lists them; run a trivial docs lookup / browser snapshot on a job
9. **Recovery:** reconnect the server; rotate the key in the host if it expired

## 2. Magic UI free MCP (FREE NOW)

1. **Service:** Magic UI MCP
2. **Purpose:** official free component catalog
3. **Free/paid:** FREE MCP. Pro is unpaid and not claimed
4. **Authentication:** none for the free MCP
5. **Setup:** `npx @magicuidesign/cli@latest install cursor` **or** merge `magicuidesign-mcp` from `mcp_config.example.json` without deleting existing servers
6. **MCP URL:** stdio `npx -y @magicuidesign/mcp@latest`
7. **Expected state:** OPTIONAL, connected when named. Never always-on
8. **Verification:** host lists `magicuidesign-mcp`. If it does not, state EXTERNAL_AUTH_VERIFICATION_REQUIRED — here it means host-connect verification
9. **Recovery:** re-merge the example JSON; do not run a CLI that rewrites the whole file blindly

## 3. 21st.dev (FREE NOW after login; quota)

1. **Service:** 21st.dev CLI/MCP
2. **Purpose:** catalog search and limited free installs
3. **Free/paid:** search free; **2 installs/day**; Builder is paid — do not buy this pass
4. **Authentication:** browser login / API key in host. AUTH REQUIRED until login
5. **Setup:** follow https://docs.21st.dev/mcp — login locally. Do not put the API key in git
6. **MCP URL:** official 21st MCP as documented there
7. **Expected state:** OPTIONAL_FREE_MCP. Never always-on
8. **Verification:** one search on a real job after login. Do not burn the daily install quota on tests
9. **Recovery:** re-login; fall back to shadcn / Magic UI OSS

## 4. Refero (PAID BUT RESOURCE-AVAILABLE)

1. **Service:** Refero MCP vs public pages
2. **Purpose:** product UI research
3. **Free/paid:** public pages FREE NOW. MCP is Pro
4. **Authentication:** Bearer token only if you have Pro
5. **Setup:** do **not** Connect MCP without Pro. Use the `refero-design` skill + public pages
6. **MCP URL:** `https://api.refero.design/mcp`
7. **Expected state:** MCP AUTH_REQUIRED / PRO. Skill ACTIVE
8. **Verification:** skill present under `~/.cursor/skills/refero-design`. MCP absent is correct without Pro
9. **Recovery:** use Godly / Lapa / SaaSFrame

## 5. Tailkit (PAID BUT RESOURCE-AVAILABLE)

1. **Service:** Tailkit
2. **Purpose:** Tailwind patterns
3. **Free/paid:** public pages FREE NOW. MCP is paid
4. **Authentication:** paid license for MCP
5. **Setup:** do not invent credentials. Register public URLs only
6. **MCP URL:** vendor MCP if you later have a license — not this pass
7. **Expected state:** MCP AUTH_REQUIRED / PAID_ACCESS
8. **Verification:** no Tailkit MCP in the host panel is correct
9. **Recovery:** HyperUI / shadcn

## 6. GetLayers (PURCHASE PENDING)

1. **Service:** GetLayers
2. **Purpose:** future generative sections
3. **Free/paid:** unpaid
4. **Authentication:** n/a
5. **Setup:** **Do not Connect.** Do not scrape. After purchase, say **update**
6. **MCP URL:** do not add `mcp.getlayers.ai` today
7. **Expected state:** PURCHASE_PENDING / BLOCKED_UNTIL_PURCHASE
8. **Verification:** host panel must **not** list GetLayers
9. **Recovery:** Three.js / R3F / drei

## 7. Optional Cursor plugins

Gmail, Calendar, Drive, Stripe, Zapier, GitHub stay on Cursor. OAuth in the host. Do not install them on Antigravity “to match.”

Optional Cursor User Rule for chats outside this folder: `GLOBAL until I say skip orchestra.`

College VS Code is a different host.

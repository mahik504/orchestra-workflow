# Cursor — auth only

1. Clone `orchestra-workflow`. Run `kit/bootstrap.ps1` and pick Cursor (and any other hosts).
2. **Settings → Cursor Settings → Tools & MCP.** Connect the servers you actually use (Playwright, Context7, Stitch, vault filesystem). Paste keys in the host UI, never in this repo.
3. Optional plugins (Gmail, Calendar, Drive, Stripe, Zapier, shadcn) stay on Cursor. Do not install them on Antigravity "to match."
4. Refero / 21st / Firecrawl: Connect when a job names them. GetLayers: **do not Connect** while purchase is pending.
5. Optional Cursor User Rule for chats outside this folder: `GLOBAL until I say skip orchestra.`
6. Restart Cursor once so `jcode` is on PATH if you installed it.

College VS Code is a different host. Cursor User Rules do not apply there.

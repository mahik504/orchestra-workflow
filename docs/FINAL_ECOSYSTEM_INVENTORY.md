# Final Ecosystem Inventory (Orchestra V3)

## A. ACTIVE WORKING SET (Curated core loaded by agents)

**Foundation / Orchestra Core**
- `orchestra-conductor`: Master planning & orchestration
- `orchestra-vault`: Private brain memory interactions
- `orchestra-ship`: End-of-task delivery & commits
- `orchestra-docs`: Architecture and README synthesis
- `ship-safe`: Pre-deployment checks
- `superpowers-planning`: Manifesto validation
- `web-design-guidelines`: High-level UI compliance

**Visual Specialists**
- `taste-design`: Semantic typography and structure
- `emil-design-eng`: Micro-interactions and craft
- `impeccable`: Design polish
- `animate` & `review-animations`: Motion planning and critique

**Security**
- `semgrep-adapter`: Static analysis rules (added explicitly)
- `penetration-testing-with-strix`
- `fix-security-vulnerabilities-with-strix`
- `ci-security-scanning-with-strix`

**Mobile / Expo**
- `expo-router`, `expo-project-structure`, `expo-data-fetching`, `expo-native-ui`, `expo-upgrade`, `expo-dev-client`, `eas-app-stores`

**Visual Pipeline (Stitch)**
- `stitch-generate-design`, `stitch-manage-design-system`, `stitch-extract-design-md`, `stitch-extract-static-html`, `stitch-code-to-design`, `stitch-upload-to-stitch`, `stitch-react-components`

## B. CURATED OPTIONAL CATALOG (Not automatically loaded)
Located in `~/.gemini/config/curated_catalog`
- `r3f-threejs`: React Three Fiber, loaded only upon WebGL request.
- `shader-gradient`: WebGL shader liquid effects.

## C. QUARANTINED ITEMS
- The `~/.gemini/config/skills_library` (1,598 items) remains quarantined on disk. It is strictly excluded from Orchestra's context window. 

## D. DELETED ITEMS
- `github-mcp-server` (Docker connection failures; removed to prevent silent token drain).
- Unapproved Claude-only skills (e.g. `academic-paper-researcher`, `fullstack-nextjs-architect`, `hackathon-vibe-coder`, `no-ai-slop`, `permissioned-github`) deleted from `.claude/skills` to guarantee 1:1 cross-agent parity.

## E. PLUGINS
- `chrome-devtools-plugin`, `data-agent-kit-plugin`, `firebase`, `gemini-api`, `google-antigravity-sdk`, `modern-web-guidance-plugin`.

## F. MCP SERVERS (Approved)
- `StitchMCP`
- `context7`
- `orchestra-brain`
- `playwright`

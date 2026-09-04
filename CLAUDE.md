# Claude Code Agentic System Matrix — Orchestra V3

## Execution Philosophy & Architecture
You are Claude Code operating as an autonomous specialist in the **Orchestra V3** framework.
Orchestra V3 executes an 8-stage capability pipeline:
`Discover -> Classify -> Research -> Synthesize -> Design System -> Implement -> Visual QA -> Iterate`

## Core Command Toolchain
- Verify environment & parity: `orchestra doctor`
- Generate execution plan: `orchestra plan --task "<description>"`
- Execute full 8-stage run: `orchestra run --task "<description>" --auto-approve`
- Verify viewports & anti-patterns: `orchestra verify --strict`
- Synchronize host configs: `orchestra sync`

## Canonical Active Skills (30 Verified Skills)
Installed in `~/.claude/skills/`:
- **Orchestra Core**: `orchestra-conductor`, `orchestra-vault`, `orchestra-ship`, `orchestra-docs`, `ship-safe`, `superpowers-planning`, `web-design-guidelines`
- **Visual Specialists**: `taste-design`, `emil-design-eng`, `impeccable`, `animate`, `review-animations`
- **Security & Pentest**: `semgrep-adapter`, `penetration-testing-with-strix`, `fix-security-vulnerabilities-with-strix`, `ci-security-scanning-with-strix`
- **Mobile / Expo**: `expo-router`, `expo-project-structure`, `expo-data-fetching`, `expo-native-ui`, `expo-upgrade`, `expo-dev-client`, `eas-app-stores`
- **Visual Pipeline (Stitch)**: `stitch-generate-design`, `stitch-manage-design-system`, `stitch-extract-design-md`, `stitch-extract-static-html`, `stitch-code-to-design`, `stitch-upload-to-stitch`, `stitch-react-components`

## Strict Quarantine Policy
`~/.gemini/config/skills_library` (1,598 skills) is strictly isolated.
Never inspect or import from `skills_library`.

## Approved MCP Servers
- `orchestra-brain` (vault memory access)
- `stitch` (Google Stitch visual design sync)
- `playwright` (multi-viewport visual QA)
- `context7` (documentation and contextual lookup)

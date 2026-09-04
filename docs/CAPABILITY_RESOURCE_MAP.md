# Capability to Resource Map

This map defines how capability gaps are satisfied by specific curated resources.

| Capability Required | Primary Resource | Fallback / Context |
|---|---|---|
| **Frontend Architecture** | `web-design-guidelines` | React/Next.js docs |
| **Visual Design Polish** | `taste-design` + `impeccable` | `emil-design-eng` for interaction |
| **Visual QA / Browser Auto** | `playwright` (MCP) | `chrome-devtools-plugin` |
| **Security Validation** | `semgrep-adapter` | `strix` skills |
| **3D / WebGL** | `r3f` + `drei` | `shadergradient` |
| **Memory Extraction** | `orchestra-vault` | `vault-memory` (MCP) |
| **Figma / Design Sync** | `stitch-code-to-design` | `StitchMCP` |
| **Mobile Development** | `expo-router` + `expo-native-ui` | `eas-app-stores` |
| **Task Planning** | `superpowers-planning` | `orchestra-conductor` |

*Note: If a required capability is not listed here, Orchestra will invoke gap-analysis and search the `curated_catalog`, then quarantine, or perform live web research rather than blindly guessing.*

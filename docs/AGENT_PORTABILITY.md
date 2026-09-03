# Agent Portability Matrix

Orchestra V3 operates as an abstracted capability router. The capability set is verified to be identical across supported environments.

| Environment | Skill Dir | MCP Config | Adapter Role | Status |
|---|---|---|---|---|
| **Cursor** | `.cursor/skills` | `.cursor/mcp.json` | Implementation, component tuning | `SYNCED` |
| **Antigravity** | `.gemini/config/skills` | `.gemini/config/mcp_config.json` | Master orchestration, QA, Architecture | `SYNCED` |
| **Claude Code** | `.claude/skills` | `.claude/mcp.json` (via env) | Terminal automation, server-side debugging | `SYNCED` |
| **OpenCode/Agents** | `.agents/skills` | N/A | Batch tasks | `SYNCED` |

## Enforcement Policy
The `skills` directories have been programmatically truncated. Unauthorized generic AI skills (`hackathon-vibe-coder`, `no-ai-slop`, etc.) were eradicated to guarantee a 1:1 consistent Execution Manifest across all adapters.

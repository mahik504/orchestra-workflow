# Control plane

Orchestra is the only control plane. Skills, MCPs, plugins, and libraries are capabilities. Agents are executors. The private Brain is memory. The registry is resource knowledge.

```mermaid
flowchart TD
  human[Human]
  cond[One conductor this session]
  graph[Capability graph]
  reg[Resource registry]
  skills[Curated skills]
  mcp[MCP optional or auth]
  exec[Executors: Cursor / Antigravity / Claude Code / jcode]
  brain[Private Brain]
  human --> cond
  cond --> graph
  cond --> reg
  cond --> skills
  cond --> mcp
  cond --> exec
  cond --> brain
```

Available ≠ loaded. Registered ≠ active. Discovered ≠ verified. Another agent's summary ≠ evidence.

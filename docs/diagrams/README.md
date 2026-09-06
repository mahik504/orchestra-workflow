# Diagrams

README shows **three** committed SVGs in [`assets/`](assets/). GitHub’s Mermaid parser treats ids such as `graph` and `click` as keywords, which produced the red “Unable to render” boxes on the old README.

Mermaid sources for those three figures sit next to the SVGs (`.mmd`). Re-render:

```bash
node docs/diagrams/assets/render-readme-diagrams.mjs
```

The eight notes below are the architecture set for editors. Do not paste all of them into README.

| File | What it shows |
| --- | --- |
| [01-control-plane.md](01-control-plane.md) | Human → one conductor → graph / skills / executors / memory |
| [02-capability-graph.md](02-capability-graph.md) | Brief → one route, or research then wait for **update** |
| [03-resource-lifecycle.md](03-resource-lifecycle.md) | discovered → selected → acquired → used → verified |
| [04-design-lab.md](04-design-lab.md) | 23 cards → one DESIGN.md → human gate |
| [05-research-layer.md](05-research-layer.md) | Public galleries + Scrapling. No fake MCP |
| [06-four-hosts.md](06-four-hosts.md) | Allowlist sync to four hosts |
| [07-verification-loop.md](07-verification-loop.md) | Launch, exercise, stills, SOS, ledger |
| [08-resource-combination.md](08-resource-combination.md) | One DESIGN.md. Not a collage |

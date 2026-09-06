# Contributing

This is a **methodology template**, not an application. You clone the method. Your products stay in *your* private workspace.

## Welcome

- Fixes to `kit/bootstrap.ps1` / `kit/bootstrap.sh` and host adapters
- Clearer docs (especially MCP honesty and Design Lab)
- Registry rows in `registries/resources.json` with license + a maintenance signal (last commit, docs URL)
- Graph rows only when a real job needed a new route — trigger **and** skip, no defaults in disguise
- Protocol clarifications that do not add skill dumps
- Go tests for classify, the Design Lab write-lock, and `AUTH_REQUIRED` vs `HEALTHY`

Open an issue first for a graph miss or a resource documented as free when the MCP is paid.

## How to check a change

```bash
cd runtime && go test ./...
node runtime/tools/hygiene.js
node runtime/tools/check-registry.js
node runtime/tools/check-host-stack.js
```

Do not commit `.env`, live `mcp.json`, or a populated `projects/` tree. Hygiene CI will fail the build.

## Not welcome

- Skill packs (`--all`, ECC catalogs, marketplace dumps, “install 88”)
- Always-on UI kit MCP
- Screenshot-to-code factories
- Personal brains, product briefs, career notes
- Secrets, `.env`, MCP JSON with keys
- Fake benchmarks (“X% faster”, “zero errors”) without a described measurement
- Contribution-graph painting

## License

MIT. See [LICENSE](./LICENSE).

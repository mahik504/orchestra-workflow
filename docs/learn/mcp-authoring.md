# Writing an MCP contract

Add a row to `registries/mcp-contracts.json`. Required text:

- `invoke_when` names the task that may call it.
- `skip_when` names the task that must not.
- `install` includes "Do not install globally".
- `usage_evidence` stays "none recorded in this registry" until a run exists.
- `permission_profile` is one of `project-write`, `keyed`, `read-only-until-the-task-names-a-write`, `local-authorized`, `optional`.

A key lives in the host. The row does not contain the key. A template for a keyed server is `mcp_config.example.json` with empty values.

"""Merge Playwright, Context7, and vault filesystem into Antigravity MCP config.

Never prints secret values. Never overwrites existing servers.
"""
from __future__ import annotations

import argparse
import json
from pathlib import Path


ADD = {
    "orchestra-brain": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-filesystem", "__VAULT__"],
    },
    "playwright": {
        "command": "npx",
        "args": ["-y", "@playwright/mcp@latest"],
    },
    "context7": {
        "command": "npx",
        "args": ["-y", "@upstash/context7-mcp"],
    },
}


def main() -> None:
    p = argparse.ArgumentParser()
    p.add_argument("--mcp", required=True)
    p.add_argument("--vault", required=True)
    args = p.parse_args()
    path = Path(args.mcp)
    data = {"mcpServers": {}}
    if path.exists():
        data = json.loads(path.read_text(encoding="utf-8"))
    servers = data.setdefault("mcpServers", {})
    added = []
    skipped = []
    for name, spec in ADD.items():
        if name in servers:
            skipped.append(name)
            continue
        # Also skip if they already have a stitch-like or filesystem vault under another name
        if name == "orchestra-brain" and any(
            k.lower() in ("orchestra-brain", "vault", "filesystem") for k in servers
        ):
            skipped.append(name + "(similar exists)")
            continue
        entry = json.loads(json.dumps(spec))
        if name == "orchestra-brain":
            entry["args"][-1] = args.vault
        servers[name] = entry
        added.append(name)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2) + "\n", encoding="utf-8")
    print("mcp servers now:", ", ".join(sorted(servers.keys())))
    print("added:", ", ".join(added) or "(none)")
    print("kept:", ", ".join(skipped) or "(none)")


if __name__ == "__main__":
    main()

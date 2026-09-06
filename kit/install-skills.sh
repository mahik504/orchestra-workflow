#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SKILLS="$ROOT/skills"
STACK="$ROOT/registries/host-stack.json"
if [[ ! -d "$SKILLS" ]]; then
  echo "missing skills/" >&2
  exit 1
fi
if [[ ! -f "$STACK" ]]; then
  echo "missing registries/host-stack.json" >&2
  exit 1
fi
DESTS=(
  "$HOME/.cursor/skills"
  "$HOME/.claude/skills"
  "$HOME/.agents/skills"
  "$HOME/.gemini/config/skills"
  "$HOME/.jcode/skills"
)
node -e '
  const s=require(process.argv[1]);
  const fs=require("fs"); const path=require("path");
  const root=process.argv[2];
  const dests=process.argv.slice(3);
  for (const n of s.skills) {
    const src=path.join(root,"skills",n);
    if (!fs.existsSync(src)) { console.error("missing skill "+n); process.exit(1); }
    for (const d of dests) {
      if (!fs.existsSync(d)) continue;
      const dst=path.join(d,n);
      fs.rmSync(dst,{recursive:true,force:true});
      fs.cpSync(src,dst,{recursive:true});
      console.log("installed "+n+" -> "+dst);
    }
  }
  console.log("Done. Orchestra "+s.version+". Restart the agent.");
' "$STACK" "$ROOT" "${DESTS[@]}"

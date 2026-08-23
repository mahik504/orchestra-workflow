#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SKILLS="$ROOT/skills"
if [[ ! -d "$SKILLS" ]]; then
  echo "missing skills/" >&2
  exit 1
fi
DESTS=(
  "$HOME/.cursor/skills"
  "$HOME/.claude/skills"
  "$HOME/.agents/skills"
  "$HOME/.gemini/config/skills"
)
for src in "$SKILLS"/*; do
  [[ -d "$src" ]] || continue
  name="$(basename "$src")"
  for d in "${DESTS[@]}"; do
    [[ -d "$d" ]] || continue
    mkdir -p "$d/$name"
    cp -R "$src/." "$d/$name/"
    echo "installed $name -> $d/$name"
  done
done
echo "Done. Restart the agent."

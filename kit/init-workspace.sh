#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TARGET="${ORCHESTRA_HOME:-$(cd "$ROOT/.." && pwd)/orchestra-workspace}"
TEMPLATE="$ROOT/workspace-template"
if [[ ! -d "$TEMPLATE" ]]; then
  echo "missing workspace-template" >&2
  exit 1
fi
mkdir -p "$TARGET"
cp -R "$TEMPLATE/." "$TARGET/"
for dir in protocols registries templates docs; do
  if [[ -d "$ROOT/$dir" ]]; then
    rm -rf "$TARGET/$dir"
    cp -R "$ROOT/$dir" "$TARGET/$dir"
  fi
done
for f in AGENTS.md Preferences.md WORKFLOW.md START.md; do
  if [[ -f "$ROOT/$f" ]]; then
    cp "$ROOT/$f" "$TARGET/$f"
  fi
done
echo "Private workspace: $TARGET"
echo "Point your agent at that folder + your app repo. Do not commit secrets."

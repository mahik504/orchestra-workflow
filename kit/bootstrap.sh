#!/bin/sh
# Non-interactive bootstrap. Dry skill sync is not enough here: --apply copies
# only into a conflict-safe sync. Adapters are stamped between orchestra markers.
set -e
root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
hosts="${1:-cursor}"
target="${2:-${ORCHESTRA_HOME:-$root/../orchestra-workspace}}"
mkdir -p "$target"
old_ifs=$IFS
IFS=,
for host in $hosts; do
  node "$root/runtime/tools/ow-sync.mjs" --apply --host "$host"
done
IFS=$old_ifs
node "$root/runtime/tools/stamp-adapter.mjs" --dest "$target/AGENTS.md" --source "$root/AGENTS.md"
node "$root/runtime/tools/stamp-adapter.mjs" --dest "$target/CLAUDE.md" --source "$root/CLAUDE.md"
node "$root/runtime/tools/stamp-adapter.mjs" --dest "$target/.cursorrules" --source "$root/.cursorrules"
echo "Orchestra $(tr -d '\r\n' < "$root/VERSION") stamped into $target"

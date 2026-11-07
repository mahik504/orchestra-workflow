#!/bin/sh
# Sync canonical skills through runtime/tools/ow-sync.mjs. Dry-run unless --apply.
# --prune is refused.
set -e
root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
apply=""
replace=""
for arg in "$@"; do
  case "$arg" in
    --apply) apply="--apply" ;;
    --backup-and-replace) replace="--backup-and-replace" ;;
    --prune)
      echo "Prune is refused. User-installed skills stay on the host."
      exit 2
      ;;
    *)
      echo "unknown argument $arg" >&2
      exit 2
      ;;
  esac
done
node "$root/runtime/tools/ow-sync.mjs" $apply $replace

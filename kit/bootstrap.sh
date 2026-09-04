#!/usr/bin/env bash
# Wire Orchestra 3.1.0 onto chosen hosts. No secrets. No skills add --all.
#
#   ./kit/bootstrap.sh
#   ./kit/bootstrap.sh --hosts cursor,antigravity --target /path/to/workspace
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
HOSTS=""
TARGET="${ORCHESTRA_HOME:-}"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --hosts) HOSTS="$2"; shift 2 ;;
    --target) TARGET="$2"; shift 2 ;;
    *) echo "unknown arg: $1" >&2; exit 1 ;;
  esac
done

STACK="$ROOT/registries/host-stack.json"
if [[ ! -f "$STACK" ]]; then
  echo "missing registries/host-stack.json" >&2
  exit 1
fi

pick_hosts() {
  if [[ -n "$HOSTS" ]]; then
    echo "$HOSTS" | tr ',' ' '
    return
  fi
  echo "Which hosts should use Orchestra 3.1.0?" >&2
  echo "  [1] Cursor" >&2
  echo "  [2] Antigravity" >&2
  echo "  [3] Claude Code" >&2
  echo "  [4] Codex / Hermes / OpenCode (AGENTS.md only)" >&2
  read -r -p "Enter numbers (e.g. 1 2) " ans
  out=()
  for tok in $ans; do
    case "$tok" in
      1) out+=(cursor) ;;
      2) out+=(antigravity) ;;
      3) out+=(claude) ;;
      4) out+=(agents) ;;
    esac
  done
  if [[ ${#out[@]} -eq 0 ]]; then
    echo "No hosts selected." >&2
    exit 1
  fi
  echo "${out[*]}"
}

CHOSEN=($(pick_hosts))
echo "Hosts: ${CHOSEN[*]}"

if [[ -z "$TARGET" ]]; then
  TARGET="${ORCHESTRA_HOME:-$(cd "$ROOT/.." && pwd)/orchestra-workspace}"
fi
if [[ ! -d "$TARGET" ]] || [[ -z "$(ls -A "$TARGET" 2>/dev/null)" ]]; then
  ORCHESTRA_HOME="$TARGET" "$ROOT/kit/init-workspace.sh"
else
  echo "workspace already exists, leaving $TARGET untouched"
fi

skill_dest_for() {
  case "$1" in
    cursor) echo "$HOME/.cursor/skills" ;;
    antigravity) echo "$HOME/.gemini/config/skills" ;;
    claude) echo "$HOME/.claude/skills" ;;
    agents) echo "$HOME/.agents/skills" ;;
  esac
}

copy_skills() {
  local dest="$1"
  mkdir -p "$dest"
  node -e "
    const s=require(process.argv[1]);
    const fs=require('fs'); const path=require('path');
    const root=process.argv[2]; const dest=process.argv[3];
    for (const n of s.skills) {
      const src=path.join(root,'skills',n);
      const dst=path.join(dest,n);
      fs.rmSync(dst,{recursive:true,force:true});
      fs.cpSync(src,dst,{recursive:true});
    }
    console.log('copied '+s.skills.length+' skills -> '+dest);
  " "$STACK" "$ROOT" "$dest"
}

for h in "${CHOSEN[@]}"; do
  d="$(skill_dest_for "$h")"
  if [[ -n "$d" ]]; then copy_skills "$d"; fi
done

mkdir -p "$TARGET/kit"
for h in "${CHOSEN[@]}"; do
  case "$h" in
    cursor)
      cp "$ROOT/.cursorrules" "$TARGET/.cursorrules"
      echo "wrote workspace .cursorrules"
      ;;
    claude)
      cp "$ROOT/CLAUDE.md" "$TARGET/CLAUDE.md"
      echo "wrote workspace CLAUDE.md"
      ;;
    antigravity)
      mkdir -p "$TARGET/kit/antigravity"
      cp "$ROOT/kit/antigravity/MASTER-PROMPT.md" "$TARGET/kit/antigravity/"
      cp "$ROOT/kit/antigravity/ALWAYS-ON.md" "$TARGET/kit/antigravity/"
      cp "$ROOT/kit/antigravity/mcp_config.example.json" "$TARGET/kit/antigravity/"
      mkdir -p "$HOME/.gemini"
      if [[ -f "$HOME/.gemini/mcp_config.json" ]]; then
        echo "left existing $HOME/.gemini/mcp_config.json untouched (may contain keys)"
      else
        cp "$ROOT/kit/antigravity/mcp_config.example.json" "$HOME/.gemini/mcp_config.example.json"
        echo "wrote $HOME/.gemini/mcp_config.example.json (template only)"
      fi
      ;;
  esac
done

cp "$ROOT/mcp_config.example.json" "$TARGET/mcp_config.example.json"
cp "$ROOT/AGENTS.md" "$TARGET/AGENTS.md"

echo
echo "Marketplace plugins (install yourself; we cannot log into Google, Stripe, or Stitch for you):"
node -e "
  const s=require(process.argv[1]);
  const chosen=process.argv.slice(2);
  if (chosen.includes('cursor')) {
    console.log('  Cursor:');
    for (const p of s.plugins.cursor) console.log('    - '+p);
  }
  if (chosen.includes('antigravity')) {
    console.log('  Antigravity:');
    for (const p of s.plugins.antigravity) console.log('    - '+p);
  }
" "$STACK" "${CHOSEN[@]}"

echo
echo "Restart the IDE. Next chat in this workspace uses Orchestra 3.1.0 (AGENTS.md)."
echo "Optional Cursor User Rule for chats outside this folder: GLOBAL until I say skip orchestra."
echo "We cannot log into Google, Stripe, or Stitch for you."
echo "Private workspace: $TARGET"

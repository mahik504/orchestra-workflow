#!/usr/bin/env bash
# Persist ORCHESTRA_HOME in the user shell profile and install orchestra onto GOPATH/bin.
# Pass HOME_DIR as the first argument, or set ORCHESTRA_HOME first.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
HOME_DIR="${1:-${ORCHESTRA_HOME:-}}"
if [ -z "$HOME_DIR" ]; then
  echo "Set ORCHESTRA_HOME or pass the private workspace path as argument 1." >&2
  exit 1
fi
if [ ! -d "$HOME_DIR" ]; then
  echo "Home dir does not exist: $HOME_DIR" >&2
  exit 1
fi

export GOBIN="${GOBIN:-$HOME/go/bin}"
mkdir -p "$GOBIN"
(cd "$ROOT/runtime" && go install ./cmd/orchestra)

PROFILE=""
if [ -f "$HOME/.zshrc" ]; then PROFILE="$HOME/.zshrc"
elif [ -f "$HOME/.bashrc" ]; then PROFILE="$HOME/.bashrc"
else PROFILE="$HOME/.profile"
fi
touch "$PROFILE"
if ! grep -q 'export ORCHESTRA_HOME=' "$PROFILE"; then
  {
    echo ""
    echo "# Orchestra 3.1.0"
    echo "export ORCHESTRA_HOME=\"$HOME_DIR\""
    echo "export PATH=\"\$PATH:$GOBIN\""
  } >> "$PROFILE"
  echo "Appended ORCHESTRA_HOME to $PROFILE"
else
  echo "ORCHESTRA_HOME already set in $PROFILE"
fi
echo "installed $GOBIN/orchestra"
echo "Reload the shell, then: orchestra doctor"

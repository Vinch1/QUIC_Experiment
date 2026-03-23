#!/usr/bin/env bash

set -euo pipefail

echo "==> Environment check"

check() {
  local name="$1"
  local cmd="$2"
  if command -v "$cmd" >/dev/null 2>&1; then
    echo "[ok] ${name}: $(command -v "$cmd")"
  else
    echo "[missing] ${name}: command '${cmd}' not found"
  fi
}

check "go" "go"
check "docker" "docker"
check "protoc" "protoc"
check "make" "make"

echo
echo "Recommended minimum stack:"
echo "- Go 1.26+"
echo "- Docker Desktop"
echo "- protoc"
echo "- GNU make"


#!/usr/bin/env bash
# Build production artifacts for both apps.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> Building API"
mkdir -p "$root/apps/api/bin"
(cd "$root/apps/api" && go build -o bin/server ./cmd/server)

echo "==> Building web"
(cd "$root/apps/web" && npm run build)

echo "Artifacts:"
echo "  apps/api/bin/server"
echo "  apps/web/dist/"

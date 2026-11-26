#!/usr/bin/env bash
# Install dependencies for both apps.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> Downloading Go modules"
(cd "$root/apps/api" && go mod download)

echo "==> Installing web dependencies"
(cd "$root/apps/web" && npm install)

echo "Done."

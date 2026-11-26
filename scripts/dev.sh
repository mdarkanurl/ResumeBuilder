#!/usr/bin/env bash
# Run the API and the web dev server together.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -f "$root/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$root/.env"
  set +a
fi

if [[ ! -d "$root/apps/web/node_modules" ]]; then
  echo "Web dependencies are missing. Run scripts/setup.sh first." >&2
  exit 1
fi

pids=()
cleanup() {
  for pid in "${pids[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
}
trap cleanup EXIT INT TERM

echo "==> API  http://localhost:${PORT:-8080}"
(cd "$root/apps/api" && go run ./cmd/server) &
pids+=("$!")

echo "==> Web  http://localhost:5173"
(cd "$root/apps/web" && npm run dev) &
pids+=("$!")

wait

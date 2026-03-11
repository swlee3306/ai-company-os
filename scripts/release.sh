#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="$ROOT/release"

mkdir -p "$OUT"

COMMIT="$(git -C "$ROOT" rev-parse --short HEAD 2>/dev/null || echo unknown)"
BUILDTIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
VERSION="0.1.0"

echo "[release] build backend"
(
  cd "$ROOT"
  go build -o "$OUT/company" \
    -ldflags "-X 'github.com/swlee3306/ai-company-os/internal/api.Version=$VERSION' -X 'github.com/swlee3306/ai-company-os/internal/api.Commit=$COMMIT' -X 'github.com/swlee3306/ai-company-os/internal/api.BuildTime=$BUILDTIME'" \
    ./cmd/company
)

echo "[release] build web"
(
  cd "$ROOT/web"
  npm ci
  npm run build
  rm -rf "$OUT/web"
  cp -R dist "$OUT/web"
)

echo "[release] done"
echo "- backend: $OUT/company"
echo "- web:     $OUT/web"

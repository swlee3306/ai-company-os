#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PEN="$ROOT_DIR/design/company-ui.pen"
EXPORT_DIR="$ROOT_DIR/design/exports"

required=(
  "install-onboarding.png"
  "operations-dashboard.png"
  "task-board-detail.png"
  "logs-audit.png"
  "diagnostics-settings.png"
)

if [ ! -f "$PEN" ]; then
  echo "ERROR: missing $PEN" >&2
  exit 1
fi

pen_mtime=$(stat -c %Y "$PEN")

missing=0
stale=0

for f in "${required[@]}"; do
  p="$EXPORT_DIR/$f"
  if [ ! -f "$p" ]; then
    echo "MISSING: $p"
    missing=1
    continue
  fi
  m=$(stat -c %Y "$p")
  if [ "$m" -lt "$pen_mtime" ]; then
    echo "STALE:   $f (export is older than company-ui.pen)"
    stale=1
  else
    echo "OK:      $f"
  fi
done

if [ "$missing" -eq 1 ]; then
  echo "\nOne or more required exports are missing." >&2
  exit 2
fi

if [ "$stale" -eq 1 ]; then
  echo "\nOne or more exports are stale. Re-export from Pencil." >&2
  exit 3
fi

echo "\nAll required exports present and up to date."

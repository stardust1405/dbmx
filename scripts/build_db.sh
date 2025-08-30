#!/usr/bin/env bash
set -euo pipefail

echo "[build_db] shell: $0"
echo "[build_db] whoami: $(whoami)"
echo "[build_db] pwd:    $(pwd)"

# Resolve repo root from THIS script's location (works no matter where hook runs)
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
MIG="$ROOT/migrations"
OUT_DIR="$ROOT/config/database"
OUT_DB="$OUT_DIR/app.db"

echo "[build_db] ROOT:    $ROOT"
echo "[build_db] MIG:     $MIG"
echo "[build_db] OUT_DIR: $OUT_DIR"
echo "[build_db] OUT_DB:  $OUT_DB"

mkdir -p "$OUT_DIR"

# Prefer installed goose; otherwise fall back to 'go run' so PATH issues don't kill the build
if command -v goose >/dev/null 2>&1; then
  GOOSE="goose"
else
  echo "[build_db] goose not found in PATH, using 'go run' fallback"
  GOOSE="go run github.com/pressly/goose/v3/cmd/goose@latest"
fi

# Show versions for diagnostics
set -x
go version
$GOOSE -version || true
set +x

# Clean and build
rm -f "$OUT_DB"
set -x
$GOOSE -dir "$MIG" sqlite3 "$OUT_DB" up
set +x

ls -lh "$OUT_DB"
echo "[build_db] DONE"

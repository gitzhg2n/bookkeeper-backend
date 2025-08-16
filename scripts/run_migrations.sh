#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/internal/db/migrations"
DB_FILE="${DB_FILE:-$ROOT_DIR/bookkeeper.db}"

echo "Running migrations from: $MIGRATIONS_DIR"
if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "No migrations directory found at $MIGRATIONS_DIR. Nothing to do."
  exit 0
fi

shopt -s nullglob
SQL_FILES=("$MIGRATIONS_DIR"/*.sql)
if [ ${#SQL_FILES[@]} -eq 0 ]; then
  echo "No .sql migration files found in $MIGRATIONS_DIR"
  exit 0
fi

# If sqlite3 CLI is available, apply .sql files against DB_FILE.
if command -v sqlite3 >/dev/null 2>&1; then
  echo "Found sqlite3. Applying migrations to $DB_FILE"
  for f in "${SQL_FILES[@]}"; do
    echo "Applying $f"
    sqlite3 "$DB_FILE" < "$f"
  done
  echo "Migrations applied."
else
  echo "sqlite3 CLI not found on runner. Listing migration files instead."
  for f in "${SQL_FILES[@]}"; do
    echo "$f"
  done
  echo "To actually apply migrations in CI, install sqlite3 or provide a migration runner."
fi

exit 0

#!/usr/bin/env bash
set -euo pipefail

# Usage: scripts/seed_db.sh [container_name]
# If container_name is provided, the script will run the seed SQL inside that container using psql.

CONTAINER_NAME=${1:-bookkeeper-backend_db_1}
SEED_FILE="$(pwd)/internal/db/seeds/seed_test_data.sql"

if docker ps --format '{{.Names}}' | grep -q "${CONTAINER_NAME}"; then
  echo "Seeding DB inside container: $CONTAINER_NAME"
  cat "$SEED_FILE" | docker exec -i "$CONTAINER_NAME" psql -U bookkeeper -d bookkeeper
  echo "Seed applied inside container."
else
  echo "Container $CONTAINER_NAME not found. Attempting to run psql locally if available..."
  if command -v psql >/dev/null 2>&1; then
    PSQL_CMD=${PSQL_CMD:-psql}
    $PSQL_CMD -h ${DB_HOST:-localhost} -U ${DB_USER:-bookkeeper} -d ${DB_NAME:-bookkeeper} -f "$SEED_FILE"
    echo "Seed applied via local psql."
  else
    echo "No container found and local psql not available. Skipping seed.";
  fi
fi

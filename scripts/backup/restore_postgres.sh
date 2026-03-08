#!/usr/bin/env bash
# NEMORIS — PostgreSQL restore script
# Restores from backup file via psql (docker exec).
# Usage: ./restore_postgres.sh [--reset] <backup_filename>
#   --reset  DROP SCHEMA public CASCADE; CREATE SCHEMA public; before restore
# Example: ./restore_postgres.sh nemoris_20260308_120000.sql
# Example: ./restore_postgres.sh --reset nemoris_20260308_120000.sql

set -e

CONTAINER="nemoris-postgres"
BACKUP_DIR="backups/postgres"
RESET_SCHEMA=""
BACKUP_FILENAME=""

# Parse args: --reset and filename
for arg in "$@"; do
  if [ "${arg}" = "--reset" ]; then
    RESET_SCHEMA=1
  elif [ -z "${BACKUP_FILENAME}" ]; then
    BACKUP_FILENAME="${arg}"
  fi
done

# Resolve script dir and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
cd "${PROJECT_ROOT}"

# Fail clearly if filename missing
if [ -z "${BACKUP_FILENAME}" ]; then
  echo "ERROR: Backup filename required."
  echo "Usage: $0 [--reset] <backup_filename>"
  echo "Example: $0 nemoris_20260308_120000.sql"
  echo "Example: $0 --reset nemoris_20260308_120000.sql"
  echo ""
  echo "Available backups in ${BACKUP_DIR}:"
  ls -la "${BACKUP_DIR}" 2>/dev/null || echo "  (none)"
  exit 1
fi

BACKUP_FILE="${BACKUP_DIR}/${BACKUP_FILENAME}"
if [ ! -f "${BACKUP_FILE}" ]; then
  echo "ERROR: Backup file not found: ${BACKUP_FILE}"
  exit 1
fi

# Fail clearly if container missing
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
  echo "ERROR: Container '${CONTAINER}' is not running. Start infra first (make infra-dev or make infra-prod)."
  exit 1
fi

# Load DB credentials from .env if present
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi
DB_USER="${DB_USER:-nemoris}"
DB_NAME="${DB_NAME:-nemoris}"

echo "WARNING: This will overwrite database '${DB_NAME}'. Press Ctrl+C to cancel, or Enter to continue."
read -r

if [ -n "${RESET_SCHEMA}" ]; then
  echo "Resetting schema (DROP SCHEMA public CASCADE; CREATE SCHEMA public)..."
  docker exec -i "${CONTAINER}" psql -U "${DB_USER}" -d "${DB_NAME}" <<'SQL'
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
SQL
fi

echo "Restoring from ${BACKUP_FILE}..."
docker exec -i "${CONTAINER}" psql -U "${DB_USER}" -d "${DB_NAME}" < "${BACKUP_FILE}"

echo "OK: Restore completed."

#!/usr/bin/env bash
# NEMORIS — PostgreSQL backup script
# Uses pg_dump via docker exec. Output: backups/postgres/nemoris_YYYYMMDD_HHMMSS.sql
# Run from project root.

set -e

CONTAINER="nemoris-postgres"
BACKUP_DIR="backups/postgres"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${BACKUP_DIR}/nemoris_${TIMESTAMP}.sql"

# Resolve script dir and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
cd "${PROJECT_ROOT}"

# Create backup dir if missing
mkdir -p "${BACKUP_DIR}"

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

echo "Backing up PostgreSQL to ${OUTPUT_FILE}..."
docker exec "${CONTAINER}" pg_dump -U "${DB_USER}" -d "${DB_NAME}" --no-owner --no-acl > "${OUTPUT_FILE}"

if [ -s "${OUTPUT_FILE}" ]; then
  echo "OK: Backup saved to ${OUTPUT_FILE}"
else
  echo "ERROR: Backup file is empty. Aborting."
  rm -f "${OUTPUT_FILE}"
  exit 1
fi

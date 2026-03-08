#!/usr/bin/env bash
# NEMORIS — WAHA session backup script
# Archives WAHA session volume via docker cp + tar.gz.
# Output: backups/waha/waha_sessions_YYYYMMDD_HHMMSS.tar.gz
# Run from project root.

set -e

CONTAINER="waha"
WAHA_SESSIONS_PATH="/app/.sessions"
BACKUP_DIR="backups/waha"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${BACKUP_DIR}/waha_sessions_${TIMESTAMP}.tar.gz"
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "${TEMP_DIR}"' EXIT

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

echo "Backing up WAHA sessions to ${OUTPUT_FILE}..."
docker cp "${CONTAINER}:${WAHA_SESSIONS_PATH}/." "${TEMP_DIR}/sessions"
tar -czf "${OUTPUT_FILE}" -C "${TEMP_DIR}" sessions

if [ -f "${OUTPUT_FILE}" ] && [ -s "${OUTPUT_FILE}" ]; then
  echo "OK: Backup saved to ${OUTPUT_FILE}"
else
  echo "ERROR: Backup file is empty or missing. Aborting."
  exit 1
fi

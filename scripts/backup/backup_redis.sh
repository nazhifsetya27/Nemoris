#!/usr/bin/env bash
# NEMORIS — Redis backup script
# Copies full /data directory from Redis container via docker cp.
# Captures appendonly.aof, appendonlydir/, and any RDB files.
# Output: backups/redis/redis_data_YYYYMMDD_HHMMSS.tar.gz
# Run from project root.

set -e

CONTAINER="nemoris-redis"
REDIS_DATA_PATH="/data"
BACKUP_DIR="backups/redis"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${BACKUP_DIR}/redis_data_${TIMESTAMP}.tar.gz"
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

echo "Backing up Redis /data to ${OUTPUT_FILE}..."
docker cp "${CONTAINER}:${REDIS_DATA_PATH}/." "${TEMP_DIR}/data"
tar -czf "${OUTPUT_FILE}" -C "${TEMP_DIR}" data

if [ -f "${OUTPUT_FILE}" ] && [ -s "${OUTPUT_FILE}" ]; then
  echo "OK: Backup saved to ${OUTPUT_FILE}"
else
  echo "ERROR: Backup file is empty or missing. Aborting."
  exit 1
fi

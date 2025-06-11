#!/bin/bash
# Litestream initialization script for Skyclerk
# This script handles database restoration and starts continuous replication

set -e

echo "Starting Litestream initialization..."

# Database path from environment or default
DB_PATH="${DB_PATH:-/data/skyclerk.sqlite}"
DB_DIR=$(dirname "$DB_PATH")

# Ensure the data directory exists
mkdir -p "$DB_DIR"

# Check if we need to restore the database
if [ ! -f "$DB_PATH" ]; then
    echo "Database not found at $DB_PATH. Attempting to restore from S3..."
    
    # Attempt to restore the database
    if litestream restore -config /app/litestream.yml "$DB_PATH"; then
        echo "Database successfully restored from S3"
    else
        echo "Failed to restore database from S3. This might be the first deployment."
        echo "Creating new database..."
        # The app will create a new database on first startup
    fi
else
    echo "Database found at $DB_PATH"
fi

# Ensure proper permissions
chmod 644 "$DB_PATH" 2>/dev/null || true

echo "Litestream initialization complete"
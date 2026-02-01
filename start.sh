#!/bin/sh

set -e

echo "Running DB migrations..."

# Check if migrate binary is present
if [ ! -f "./migrate" ]; then
    echo "❌ migrate binary not found"
    exit 1
fi

# Show migration directories
ls -R ./migration || echo "❌ migration folder missing"

# Run migrations
./migrate -path ./migration -database "$DB_SOURCE" -verbose up

echo "Starting API..."

if [ "$PRODUCTION" = "true" ]; then
    echo "Running in PRODUCTION mode"
    exec /app/main
else
    echo "Running in DEVELOPMENT mode with Air"
    exec air -c .air.toml
fi

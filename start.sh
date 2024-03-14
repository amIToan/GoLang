#!/bin/sh

set -e

echo "start the app"
# /app/migrate -path /app/migration -database "$DB_SOURCE" --verbose up 
echo "starting server with parameters : $1"
exec "$@"
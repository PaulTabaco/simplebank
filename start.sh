#!/bin/sh

set -e

echo "run db migration"
# load environments from file to current shell
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the simplebank app"
exec "$@"

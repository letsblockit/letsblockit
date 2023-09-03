#!/usr/bin/env bash
# This script runs linting and tests.
## Run it with `nix run .#run-tests`, or install the dependencies manually.

set -euox pipefail
export GOGC=400
export TEST_DATABASE_URL=${TEST_DATABASE_URL:-"postgresql:///lbitests"}

if [ "$*" != "notest" ]; then
  # Fail-fast on filter changes
  go test ./src/filters
  # Run unit and integration tests on the lbitests DB after purging it
  psql --quiet "$TEST_DATABASE_URL" -c "DROP owned BY $(whoami)"
  go test -race ./...
fi

if [ "$*" != "nolint" ]; then
  # Linting golang code
  golangci-lint run --timeout 5m
fi

echo "OK"

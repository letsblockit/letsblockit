#!/usr/bin/env bash
# This script runs linting and tests.
## Run it with `nix run .#run-tests`, or install the dependencies manually.

set -euox pipefail
export GOGC=400

# Fail-fast on filter changes
go test ./src/filters

# Full tests for backend changes
golangci-lint run --timeout 5m
go test -v -race ./...
echo "OK"

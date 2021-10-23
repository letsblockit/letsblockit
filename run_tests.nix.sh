#! /usr/bin/env nix-shell
#! nix-shell -i bash -p gcc -p go_1_16 -p golangci-lint -p sqlite
#! nix-shell --pure --quiet

# This script runs linting and tests

set -euox pipefail

golangci-lint run --build-tags=libsqlite3 --timeout=5m
go test -tags libsqlite3 -v -race ./...
go run -tags libsqlite3 main.go --dry-run

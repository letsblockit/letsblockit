#! /usr/bin/env nix-shell
#! nix-shell -i bash -p gcc -p go_1_17 -p golangci-lint
#! nix-shell --pure --quiet

# This script runs linting and tests

set -euox pipefail
export GOGC=400

golangci-lint run --timeout 5m
go test -v -race ./...
echo "OK"

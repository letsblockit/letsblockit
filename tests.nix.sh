#! /usr/bin/env nix-shell
#! nix-shell -i bash -p gcc -p go_1_16 -p golangci-lint
#! nix-shell --quiet

# This script runs linting and tests

set -euox pipefail

golangci-lint run
go test -v -race ./...
go run main.go --dry-run

#!/usr/bin/env bash
# This script updates the sqlc data access layer and mockgen mocks.
## Run it with `nix run .#update-codegen`, or install the dependencies manually.

set -euox pipefail

sqlc generate -f src/db/sqlc.yaml
mockgen -source ./src/server/deps.go -destination ./src/server/mocks/deps.go -package mocks
mockgen -source ./src/filters/deps.go -destination ./src/filters/mocks/deps.go -package mocks
mockgen -source ./src/db/interface.go -destination ./src/server/mocks/querier.go -package mocks

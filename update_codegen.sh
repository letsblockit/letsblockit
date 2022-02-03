#! /usr/bin/env nix-shell
#! nix-shell -i bash -p mockgen  -p "(callPackage ./nix/sqlc.nix {})"
#! nix-shell --quiet

# This script updates the sqlc data access layer and mockgen mocks

set -euox pipefail

sqlc generate -f src/db/sqlc.yaml
mockgen -source ./src/server/deps.go -destination ./src/server/mocks/deps.go -package mocks
mockgen -source ./src/filters/deps.go -destination ./src/filters/mocks/deps.go -package mocks
mockgen -source ./src/db/querier.go -destination ./src/server/mocks/querier.go -package mocks

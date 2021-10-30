#! /usr/bin/env nix-shell
#! nix-shell -i bash -p mockgen
#! nix-shell --quiet

# This script updates the mocks with mockgen

set -euox pipefail

mockgen -source ./src/server/deps.go -destination ./src/server/deps_test.go -package server

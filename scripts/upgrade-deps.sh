#!/usr/bin/env bash
# This script updates the vendorSha256 property in any buildGoModule
# derivation to keep it in sync with go.mod changes
## Run it with `nix run .#upgrade-dps`, or install the dependencies manually.

set -euox pipefail

# Upgrade nix flake deps
nix flake update

# Upgrade npm deps
cd src/assets
npm upgrade -S
npm prune
cd ../..

# Rebuild assets, accounting for file renames
git rm -rf data/assets/dist/
./scripts/update-assets.sh
git add data/assets/dist/

# Upgrade golang deps
go get -t -u ./...
go mod tidy
./scripts/update-vendorsha.sh

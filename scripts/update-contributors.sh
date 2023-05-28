#!/usr/bin/env bash
# This script updates the assets by running npm.
## Run it with `nix run .#update-assets`, or install the dependencies manually.

set -euox pipefail

go run ./cmd/utils/ download-avatars
go run ./cmd/utils/ hash-assets

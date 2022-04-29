#!/usr/bin/env bash
# This script wraps golang-migrate while setting the right options.
## Run it with `nix run .#add-migration`, or install the golang-migrate manually.

set -euox pipefail

OUT="./src/db/migrations"

migrate create -dir "$OUT"  -ext sql -digits 4 -seq "$@"
git add "$OUT"

#!/usr/bin/env bash
# This script wraps golang-migrate while setting the source and database
# for local development.
## Run it with `nix run .#migrate`, or install the golang-migrate manually.

set -euox pipefail

migrate -path ./src/db/migrations -database pgx:///letsblockit?host=/var/run/postgresql "$@"

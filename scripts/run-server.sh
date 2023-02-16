#!/usr/bin/env bash
# This script runs the auth proxy and the server, recompiling the server on code changes.
## Run it with `nix run .#run-server`, or install the dependencies manually.

export ORY_PROJECT_SLUG=${ORY_PROJECT_SLUG:-inspiring-gates-6gd6xqshnz}
ory proxy --port 4000 http://localhost:8765/ --no-jwt --quiet &

export LETSBLOCKIT_AUTH_METHOD=kratos
export LETSBLOCKIT_AUTH_KRATOS_URL="http://localhost:4000/.ory"
export LETSBLOCKIT_CACHE_DIR="/tmp/lbi-cache"
export LETSBLOCKIT_OFFICIAL_INSTANCE=true

mkdir -p $LETSBLOCKIT_CACHE_DIR
reflex -r "(cmd|src|data)/" -s -- go run -race ./cmd/server --hot-reload

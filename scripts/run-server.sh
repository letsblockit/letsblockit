#!/usr/bin/env bash
# This script runs the auth proxy and the server, recompiling the server on code changes.
## Run it with `nix run .#run-server`, or install the dependencies manually.

export ORY_SDK_URL=${ORY_SDK_URL:-https://inspiring-gates-6gd6xqshnz.projects.oryapis.com}
ory proxy --port 4000 http://localhost:8765/ --no-jwt &
reflex -r "(cmd|src|data)/" -s -- go run -race ./cmd/server --auth-method=kratos --hot-reload

#!/usr/bin/env bash
# This script runs the auth proxy and the server, recompiling the server on code changes.
## Run it with `nix run .#run-server`, or install the dependencies manually.

ory proxy --port 4000 http://localhost:8765/ --sdk-url https://inspiring-gates-6gd6xqshnz.projects.oryapis.com --no-jwt &
reflex -r "(cmd|src|data)/" -s -- go run -race ./cmd/server --reload

#! /usr/bin/env nix-shell
#! nix-shell -i bash -p "(callPackage ./nix/ory.nix {})"
#! nix-shell --pure --quiet

# Runs the Ory proxy on the letsblockit-dev project
ory proxy --port 4000 http://localhost:8765/ --sdk-url https://inspiring-gates-6gd6xqshnz.projects.oryapis.com --no-jwt

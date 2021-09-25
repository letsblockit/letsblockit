#! /usr/bin/env nix-shell
#! nix-shell -i bash -p "(callPackage ./ory-cli.nix {})"
#! nix-shell --pure --quiet

ory proxy local --port 4000 http://localhost:8765/ --dont-install-cert --no-open --project playground

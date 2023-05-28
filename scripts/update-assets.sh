#!/usr/bin/env bash
# This script updates the assets by running npm.
## Run it with `nix run .#update-assets`, or install the dependencies manually.

set -euox pipefail

npm install --prefix ./src/assets/
npm run lint --prefix ./src/assets/

npm run clean --prefix ./src/assets/

if [ "$*" == "watch" ]; then
  go run ./cmd/utils/ extract-icons --all
  npm run watch --prefix ./src/assets/
else
  go run ./cmd/utils/ extract-icons
  npm run build --prefix ./src/assets/
fi

go run ./cmd/utils/ hash-assets

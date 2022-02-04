#!/usr/bin/env bash
# This script updates the assets by running npm.
## Run it with `nix run .#update-assets`, or install the dependencies manually.

set -euox pipefail

npm install --prefix ./src/assets/
npm run clean --prefix ./src/assets/

if [ "$*" == "watch" ]; then
  npm run watch --prefix ./src/assets/
else
  npm run build --prefix ./src/assets/
fi

#!/usr/bin/env bash
# This script updates the assets by running npm.
## Run it with `nix run .#update-assets`, or install the dependencies manually.

set -euox pipefail

npm install --prefix ./src/assets/

if [ "$*" == "--watch" ]; then
  reflex -r "src/assets" -s -- npm run assets --prefix ./src/assets/
else
  npm run assets --prefix ./src/assets/
  git add ./data/assets/
fi

npm run assets --prefix ./src/assets/
git add ./data/assets/

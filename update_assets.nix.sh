#! /usr/bin/env nix-shell
#! nix-shell -i bash -p nodejs -p nodePackages.npm -p git
#! nix-shell --pure --quiet

# This script updates the assets by running npm

set -euox pipefail

npm install --prefix ./src/assets/
npm run assets --prefix ./src/assets/
git add ./data/assets/

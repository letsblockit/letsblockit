#! /usr/bin/env nix-shell
#! nix-shell -i bash -p nodejs -p nodePackages.npm -p git -p reflex
#! nix-shell --pure --quiet

# This script updates the assets by running npm

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

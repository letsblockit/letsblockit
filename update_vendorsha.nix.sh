#! /usr/bin/env nix-shell
#! nix-shell -i bash -p nix-prefetch
#! nix-shell --quiet

# This script updates the vendorSha256 property in the default.nix
# build derivation to keep it in sync with go.mod changes

set -euox pipefail

oldSha=$(nix-instantiate ./default.nix -A vendorSha256 --eval 2>/dev/null | sed -e 's/"//g')
newSha=$(nix-prefetch '{ sha256 }: (callPackage (import ./default.nix) { }).go-modules.overrideAttrs (_: { modSha256 = sha256; })' 2>/dev/null)

if [ "$oldSha" == "$newSha" ]; then
  echo "Nothing to update"
  exit 0
elif [ "$*" == "--check" ]; then
  echo "default.nix out of sync, please run update_vendorsha.nix.sh"
  exit 1
else
  echo "Updating vendorSha256: $oldSha -> $newSha"
fi

sed -i "s|$oldSha|$newSha|" default.nix

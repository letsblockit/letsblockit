#!/usr/bin/env bash
# This script updates the vendorSha256 property in any buildGoModule
# derivation to keep it in sync with go.mod changes
## Run it with `nix run .#update-vendorsha`, or install the dependencies manually.

set -euox pipefail

nixFile=${1:-./nix/letsblockit.nix}

nix-prefetch "{ sha256 }: (callPackage (import $nixFile) { }).goModules.overrideAttrs (_: { modSha256 = sha256; })"

oldSha=$(grep vendorHash "$nixFile" | sed 's/[^"]*"\([^"]*\).*/\1/')
newSha=$(nix-prefetch "{ sha256 }: (callPackage (import $nixFile) { }).goModules.overrideAttrs (_: { modSha256 = sha256; })" 2>/dev/null)

if [ "$oldSha" == "$newSha" ]; then
  echo "Nothing to update"
  exit 0
elif [ "${2:-none}" == "--check" ]; then
  echo "nix file out of sync, please run update_vendorsha.nix.sh"
  exit 1
else
  echo "Updating vendorHash: $oldSha -> $newSha"
fi

sed -i "s|$oldSha|$newSha|" "$nixFile"

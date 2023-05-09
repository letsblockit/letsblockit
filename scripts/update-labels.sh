#!/usr/bin/env bash
# This script updates the image labels.
## Run it with `nix run .#update-labels`, or install the dependencies manually.

set -euox pipefail

outputFile=${1:-./nix/labels.nix}

cat > "$outputFile" <<EOF
rec {
  "org.opencontainers.image.licenses" = "Apache 2.0";
  "org.opencontainers.image.source"   = "https://github.com/letsblockit/letsblockit";
  "org.opencontainers.image.url"      = "https://github.com/letsblockit/letsblockit";
  "org.opencontainers.image.vendor"   = "letsblock.it";
EOF


if [ "${GITHUB_ACTIONS:-}" == "true" ]; then
cat >> "$outputFile" <<EOF
  "org.opencontainers.image.created"  = "$(date --rfc-3339=sec)";
  "org.opencontainers.image.revision" = "${GITHUB_SHA}";
  "org.opencontainers.image.version"  = "${GITHUB_REF_NAME}";
EOF
fi

echo "}" >> "$outputFile"

#!/usr/bin/env bash

set -euox pipefail

PKG_PREFIX="ghcr.io/letsblockit"
USAGE="Usage: $0 SOURCE_IMAGE [TARGET_IMAGE]"

case "$#" in

  "1")
    docker push "$PKG_PREFIX/$1"
    ;;

  "2")
    docker tag "$PKG_PREFIX/$1" "$PKG_PREFIX/$2"
    docker push "$PKG_PREFIX/$2"
    ;;

  *)
    echo "Error: no source image provided."
    echo "$USAGE"
    exit 1
    ;;
esac
exit 0

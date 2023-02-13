# Copied from https://raw.githubusercontent.com/NixOS/nixpkgs/nixos-22.11/pkgs/tools/misc/vector/default.nix
# And modified to reduce the image size
{ stdenv
, lib
, fetchFromGitHub
, rustPlatform
, pkg-config
, llvmPackages
, openssl
, protobuf
, oniguruma
, zstd
, tzdata
, cmake
, perl
, gcc
, removeReferencesTo
}:

let
  pname = "vector";
  version = "0.26.0";
in
rustPlatform.buildRustPackage {
  inherit pname version;

  src = fetchFromGitHub {
    owner = "vectordotdev";
    repo = pname;
    rev = "v${version}";
    hash = "sha256-0h9hcNgaVBDBeSKo39TvrMlloTS5ZoXrbVhm7Y43U+o=";
  };

  cargoHash = "sha256-UHc8ZyLJ1pxaBuP6bOXdbAI1oVZD4CVHAIa8URnNdaI=";
  nativeBuildInputs = [ pkg-config cmake perl protobuf removeReferencesTo ];
  buildInputs = [ oniguruma openssl zstd ];

  PROTOC = "${protobuf}/bin/protoc";
  PROTOC_INCLUDE = "${protobuf}/include";
  RUSTONIG_SYSTEM_LIBONIG = true;
  LIBCLANG_PATH = "${llvmPackages.libclang.lib}/lib";
  TZDIR = "${tzdata}/share/zoneinfo";

  buildNoDefaultFeatures = true;
  buildFeatures = [
    "enterprise"
    "sinks-blackhole"
    "sinks-datadog_logs"
    "sinks-datadog_metrics"
    "sources-demo_logs"
    "sources-host_metrics"
    "sources-exec"
    "sources-statsd"
    "sources-stdin"
    "transforms"
  ];

  doCheck = false;

  postInstall = ''
    find "$out" -type f -exec remove-references-to -t ${stdenv.cc} '{}' +
  '';
}

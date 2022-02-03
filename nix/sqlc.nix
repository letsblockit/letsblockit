{ buildGoModule, fetchFromGitHub, go_1_17 }:
buildGoModule.override { go = go_1_17; } rec {
  pname = "sqlc";
  version = "1.11.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-EzV0h5YPaZdsPFjihXX5gDMiWqlCKFVlN39c9/eerAU=";
  };

  vendorSha256 = "sha256-lxOrCRyk3zHRw8WZLs0p7pN6dCyOM2pqy3bax8+PCi0=";
  runVend = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

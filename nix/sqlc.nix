{ buildGoModule, fetchFromGitHub, go_1_17 }:
buildGoModule.override { go = go_1_17; } rec {
  pname = "sqlc";
  version = "1.12.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-YlOkjqkhN+4hL1+KJ0TuqcbQXJad/bHZclgpgFPr4to=";
  };

  vendorSha256 = "sha256-MX140KUb+FGtorzOR46NHlWFhfzvlsSbg41sn3XR8Ys=";
  runVend = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

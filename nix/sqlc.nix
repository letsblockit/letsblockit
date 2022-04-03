{ buildGoModule, fetchFromGitHub, go_1_17 }:
buildGoModule.override { go = go_1_17; } rec {
  pname = "sqlc";
  version = "1.13.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-HPCt47tctVV8Oz9/7AoVMezIAv6wEsaB7B4rgo9/fNU=";
  };

  vendorSha256 = "sha256-ZEVtc5FMiRTuTLtgbYJeuIWGYXKGbxZZ148hjuzN2wM=";
  runVend = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

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

  vendorSha256 = "sha256-zZ0IrtfQvczoB7th9ZCUlYOtyZr3Y3yF0pKzRCqmCjo=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

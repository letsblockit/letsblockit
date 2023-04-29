{ buildGoModule, fetchFromGitHub, go_1_19 }:
buildGoModule.override { go = go_1_19; } rec {
  pname = "sqlc";
  version = "1.18.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-5MC7D9+33x/l76j186FCnzo0Hnx0wY6BPdneW7E7MpE=";
  };

  vendorSha256 = "sha256-gDePB+IZSyVIILDAj+O0Q8hgL0N/0Mwp1Xsrlh3B914=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

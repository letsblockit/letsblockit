{ buildGoModule, fetchFromGitHub, go_1_19 }:
buildGoModule.override { go = go_1_19; } rec {
  pname = "sqlc";
  version = "1.17.1";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-lz9Y4HyCwJEB+OR/a02eB0Xr91NC3l3ANeqYf6Zq2Kg=";
  };

  vendorSha256 = "sha256-y5OYq1X4Y0DxFYW2CiedcIjhOyeHgMhJ3dMa+2PUCUY=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

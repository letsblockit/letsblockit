{ buildGoModule, fetchFromGitHub, go_1_19 }:
buildGoModule.override { go = go_1_19; } rec {
  pname = "sqlc";
  version = "1.17.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-knblQwO+c8AD0WJ+1l6FJP8j8pdsVhKa/oiPqUJfsVY=";
  };

  vendorSha256 = "sha256-y5OYq1X4Y0DxFYW2CiedcIjhOyeHgMhJ3dMa+2PUCUY=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

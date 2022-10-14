{ buildGoModule, fetchFromGitHub, go_1_18 }:
buildGoModule.override { go = go_1_18; } rec {
  pname = "sqlc";
  version = "1.14.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-+JkNuN5Hv1g1+UpJEBZpf7QV/3A85IVzMa5cfeRSQRo=";
  };

  vendorSha256 = "sha256-eMghAqOiyX/EbXg/Q3Bxb3Xx8N5ekFGBwn1AMIGY+hw=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

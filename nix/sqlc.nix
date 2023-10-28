{ buildGoModule, fetchFromGitHub, pinnedGo }:
buildGoModule.override { go = pinnedGo; } rec {
  pname = "sqlc";
  version = "1.19.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-/6CqzkdZMog0ldoMN0PH8QhL1QsOBaDAnqTHlgtHdP8=";
  };

  vendorSha256 = "sha256-AsOm86apA5EiZ9Ss7RPgVn/b2/O6wPj/ur0zG91JoJo=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

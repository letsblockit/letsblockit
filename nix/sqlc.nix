{ buildGoModule, fetchFromGitHub }:
buildGoModule rec {
  pname = "sqlc";
  version = "1.24.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    sha256 = "sha256-j+pyj1CJw0L3s4Nyhy+XXUgX2wbrOWveEJQ4cFhQEvs=";
  };

  vendorHash = "sha256-xOMqZCuENGuCs+VkbCxMpXOEr4MALhlveTfUHEPnP1w=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

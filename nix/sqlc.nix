{ buildGoModule, fetchFromGitHub }:
buildGoModule rec {
  pname = "sqlc";
  version = "1.25.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v${version}";
    hash = "sha256-VrR/oSGyKtbKHfQaiLQ9oKyWC1Y7lTZO1aUSS5bCkKY=";
  };

  vendorHash = "sha256-C5OOTAYoSt4anz1B/NGDHY5NhxfyTZ6EHis04LFnMPM=";
  proxyVendor = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}

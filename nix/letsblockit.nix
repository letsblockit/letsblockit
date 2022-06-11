{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-QDKs71bruDNJrsuBppZkzZM6432NV8DcSb2Wc71E6jU=";
  version = "1.0";
}

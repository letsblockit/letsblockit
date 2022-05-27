{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-upi0jlt6FAtlwscm+5SJyUUvz6wdd1/U+6nuUe58Kq4=";
  version = "1.0";
}

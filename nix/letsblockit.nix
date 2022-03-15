{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-YLyKspZB58aHLO8AvfZ90z3xnENgN0WcpSRYsynf93Q=";
  version = "1.0";
}

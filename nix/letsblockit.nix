{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-DpitG3WXytBV4TdpQnwsGiVuj4dcoCYTlE2fIp9fh1Y=";
  version = "1.0";
}

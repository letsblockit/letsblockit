{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-Ms3LHM58RryRrVeIjr14KdTx+eEN/na/oLyTLVrzW1g=";
  version = "1.0";
}

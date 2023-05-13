{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-h/iPX3PZQjiaREkar57YBE3bZulGnL78XCQnig+E8aM=";
  version = "1.0";
}

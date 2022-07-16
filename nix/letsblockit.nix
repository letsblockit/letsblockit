{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-K0ewO8vYPQ/0umBQX4h63q8mqQqmqA233HyMF8KFX9g=";
  version = "1.0";
}

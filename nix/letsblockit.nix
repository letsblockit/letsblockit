{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-7f+bW3P/7zwvtjlAVwYceIOOd89IVGrpTsTHNbIR4XE=";
  version = "1.0";
}

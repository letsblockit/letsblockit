{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-pjU/Q7eQo06gXzYi9tXcLrzm9gtGcdCnZ05Nh1Ysqz8=";
  version = "1.0";
}

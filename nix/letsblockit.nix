{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-PALfKDOu6qpKBmIGuj7u5WtKZGS+xcJzFry/Q5Z6hX0=";
  version = "1.0";
}

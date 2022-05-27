{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-LIwfRZmD1buCYwLWCFOiR29+L0gRyrM9w3+JcsIe5Kk=";
  version = "1.0";
}

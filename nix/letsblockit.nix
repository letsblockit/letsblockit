{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-m5+9DFZFiPhzg1KcPr4tAaEwO9bbHzNTZv++cOZKORM=";
  version = "1.0";
}

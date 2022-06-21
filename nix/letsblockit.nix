{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-vkXkZAFJyA7BAY8mxgN2GYDgd7cxDyI1TK9fI+lqp8Q=";
  version = "1.0";
}

{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-gAps5grMcM0K+WIPyStgK9gG+F5mZJdaMxpQJmlkmLk=";
  version = "1.0";
}

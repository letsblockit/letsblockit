{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-+Kas8YuiJcV6iGwHF7SmGZ7KYkWKAHagE86MjZNEcyE=";
  version = "1.0";
}

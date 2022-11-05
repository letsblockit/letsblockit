{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-uLDbsF4qkXw6NXNovzgvOuHw3mQT0ozb7lPfh3l9mkg=";
  version = "1.0";
}

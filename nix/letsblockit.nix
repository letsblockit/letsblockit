{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-RgNKYoXCbRoCKd7HCc4SugQGhKD4/Hi5t030BO86Ok8=";
  version = "1.0";
}

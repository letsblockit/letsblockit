{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-509/PRcek+wZqDXZyczoOusMB7R4EMeAPUN720O12Yo=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-qxqrDigjielSF5S2OTyXeLsVLe1imHHyNXsAEAh1hOc=";
  version = "1.0";
}

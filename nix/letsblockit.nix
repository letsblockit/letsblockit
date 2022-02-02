{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-ZMGTL58tCbyG1QvmTFzTd+mkJ1t2/403e5podNr3aOY=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-Tr6vRy1sQNQ2eMYWNa6MbCNcePhKT8AND2W7J3P3wz8=";
  version = "1.0";
}

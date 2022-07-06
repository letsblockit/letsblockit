{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-ggTfw2RKe5oEDPM1Iyo4di0PdFmBsttmbVA01FMuy0E=";
  version = "1.0";
}

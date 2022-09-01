{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-NCpjoabxCf6gZ/RyePXPdnsMwjKsO2PonbXWX2jiqNk=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-Qs5DRP7mlyLn5b5+GMabyGtcwtcBmj0gTA5+iO1iLGY=";
  version = "1.0";
}

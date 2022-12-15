{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-SNexwjs50zUGcBAvul2oo22zxON6QMePMaWgOlbj6zI=";
  version = "1.0";
}

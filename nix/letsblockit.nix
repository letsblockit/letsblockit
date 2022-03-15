{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-3U3UNWNOXOXcuhM2azHeZFS6Jgqka/AZAkOw/aOgzaY=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-LaQ8yHsadkbmSpaJhi8qd+qwV5JZIFY9J5UVW0Wexro=";
  version = "1.0";
}

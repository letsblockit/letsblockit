{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-Twgl1iTRMXMxGwNk66msFyPkM6og418HtkGg7iy9Ra0=";
  version = "1.0";
}

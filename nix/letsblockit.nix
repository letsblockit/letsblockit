{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-+eZZx8Y2sx+BU9p8yYStHykh9WHk7ICc5q8yWk45N0E=";
  version = "1.0";
}

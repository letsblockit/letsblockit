{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-kdCLI+B+fk0N+VgRCGJehVyt51UffcBq3i7Erd2qfG8=";
  version = "1.0";
}

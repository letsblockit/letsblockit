{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-m/Zg/i4NlA6mX4sByT7swmKSJcxFDqLK4lrmyje9jNE=";
  version = "1.0";
}

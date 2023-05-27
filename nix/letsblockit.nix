{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-vQQUrALQAOI2WEfXpxo7lQk3al2OXp1Ay5kzGimX9gE=";
  version = "1.0";
}

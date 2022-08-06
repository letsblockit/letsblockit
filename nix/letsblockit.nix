{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-vA6Yc8DujRtsG8mpkqWlt+fNmJ4Z7l9KLH8PtF8pdb4=";
  version = "1.0";
}

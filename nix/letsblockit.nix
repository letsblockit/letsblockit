{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-1fggyEFk64RT5oJOVLOAXejCCYfhhD5OW1sXwq4QtZU=";
  version = "1.0";
}

{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-uvDVX/1Tg2llfdR9hmQzri/rmiH7ELgwqditCyrMZfk=";
  version = "1.0";
}

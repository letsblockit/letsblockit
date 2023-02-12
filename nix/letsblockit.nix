{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-FgkaFgDf6mQm9PE7Tm3ciZb8dz9CaHirIbHUqETlQFY=";
  version = "1.0";
}

{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-CQFHk6ejVzkU3WCYPWviG/pxT0OrNP2u4ZWr0PDOC1Y=";
  version = "1.0";
}

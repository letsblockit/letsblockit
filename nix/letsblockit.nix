{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-XGfbHEBP9iB6+NIsuQrPN7u/bk2d4sjFvvkyXCtr2dY=";
  version = "1.0";
}

{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-S2c6KkMAur643FRl5OaNIAEYk3Fd745JLJuGH0UAYCM=";
  version = "1.0";
}

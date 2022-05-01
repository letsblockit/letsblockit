{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-lIyEgmPvFSUymfpj0fvMndHorkp7SrDtZSb5TsIggvU=";
  version = "1.0";
}

{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-gXdU/qTmGV1RPpxJs2NclCfCvz1z2lnCHkutSSwgSxg=";
  version = "1.0";
}

{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-GVC5UMPQvsJeHYlcCtmNYZCAI1PM+qLj5kdh19yH2yE=";
  version = "1.0";
}

{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-P/qeXPgtc82ga3fqFFx1n/n9mS9DFEedr6dnWR5pVwY=";
  version = "1.0";
}

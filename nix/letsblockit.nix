{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-2UmnAHg+A0ehOY2Kkfg837eWDNNNdyqISARC/hqe9tU=";
  version = "1.0";
}

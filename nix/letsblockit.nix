{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-c/9monkRFCPPhilCqYYWJKcQ7hTlnBhMyWNUrxT+vQY=";
  version = "1.0";
}

{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-2YuJgkRpm1HldTAY9kHNUxaXsUhqsuNE8OFcqtKRdf4=";
  version = "1.0";
}

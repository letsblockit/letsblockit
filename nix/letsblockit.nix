{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-farWDLfdsM2T5kYzRf5MH9Vf6riHKYEOFOE1BhODcYA=";
  version = "1.0";
}

{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-din/oMTUwOGXK+mbiecRTLDsXmL/gex0F3USTrwZxsU=";
  version = "1.0";
}

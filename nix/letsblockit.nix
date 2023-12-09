{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-WddUHyGbjFWeUUbVxKTeG4PbVgiPMnsZuKHmdCxHArQ=";
  version = "1.0";
}

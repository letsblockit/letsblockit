{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-C6ApTMLyPvlda3tK6g9JE/QfKUA7A/earayW1bABPUA=";
  version = "1.0";
}

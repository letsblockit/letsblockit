{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-9fSnsIB0nxl7gmJaDBx3V24adwDY5G20w2tuurt0anw=";
  version = "1.0";
}

{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-HLxXL+XwgGCGPrHF7oNTOBoY9YoYJWqkkrgvdCq/Gb8=";
  version = "1.0";
}

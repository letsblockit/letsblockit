{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-rrWEtGBfDSTfyXJ7nqD8pqb42iwwvHP+NEs4md8kmKI=";
  version = "1.0";
}

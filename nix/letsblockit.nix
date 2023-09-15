{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-waScF4pcDgbZ4CkVunZYZog48pUe+tS6MDcBuYMx2zA=";
  version = "1.0";
}

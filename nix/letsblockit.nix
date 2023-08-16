{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-c3+LskTwKOK5i9+gV4tpWEPmbIMMrIRHszC0PKTKtwQ=";
  version = "1.0";
}

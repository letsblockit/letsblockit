{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-bGL2fmBXKDl0hvZqK2RBneZG75yFuwh/mFCJvHn9yKU=";
  version = "1.0";
}

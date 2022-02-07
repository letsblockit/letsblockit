{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-4LTSlbOF622/Wclm2C87gn0Ii8EyNKIVrxpc7D5XZlA=";
  version = "1.0";
}

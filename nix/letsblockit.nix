{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-8DQd1EBGiBvT+YeSVu/UeIZyNDuBSPbsgO7L2knQtoc=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-mLxYMwTROuAZ+J3xEkP+lNDpUxia5q+inCDGYqqxg+4=";
  version = "1.0";
}

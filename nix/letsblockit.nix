{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-XTgdqcRlvS0yxyKTHR/reAVsLD/rtNe3wBpi33bhQI8=";
  version = "1.0";
}

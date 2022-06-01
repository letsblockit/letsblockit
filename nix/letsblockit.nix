{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-mPEPFhldSG4qszNM+Qthgfv6XNdobwbL3OEsSF1Jei8=";
  version = "1.0";
}

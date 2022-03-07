{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-i+GqtjSnYShu46CjPt93QW/a1/5TsHs/Y68F7sZGKSQ=";
  version = "1.0";
}

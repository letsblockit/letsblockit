{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-W6gPj4E9K8tVgv2mcBn38IiXRxJc/yaszYC3tBGDm+A=";
  version = "1.0";
}

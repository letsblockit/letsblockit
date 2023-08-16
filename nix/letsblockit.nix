{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-sMNoBdT8jyNx08PRueYL3xmcgFAXqzLH2n2cjR7rWRY=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-FziB0CLtgtpJiyY8iHFHFQA5WDdpnXS+jbkkA1vmnB8=";
  version = "1.0";
}

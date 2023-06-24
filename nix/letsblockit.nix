{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-k9XoqSwtwHXCKRTd+kp5YDU/3lJw3JuMRVz1GA0TPnQ=";
  version = "1.0";
}

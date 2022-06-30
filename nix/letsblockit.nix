{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-CFSMxaohP6H3uHQZndRmQKld3FQJHwvb5Ief7OXa4aU=";
  version = "1.0";
}

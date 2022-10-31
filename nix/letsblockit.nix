{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-142XCOBzAO/wi9fahD1xoLYluXhTYLeTtmTVhfH6c5M=";
  version = "1.0";
}

{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-pLtwFR4o1hzpqLFfTDYpBXfWge+EgwNn9X6gWh1dbu4=";
  version = "1.0";
}

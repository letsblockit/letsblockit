{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-3g/LLaKN+4gkNVEslHA/YAYqHDmnMWM1vXckaupqQwI=";
  version = "1.0";
}

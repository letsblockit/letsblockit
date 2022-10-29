{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-CFzXsDveqMk765jaFk5AY8AyYJLPbGXbiv67SM8Oa+Q=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-WFsNvemHt+JhXV4TK8c1s6XWZ9E2QNgAT7arqqijjXs=";
  version = "1.0";
}

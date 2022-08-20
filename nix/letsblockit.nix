{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-icJDsc3MhYwJwxXXwH558QNH+ohQiZzmX5cPRyT0glc=";
  version = "1.0";
}

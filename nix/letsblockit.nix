{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-TBjK3Aa9wqjOXC+RVaxshOSYvLdgCNI+Kvn02sduhN0=";
  version = "1.0";
}

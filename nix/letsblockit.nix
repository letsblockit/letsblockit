{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-i7mxNdhw5w4bvDuAEScmvdQEKzFIGQbC/EXe6t2CTAw=";
  version = "1.0";
}

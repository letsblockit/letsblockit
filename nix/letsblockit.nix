{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-nnkfEA3y0FO1InRao2+mMtvQ5Oof4Ex8kDQLNer7LcA=";
  version = "1.0";
}

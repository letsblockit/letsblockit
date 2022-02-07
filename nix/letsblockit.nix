{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-m5fMYwLQckM0v2O20XtO8e80B3gmuWLNkw4/U48jtFY=";
  version = "1.0";
}

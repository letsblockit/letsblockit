{ buildGoModule, pinnedGo, cmd ? "server" }:
buildGoModule.override { go = pinnedGo; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-j0Wqj9/8/9Me/Q+hA2buOfUjfKX0hzEnap8MlP90DAo=";
  version = "1.0";
}

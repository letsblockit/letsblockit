{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-sd3dCf9TaDFyznHBphTLU6jU2DkQOOCLorynxCwQhy0=";
  version = "1.0";
}

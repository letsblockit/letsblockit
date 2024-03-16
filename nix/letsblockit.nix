{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-Il4othd7km6kizBpsR3aOx2D3Gr/kmwy/SXQW4U1wLo=";
  version = "1.0";
}

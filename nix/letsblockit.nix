{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-Pv0W3JnzkEgy54/CKv08kRr/RPpk27Hg0EjUPhUmCGQ=";
  version = "1.0";
}

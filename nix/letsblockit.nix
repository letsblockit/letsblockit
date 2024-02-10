{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-FtBMU8FSzgqyKG/H3MdnrEfEDmLFt0T7Ti2LeUP0JXs=";
  version = "1.0";
}

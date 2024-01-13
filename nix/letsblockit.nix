{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-IGqWEQ6IqGnx5WDdqLiYEasYTKCmLx3sR5LvcwcGzn0=";
  version = "1.0";
}

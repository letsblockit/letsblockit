{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-vZmO93/doz8BZHzqzOW5PxX9GWu8YDuGZ2Lpb5Bru2M=";
  version = "1.0";
}

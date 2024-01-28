{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-H8Y3a03aujES6w+Pn7oxxWlCW6T4NJK/LwThbqOxDaU=";
  version = "1.0";
}

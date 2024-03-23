{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-0iJWVFAPKeMsBqC79+AR+J4wAl3M7V6qwJ7JzvHuZh0=";
  version = "1.0";
}

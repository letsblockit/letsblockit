{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-Hqz+aWB3BRlyOH3Z5RXlYe6YDEQ7Ztrqk7xN7WsfD8Q=";
  version = "1.0";
}

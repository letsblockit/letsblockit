{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-hdz263mKawGOJQHoX5FQL707cNTC9QR9WpTUaA4UnmQ=";
  version = "1.0";
}

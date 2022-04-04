{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-74iOsh1wEwryRiN27IxRlikBEUY0KwNevd77IWcwpiA=";
  version = "1.0";
}

{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-vbbv0r3+PZAR8sciXTxqOSg8kWfXn/1vw7xJfY/JelA=";
  version = "1.0";
}

{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-0rc87DCJpgNMEErrCRmy34WtRuPA1CU6tKd7iZzOCIM=";
  version = "1.0";
}

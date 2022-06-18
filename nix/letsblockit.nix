{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-Q4aDLz3SD0f1dUmtMig5vWk9PqKFlLMSNdjTcKfDuQM=";
  version = "1.0";
}

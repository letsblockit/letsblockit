{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-pNlnY6MZ5OdAvepmEQmkMcZt7m8ZTe9HQSlbcxyE9is=";
  version = "1.0";
}

{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-ENQ+Kwd8CU3FHkWLNn8RNrTXOyx5D/bffcZSGvecBTk=";
  version = "1.0";
}

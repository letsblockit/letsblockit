{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-WJp53Xq8qqsRqiV9z7gr8YPCQcjJwQlnxM5Eu8fpnro=";
  version = "1.0";
}

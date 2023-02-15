{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-dS3bph8S+Bwuh2ie6vm/ZMQdg+7pjEpCGslPLdE1QLw=";
  version = "1.0";
}

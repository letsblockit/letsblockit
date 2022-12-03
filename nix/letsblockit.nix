{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-F9YRhVOvpf3EXmtVVfXznRu3xCebKMV6j/Rp4OfZzCE=";
  version = "1.0";
}

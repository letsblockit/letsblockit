{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-s63/DK3UHLpyrtDNDEYIFeTKWOaE1UexZG7jMgxvvHc=";
  version = "1.0";
}

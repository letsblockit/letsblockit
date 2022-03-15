{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-8vOJfAW5mUmk7I5oPb/YjoLiaJaPf7ZwZ5hBJVILJuU=";
  version = "1.0";
}

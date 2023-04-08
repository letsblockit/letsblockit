{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-WE6eVfVS6YvSVJEYehOxzg0akDz5YcFs0RA78p+fYGc=";
  version = "1.0";
}

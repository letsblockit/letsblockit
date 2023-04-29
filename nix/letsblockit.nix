{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-ssT0YxNo28t41JNiXVUuxcVI1YN9iwYVdCK9un6NKZ0=";
  version = "1.0";
}

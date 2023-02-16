{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-nkC2RyPabWuFzeKViy9cFx5wpZplHlXnRHvTKtGlY8M=";
  version = "1.0";
}

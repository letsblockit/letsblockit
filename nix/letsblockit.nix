{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-t1fb3EEEULnQo5a598BcBnZNjIHeB43oU4Q0dNCcqRs=";
  version = "1.0";
}

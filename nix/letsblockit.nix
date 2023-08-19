{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-1MM8gZdzwzp5Wbiu46Qz/lELzSBzZU3ZVDt+FhhoHPs=";
  version = "1.0";
}
